package api

import (
    "io"
    "net/http"
    "net/url"
    "os"

    "github.com/Taiterbase/vtrips/apps/frontend/web/views"
    "github.com/a-h/templ"
    "github.com/labstack/echo/v4"
)

// registerRoutes defines application routes.
func registerRoutes(e *echo.Echo) {
    // Voluntrips-ported public pages with HTMX partial support (templ components)
    e.GET("/", func(c echo.Context) error {
        // Always render full page; internal partial swaps target #partial in layout
        h := templ.Handler(views.HomePage(isLoggedIn(c)))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })
    e.GET("/browse", func(c echo.Context) error {
        h := templ.Handler(views.BrowsePage(isLoggedIn(c)))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })
    e.GET("/likes", func(c echo.Context) error {
        h := templ.Handler(views.LikesPage(isLoggedIn(c)))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })

    // Modals and dropdowns
    e.GET("/modal/login", componentHandler(views.LoginModal()))
    e.GET("/modal/sign-up", componentHandler(views.SignupModal()))
    // Auth proxy endpoints -> users service
    e.POST("/auth/sign-up", proxyUsersSignup)
    e.POST("/auth/login", proxyUsersLogin)
    e.PUT("/auth/login", proxyUsersLogin) // support hx-put from modal
    e.POST("/auth/logout", proxyUsersLogout)
    e.GET("/modal/search", componentHandler(views.SearchDropdown()))
    e.GET("/modal/more", componentHandler(views.MoreDropdown()))
    e.GET("/modal/user", func(c echo.Context) error {
        // decide from cookie
        loggedIn := isLoggedIn(c)
        h := templ.Handler(views.UserDropdown(loggedIn))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })

    // Additional pages (full pages render their body inside #partial)
    e.GET("/about", func(c echo.Context) error {
        h := templ.Handler(views.StaticPage(isLoggedIn(c), "About", "About page (placeholder)"))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })
    e.GET("/developers", func(c echo.Context) error {
        h := templ.Handler(views.StaticPage(isLoggedIn(c), "Developers", "Developers page (placeholder)"))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })
    e.GET("/tos", func(c echo.Context) error {
        h := templ.Handler(views.StaticPage(isLoggedIn(c), "Terms of Service", "Terms of Service (placeholder)"))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })
    e.GET("/policy", func(c echo.Context) error {
        h := templ.Handler(views.StaticPage(isLoggedIn(c), "Privacy Policy", "Privacy Policy (placeholder)"))
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    })

    // Trips requires authentication
    e.GET("/trips", authWrapper(views.TripsPage("Friend"), views.UnauthPage(), views.UnauthPartial()))
}

// proxy target for users service inside cluster
func usersBaseURL() string {
    if v := os.Getenv("USERS_BASE_URL"); v != "" {
        return v
    }
    // default: same namespace DNS (Tilt deploys all into default)
    return "http://users:80"
}

// postFormUsers tries multiple base URLs for robustness in local clusters
func postFormUsers(path string, values url.Values) (*http.Response, error) {
    bases := []string{}
    if v := os.Getenv("USERS_BASE_URL"); v != "" {
        bases = append(bases, v)
    }
    // Common fallbacks
    bases = append(bases,
        "http://users:80",
        "http://users.default.svc.cluster.local:80",
        "http://users.users.svc.cluster.local:80",
    )
    var resp *http.Response
    var err error
    for _, b := range bases {
        resp, err = http.PostForm(b+path, values)
        if err == nil {
            return resp, nil
        }
    }
    return nil, err
}

func proxyUsersSignup(c echo.Context) error {
    values := url.Values{}
    values.Set("username", c.FormValue("username"))
    values.Set("password", c.FormValue("password"))
    values.Set("contact", c.FormValue("contact"))
    resp, err := postFormUsers("/v1/users/auth/signup", values)
    if err != nil {
        return c.JSON(http.StatusBadGateway, err.Error())
    }
    defer resp.Body.Close()
    // propagate only safe headers; avoid forwarding Content-Length
    if ct := resp.Header.Get("Content-Type"); ct != "" {
        c.Response().Header().Set("Content-Type", ct)
    }
    for _, v := range resp.Header.Values("Set-Cookie") {
        c.Response().Header().Add("Set-Cookie", v)
    }
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        c.Response().Header().Set("HX-Redirect", "/")
    }
    c.Response().WriteHeader(resp.StatusCode)
    _, _ = io.Copy(c.Response().Writer, resp.Body)
    return nil
}

func proxyUsersLogin(c echo.Context) error {
    values := url.Values{}
    values.Set("username", c.FormValue("username"))
    values.Set("password", c.FormValue("password"))
    resp, err := postFormUsers("/v1/users/auth/login", values)
    if err != nil {
        return c.JSON(http.StatusBadGateway, err.Error())
    }
    defer resp.Body.Close()
    if ct := resp.Header.Get("Content-Type"); ct != "" {
        c.Response().Header().Set("Content-Type", ct)
    }
    for _, v := range resp.Header.Values("Set-Cookie") {
        c.Response().Header().Add("Set-Cookie", v)
    }
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        c.Response().Header().Set("HX-Redirect", "/")
    }
    c.Response().WriteHeader(resp.StatusCode)
    _, _ = io.Copy(c.Response().Writer, resp.Body)
    return nil
}

func proxyUsersLogout(c echo.Context) error {
    req, _ := http.NewRequest(http.MethodPost, usersBaseURL()+"/v1/users/auth/logout", nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return c.JSON(http.StatusBadGateway, err.Error())
    }
    defer resp.Body.Close()
    if ct := resp.Header.Get("Content-Type"); ct != "" {
        c.Response().Header().Set("Content-Type", ct)
    }
    for _, v := range resp.Header.Values("Set-Cookie") {
        c.Response().Header().Add("Set-Cookie", v)
    }
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        c.Response().Header().Set("HX-Redirect", "/")
    }
    c.Response().WriteHeader(resp.StatusCode)
    _, _ = io.Copy(c.Response().Writer, resp.Body)
    return nil
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
        if hxRequest(c) {
            h := templ.Handler(partial)
            h.ServeHTTP(c.Response().Writer, c.Request())
            return nil
        }
        h := templ.Handler(full)
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    }
}

func hxRequest(c echo.Context) bool {
    if v := c.Request().Header.Get("HX-Request"); v != "" && v != "false" {
        return true
    }
    if v := c.Request().Header.Get("X-HX-Request"); v != "" && v != "false" {
        return true
    }
    return false
}

func isLoggedIn(c echo.Context) bool {
    // users service sets cookie name 'auth_token'
    ck, err := c.Cookie("auth_token")
    if err != nil || ck == nil {
        return false
    }
    return ck.Value != ""
}

// authWrapper returns auth UI if logged in; otherwise a login prompt page.
// For HTMX requests, returns the appropriate partial.
func authWrapper(authFull templ.Component, unauthFull templ.Component, unauthPartial templ.Component) echo.HandlerFunc {
    return func(c echo.Context) error {
        if isLoggedIn(c) {
            h := templ.Handler(authFull)
            h.ServeHTTP(c.Response().Writer, c.Request())
            return nil
        }
        // not logged in
        if hxRequest(c) {
            h := templ.Handler(unauthPartial)
            h.ServeHTTP(c.Response().Writer, c.Request())
            return nil
        }
        h := templ.Handler(unauthFull)
        h.ServeHTTP(c.Response().Writer, c.Request())
        return nil
    }
}
