package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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

	e.GET("/modal/login", componentHandler(views.LoginModal()))
	e.GET("/modal/sign-up", componentHandler(views.SignupModal()))
	e.GET("/modal/search", componentHandler(views.SearchDropdown()))
	e.GET("/modal/more", componentHandler(views.MoreDropdown()))
	e.GET("/modal/user", func(c echo.Context) error {
		// decide from cookie
		loggedIn := isLoggedIn(c)
		h := templ.Handler(views.UserDropdown(loggedIn))
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})

	e.POST("/auth/sign-up", proxyUsersSignup)
	e.POST("/auth/login", proxyUsersLogin)
	e.PUT("/auth/login", proxyUsersLogin) // support hx-put from modal
	e.POST("/auth/logout", proxyUsersLogout)
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
	e.GET("/trips", authWrapper(tripsIndexHandler, views.UnauthPage(), views.UnauthPartial()))
	e.GET("/trips/new", tripsNewHandler)
	e.POST("/trips", tripsCreateHandler)
	e.PUT("/trips/:trip_id", tripsUpdateHandler)
	e.GET("/trips/:trip_id", tripsShowHandler)
}

func tripsIndexHandler(c echo.Context) error {
	if !isLoggedIn(c) {
		return c.Redirect(http.StatusFound, "/")
	}
	orgID := currentOrgID(c)
	summaries, err := fetchTripSummaries(orgID)
	if err != nil {
		return c.JSON(http.StatusBadGateway, err.Error())
	}
	data := views.TripsWizardData{Step: views.TripWizardStepBasics, OrgID: orgID, Form: map[string]string{}}
	return renderTripsPage(c, "Friend", views.TripsWizardBasics(data), views.TripsSummaryList(summaries))
}

func tripsNewHandler(c echo.Context) error {
	if !isLoggedIn(c) {
		return c.Redirect(http.StatusFound, "/")
	}
	orgID := currentOrgID(c)
	data := views.TripsWizardData{Step: views.TripWizardStepBasics, OrgID: orgID, Form: map[string]string{}}
	summaries, err := fetchTripSummaries(orgID)
	if err != nil {
		summaries = nil
	}
	wizard := views.TripsWizardBasics(data)
	summary := views.TripsSummaryList(summaries)
	if isHXRequest(c) {
		return renderWizardPartial(c, wizard, summary)
	}
	return renderTripsPage(c, "Friend", wizard, summary)
}

type tripsPayload map[string]string

func fetchTripSummaries(orgID string) ([]views.TripSummary, error) {
	req, err := http.NewRequest(http.MethodGet, tripsBaseURL()+"/v1/trips?org_id="+url.QueryEscape(orgID), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("trips service error: %s", strings.TrimSpace(string(body)))
	}

	return parseTripSummaries(body)
}

func parseTripSummaries(body []byte) ([]views.TripSummary, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return nil, nil
	}

	if trimmed[0] == '[' {
		var list []map[string]any
		if err := json.Unmarshal(trimmed, &list); err == nil {
			return buildTripSummaries(list), nil
		}
		var generic []any
		if err := json.Unmarshal(trimmed, &generic); err != nil {
			return nil, err
		}
		return buildTripSummariesFromAny(generic), nil
	}

	var payload map[string]any
	if err := json.Unmarshal(trimmed, &payload); err != nil {
		return nil, err
	}

	if tripsVal, ok := payload["trips"]; ok {
		return buildTripSummariesFromInterface(tripsVal)
	}

	if dataVal, ok := payload["data"]; ok {
		switch d := dataVal.(type) {
		case []any:
			return buildTripSummariesFromAny(d), nil
		case map[string]any:
			if inner, ok := d["trips"]; ok {
				return buildTripSummariesFromInterface(inner)
			}
		}
	}

	return nil, nil
}

func buildTripSummaries(list []map[string]any) []views.TripSummary {
	summaries := make([]views.TripSummary, 0, len(list))
	for _, entry := range list {
		summaries = append(summaries, views.NewTripSummaryFromPayload(entry))
	}
	return summaries
}

func buildTripSummariesFromAny(raw []any) []views.TripSummary {
	summaries := make([]views.TripSummary, 0, len(raw))
	for _, entry := range raw {
		if tripMap, ok := entry.(map[string]any); ok {
			summaries = append(summaries, views.NewTripSummaryFromPayload(tripMap))
		}
	}
	return summaries
}

func buildTripSummariesFromInterface(value any) ([]views.TripSummary, error) {
	switch v := value.(type) {
	case []any:
		return buildTripSummariesFromAny(v), nil
	case []map[string]any:
		return buildTripSummaries(v), nil
	default:
		return nil, nil
	}
}

