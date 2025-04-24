package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.Handler(http.MethodGet, "/", http.FileServer(http.Dir("./static")))

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/pack/calculate/:size", app.calculatePackSizeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/pack/size", app.getPackSizeHandler)
	router.HandlerFunc(http.MethodPost, "/v1/pack/size", app.addPackSizeHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/pack/size/:size", app.deletePackSizeHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
