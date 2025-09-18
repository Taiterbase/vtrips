package server

import (
	"net/http"

	"github.com/Taiterbase/vtrips/apps/frontend/web/views"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// registerRoutes defines application routes.
func registerRoutes(e *echo.Echo) {
	e.GET("/", componentHandler(views.HomePorted()))
	e.GET("/browse", componentHandler(views.BrowsePorted()))

	org := e.Group("/org")
	org.GET("/dashboard", componentHandler(views.OrgDashboardPage()))

	trips := org.Group("/trips")
	trips.GET("", componentHandler(views.TripsIndexPage()))
	trips.GET("/new", componentHandler(views.TripNewPage()))
	trips.POST("", func(c echo.Context) error {
		// Placeholder for creating a trip
		return c.Redirect(http.StatusSeeOther, "/org/trips")
	})

	apps := org.Group("/applications")
	apps.GET("", componentHandler(views.ApplicationsIndexPage()))
	apps.GET("/:id", componentHandler(views.ApplicationShowPage()))
	apps.POST("/:id/grade", func(c echo.Context) error {
		// Placeholder for grading/processing
		return c.Redirect(http.StatusSeeOther, "/org/applications")
	})
}

// componentHandler wraps a templ.Component for use with Echo.
func componentHandler(cmp templ.Component) echo.HandlerFunc {
	return func(c echo.Context) error {
		h := templ.Handler(cmp)
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}