func renderTripsPage(c echo.Context, userName string, wizard templ.Component, summary templ.Component) error {
	templ.Handler(views.TripsPage(userName, wizard, summary)).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func renderWizardPartial(c echo.Context, wizard templ.Component, summary templ.Component) error {
	templ.Handler(views.TripsDashboard(wizard, summary)).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func isHXRequest(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true" || c.Request().Header.Get("X-HX-Request") == "true"
}

func tripsCreateHandler(c echo.Context) error {
	if !isLoggedIn(c) {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	orgID := currentOrgID(c)
	payload := make(tripsPayload)
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	payload["org_id"] = orgID
	payload["status"] = defaultTripStatus(payload["status"])
	payload["step"] = string(views.TripWizardStepBasics)

	resp, err := postJSON(tripsBaseURL()+"/v1/trips", payload)
	if err != nil {
		return c.JSON(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return c.JSON(resp.StatusCode, string(body))
	}

	var trip map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&trip); err != nil {
		return c.JSON(http.StatusBadGateway, err.Error())
	}

	tripID, _ := trip["id"].(string)
	data := views.TripsWizardData{
		Step:   views.TripWizardStepLogistics,
		TripID: tripID,
		OrgID:  orgID,
		Form:   views.TripFormFromPayload(trip),
	}

	summaries, err := fetchTripSummaries(orgID)
	if err != nil {
		summaries = nil
	}
	wizard := views.TripsWizardLogistics(data)
	summary := views.TripsSummaryList(summaries)
	if isHXRequest(c) {
		return renderWizardPartial(c, wizard, summary)
	}
	return renderTripsPage(c, "Friend", wizard, summary)
}

func tripsUpdateHandler(c echo.Context) error {
	if !isLoggedIn(c) {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	tripID := c.Param("trip_id")
	payload := make(tripsPayload)
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	payload["org_id"] = currentOrgID(c)

	resp, err := putJSON(tripsBaseURL()+"/v1/trips/"+tripID+"?org_id="+url.QueryEscape(payload["org_id"]), payload)
	if err != nil {
		return c.JSON(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return c.JSON(resp.StatusCode, string(body))
	}

	summaries, err := fetchTripSummaries(payload["org_id"])
	if err != nil {
		summaries = nil
	}
	summary := views.NewTripSummaryFromPayload(map[string]any{
		"id":      tripID,
		"name":    payload["name"],
		"city":    payload["city"],
		"country": payload["country"],
		"status":  payload["status"],
	})

	data := views.TripsWizardData{Step: views.TripWizardStepReview, TripID: tripID, OrgID: payload["org_id"], Form: payload}
	wizard := views.TripsWizardReview(data, summary)
	summaryList := views.TripsSummaryList(summaries)
	if isHXRequest(c) {
		return renderWizardPartial(c, wizard, summaryList)
	}
	return renderTripsPage(c, "Friend", wizard, summaryList)
}

func tripsShowHandler(c echo.Context) error {
	if !isLoggedIn(c) {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	tripID := c.Param("trip_id")
	orgID := currentOrgID(c)
	req, err := http.NewRequest(http.MethodGet, tripsBaseURL()+"/v1/trips/"+tripID+"?org_id="+url.QueryEscape(orgID), nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadGateway, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return c.JSON(resp.StatusCode, string(body))
	}

	var trip map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&trip); err != nil {
		return c.JSON(http.StatusBadGateway, err.Error())
	}

	summary := views.NewTripSummaryFromPayload(trip)
	data := views.TripsWizardData{Step: views.TripWizardStepReview, TripID: tripID, OrgID: orgID, Form: views.TripFormFromPayload(trip)}
	wizard := views.TripsWizardReview(data, summary)
	summaryList := views.TripsSummaryList([]views.TripSummary{summary})
	if isHXRequest(c) {
		return renderWizardPartial(c, wizard, summaryList)
	}
	return renderTripsPage(c, "Friend", wizard, summaryList)
}

func tripsBaseURL() string {
	if v := os.Getenv("TRIPS_BASE_URL"); v != "" {
		return v
	}
	return "http://trips:80"
}

func currentOrgID(c echo.Context) string {
	if v := c.QueryParam("org_id"); v != "" {
		return v
	}
	if ck, err := c.Cookie("org_id"); err == nil && ck != nil {
		return ck.Value
	}
	return "org-default"
}

func defaultTripStatus(requested string) string {
	switch requested {
	case "listed", "complete", "unlisted":
		return requested
	default:
		return "draft"
	}
}

func postJSON(url string, payload any) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func putJSON(url string, payload any) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
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
func authWrapper(authHandler echo.HandlerFunc, unauthFull templ.Component, unauthPartial templ.Component) echo.HandlerFunc {
	return func(c echo.Context) error {
		if isLoggedIn(c) {
			return authHandler(c)
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
