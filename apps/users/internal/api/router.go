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

func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}

func StartAPI() {
	e := echo.New()
	e.Debug = false
	e.Logger.SetLevel(log.DEBUG)
	e.Debug = true
	e.HideBanner = true
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${uri} ${latency_human}\n",
	}))
	e.Use(middleware.Recover())

	e.Validator = &Validator{validator: validator.New()}
	setupRouters(e)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", "0.0.0.0", "8080")))
}

func setupRouters(eng *echo.Echo) {
	eng.GET("/v1/users", GetUsers)

	// item operations
	eng.POST("/v1/users", CreateUser)
	eng.GET("/v1/users/:user_id", GetUser)
	eng.PUT("/v1/users/:user_id", UpdateUser)
	eng.DELETE("/v1/users/:user_id", DeleteUser)

	authGroup := eng.Group("/v1/users/auth")
	authGroup.POST("/login", LoginHandler)
	authGroup.POST("/signup", SignUpHandler)
	authGroup.POST("/logout", LogoutHandler)
}
