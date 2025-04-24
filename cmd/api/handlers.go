package main

import (
	"errors"
	"github.com/TuralAsgar/dynamic-programming/internal/data"
	"net/http"
)

func (app *application) calculatePackSizeHandler(w http.ResponseWriter, r *http.Request) {
	orderAmount, err := app.readSizeParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	packs, err := app.models.Calculator.CalculatePacks(orderAmount)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"packages": packs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getPackSizeHandler(w http.ResponseWriter, r *http.Request) {
	packSizes, err := app.models.Calculator.GetAllPacks()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sizes": packSizes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addPackSizeHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Size int `json:"size"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Calculator.InsertPack(input.Size)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"size": input.Size}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePackSizeHandler(w http.ResponseWriter, r *http.Request) {
	size, err := app.readSizeParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Calculator.DeletePack(size)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "size successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
