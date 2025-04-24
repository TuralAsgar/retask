package data

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"sort"
	"time"
)

type CalculatorModel struct {
	DB *sql.DB
}

// dpState stores the minimum packs required to reach a specific item count
// and the last pack size used to achieve it.
type dpState struct {
	totalPacks   int
	lastPackUsed int
}

// CalculatePacks determines the optimal pack combination for a given order amount.
// It follows the rules:
// 1. Only whole packs.
// 2. Minimize total items sent (must be >= orderAmount).
// 3. Minimize number of packs (for the chosen minimum total items).
func (m *CalculatorModel) CalculatePacks(orderAmount int) (map[int]int, error) {
	packSizes, err := m.GetAllPacks()
	if err != nil {
		return nil, err
	}

	if orderAmount <= 0 {
		return nil, fmt.Errorf("order amount should be greater than zero")
	}

	// Ensure pack sizes are sorted (should be handled by constructor, but double-check)
	if !sort.IntsAreSorted(packSizes) {
		sort.Ints(packSizes)
	}

	// If no pack sizes are available, cannot fulfill the order.
	if len(packSizes) == 0 {
		return nil, fmt.Errorf("pack sizes from db could not be found, did you forget seeding?")
	}

	// --- Dynamic Programming Approach ---

	// Determine the maximum amount to calculate DP for.
	// We need to check amounts >= orderAmount. The minimum possible total items
	// won't exceed orderAmount + largestPackSize - 1.
	// Example: order=251, packs=[250, 500]. Max check needed is 500. DP up to 500.
	// Example: order=1, packs=[250, 500]. Max check needed is 250. DP up to 250.
	// Upper bound: orderAmount + largest pack size seems safe.
	largestPackSize := 0
	if len(packSizes) > 0 {
		largestPackSize = packSizes[len(packSizes)-1]
	}
	// Prevent potential overflow if orderAmount is huge, though unlikely with typical inputs.
	// Cap maxAmount reasonably if necessary, or use math/big. For this challenge, int should suffice.
	maxAmount := orderAmount + largestPackSize

	// Initialize DP table. dp[i] stores the best way to sum exactly to i items.
	dp := make([]dpState, maxAmount+1)
	for i := range dp {
		dp[i] = dpState{totalPacks: math.MaxInt32, lastPackUsed: 0} // Initialize with "infinity"
	}
	dp[0] = dpState{totalPacks: 0, lastPackUsed: 0} // Base case: 0 items need 0 packs

	// Fill the DP table
	for amount := 1; amount <= maxAmount; amount++ {
		for _, packSize := range packSizes {
			if amount >= packSize && dp[amount-packSize].totalPacks != math.MaxInt32 {
				// If we can reach (amount - packSize), consider adding this packSize
				candidatePacks := dp[amount-packSize].totalPacks + 1
				// If this path uses fewer packs than the current best for 'amount', update dp[amount]
				if candidatePacks < dp[amount].totalPacks {
					dp[amount] = dpState{
						totalPacks:   candidatePacks,
						lastPackUsed: packSize,
					}
				}
			}
		}
	}

	// Find the minimum total items >= orderAmount that is achievable
	minTotalItems := -1

	for i := orderAmount; i <= maxAmount; i++ {
		if dp[i].totalPacks != math.MaxInt32 {
			// Found an achievable amount >= orderAmount
			minTotalItems = i
			break // Since we iterate upwards, the first one found has the minimum total items
		}
	}

	// If no solution found (shouldn't happen with positive packs and order > 0)
	if minTotalItems == -1 {
		// This indicates an issue, maybe maxAmount was too small or no packs available.
		// For this problem, it likely means orderAmount > 0 but no packs were provided.
		return nil, fmt.Errorf("maxAmount was too small or no packs available")
	}

	// Reconstruct the pack combination by backtracking
	result := make(map[int]int)
	currentAmount := minTotalItems
	for currentAmount > 0 && dp[currentAmount].totalPacks > 0 {
		packUsed := dp[currentAmount].lastPackUsed
		if packUsed <= 0 {
			// Should not happen if dp table is correct and minTotalItems > 0
			break // Avoid infinite loop if something went wrong
		}
		result[packUsed]++
		currentAmount -= packUsed
	}

	return result, nil
}

func (m *CalculatorModel) GetAllPacks() ([]int, error) {
	query := `
        SELECT size
        FROM pack_size
        `

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	packSizes := []int{}

	for rows.Next() {
		var packSize int

		err = rows.Scan(
			&packSize,
		)
		if err != nil {
			return nil, err
		}

		packSizes = append(packSizes, packSize)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return packSizes, nil

}

func (m *CalculatorModel) InsertPack(amount int) error {
	query := `
        INSERT INTO pack_size (size) 
        VALUES (?)
        RETURNING size`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, amount).Scan(&amount)
}

func (m *CalculatorModel) DeletePack(amount int) error {

	if amount <= 0 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM pack_size
			  WHERE size = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, amount)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
