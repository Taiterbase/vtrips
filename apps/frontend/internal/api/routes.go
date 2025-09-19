package api

import (
    "net/http"

    "github.com/Taiterbase/vtrips/apps/frontend/web/views"
    "github.com/a-h/templ"
    "github.com/labstack/echo/v4"
)

// registerRoutes defines application routes.
func registerRoutes(e *echo.Echo) {
    // Voluntrips-ported public pages with HTMX partial support
    e.GET("/", componentOrPartial(views.VolunHomePage(false), views.VolunHomeContent()))
    e.GET("/browse", componentOrPartial(views.VolunBrowsePage(false), views.VolunBrowseContent()))
    e.GET("/likes", componentOrPartial(views.VolunLikesPage(true), views.VolunLikesContent()))

    // Modals and dropdowns
    e.GET("/modal/login", componentHandler(views.ModalLogin()))
    e.GET("/modal/sign-up", componentHandler(views.ModalSignup()))
    e.GET("/modal/search", componentHandler(views.DropdownSearch()))
    e.GET("/modal/more", componentHandler(views.DropdownMore()))
    e.GET("/modal/user", func(c echo.Context) error {
        // optional layout param, not used for rendering decision here
        loggedIn := c.QueryParam("user") == "1" // simple toggle if needed
        h := templ.Handler(views.DropdownUser(loggedIn))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })

    // Additional pages (support HTMX partials for dropdown links)
    e.GET("/about", componentOrPartial(views.VolunAboutPage(), views.VolunStaticContent("About", "About page (placeholder)")))
    e.GET("/developers", componentOrPartial(views.VolunDevelopersPage(), views.VolunStaticContent("Developers", "Developers page (placeholder)")))
    e.GET("/tos", componentOrPartial(views.VolunTOSPage(), views.VolunStaticContent("Terms of Service", "Terms of Service (placeholder)")))
    e.GET("/policy", componentOrPartial(views.VolunPolicyPage(), views.VolunStaticContent("Privacy Policy", "Privacy Policy (placeholder)")))

    // Trips dashboard-like page (protected look & feel)
    e.GET("/trips", componentHandler(views.VolunTripsPage("Friend")))

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

// componentOrPartial renders a full page normally, but returns only the partial body when HTMX requests it.
func componentOrPartial(full, partial templ.Component) echo.HandlerFunc {
    return func(c echo.Context) error {
        if c.Request().Header.Get("HX-Request") == "true" || c.Request().Header.Get("X-HX-Request") == "true" {
            h := templ.Handler(partial)
            h.ServeHTTP(c.Response().Writer, c.Request())
            return nil
        }
        h := templ.Handler(full)
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    }
}
