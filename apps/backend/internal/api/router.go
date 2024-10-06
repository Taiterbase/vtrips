package api

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func StartAPI() {
	e := echo.New()
	e.Debug = false
	e.Logger.SetLevel(log.DEBUG)
	e.Debug = true
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${uri} ${latency_human}\n",
	}))
	e.Use(middleware.Recover())

	e.Validator = &Validator{validator: validator.New()}
	setupRouters(e)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", "0.0.0.0", "8080")))
}

func setupRouters(eng *echo.Echo) {
	eng.GET("/v1/trips", ListTrips)

	// bulk operations
	eng.PUT("/v1/trips", UpdateTrips)
	eng.DELETE("/v1/trips", DeleteTrips)

	// item operations
	eng.POST("/v1/trips", CreateTrip)
	eng.GET("/v1/trips/:trip_id", GetTrip)
	eng.PUT("/v1/trips/:trip_id", UpdateTrip)
	eng.DELETE("/v1/trips/:trip_id", DeleteTrip)
}
