package views

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// Layout composes a basic HTML page with optional Tailwind via CDN and local CSS.
func Layout(title string, body templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<!doctype html><html lang=\"en\"><head>")
		_, _ = io.WriteString(w, "<meta charset=\"utf-8\">")
		_, _ = io.WriteString(w, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
		_, _ = io.WriteString(w, "<title>")
		_, _ = io.WriteString(w, templ.EscapeString(title))
		_, _ = io.WriteString(w, "</title>")
		// Tailwind via CDN for development; replace with compiled CSS for production.
		_, _ = io.WriteString(w, "<script src=\"https://cdn.tailwindcss.com\"></script>")
		_, _ = io.WriteString(w, "<link rel=\"stylesheet\" href=\"/static/styles.css\">")
		_, _ = io.WriteString(w, "</head><body class=\"min-h-screen bg-gray-50 text-gray-900\">")
		// Simple nav
		_, _ = io.WriteString(w, "<header class=\"bg-white border-b\"><div class=\"max-w-7xl mx-auto px-4 py-4 flex items-center justify-between\">")
		_, _ = io.WriteString(w, "<a href=\"/\" class=\"font-semibold text-lg\">vTrips</a>")
		_, _ = io.WriteString(w, "<nav class=\"space-x-4\"><a class=\"text-gray-600 hover:text-gray-900\" href=\"/org/dashboard\">Dashboard</a><a class=\"text-gray-600 hover:text-gray-900\" href=\"/org/trips\">Trips</a><a class=\"text-gray-600 hover:text-gray-900\" href=\"/org/applications\">Applications</a></nav>")
		_, _ = io.WriteString(w, "</div></header>")
		_, _ = io.WriteString(w, "<main class=\"max-w-7xl mx-auto px-4 py-8\">")
		if err := body.Render(ctx, w); err != nil {
			return err
		}
		_, _ = io.WriteString(w, "</main>")
		_, _ = io.WriteString(w, "</body></html>")
		return nil
	})
}

func HomePage() templ.Component {
	return Layout("vTrips", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<section class=\"text-center space-y-4\">")
		_, _ = io.WriteString(w, "<h1 class=\"text-3xl font-bold\">Volunteer trips made simple</h1>")
		_, _ = io.WriteString(w, "<p class=\"text-gray-600\">Organize and manage service trips worldwide.</p>")
		_, _ = io.WriteString(w, "<div class=\"pt-4\"><a href=\"/org/dashboard\" class=\"inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700\">Go to Dashboard</a></div>")
		_, _ = io.WriteString(w, "</section>")
		return nil
	}))
}

func OrgDashboardPage() templ.Component {
	return Layout("Organization Dashboard", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<div class=\"grid grid-cols-1 md:grid-cols-3 gap-6\">")
		card := func(title, value, href string) templ.Component {
			return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
				_, _ = io.WriteString(w, "<a href=\"")
				_, _ = io.WriteString(w, templ.EscapeString(href))
				_, _ = io.WriteString(w, "\" class=\"block p-6 bg-white rounded-lg shadow hover:shadow-md transition\">")
				_, _ = io.WriteString(w, "<div class=\"text-sm text-gray-500\">")
				_, _ = io.WriteString(w, templ.EscapeString(title))
				_, _ = io.WriteString(w, "</div><div class=\"mt-2 text-2xl font-semibold\">")
				_, _ = io.WriteString(w, templ.EscapeString(value))
				_, _ = io.WriteString(w, "</div></a>")
				return nil
			})
		}
		_ = card("Active Trips", "3", "/org/trips").Render(ctx, w)
		_ = card("Pending Applications", "12", "/org/applications").Render(ctx, w)
		_ = card("Create a Trip", "+", "/org/trips/new").Render(ctx, w)
		_, _ = io.WriteString(w, "</div>")
		return nil
	}))
}

func TripsIndexPage() templ.Component {
	return Layout("Trips", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<div class=\"flex items-center justify-between\">")
		_, _ = io.WriteString(w, "<h2 class=\"text-xl font-semibold\">Trips</h2>")
		_, _ = io.WriteString(w, "<a href=\"/org/trips/new\" class=\"px-3 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700\">New Trip</a>")
		_, _ = io.WriteString(w, "</div>")
		_, _ = io.WriteString(w, "<div class=\"mt-6 text-gray-500\">No trips to show yet.</div>")
		return nil
	}))
}

func TripNewPage() templ.Component {
	return Layout("Create Trip", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<form method=\"post\" action=\"/org/trips\" class=\"space-y-4 max-w-xl\">")
		_, _ = io.WriteString(w, "<div><label class=\"block text-sm text-gray-600\">Title</label><input name=\"title\" class=\"mt-1 w-full border rounded-md px-3 py-2\" required></div>")
		_, _ = io.WriteString(w, "<div><label class=\"block text-sm text-gray-600\">Location</label><input name=\"location\" class=\"mt-1 w-full border rounded-md px-3 py-2\" required></div>")
		_, _ = io.WriteString(w, "<div><label class=\"block text-sm text-gray-600\">Description</label><textarea name=\"description\" class=\"mt-1 w-full border rounded-md px-3 py-2\" rows=\"5\"></textarea></div>")
		_, _ = io.WriteString(w, "<div class=\"pt-2\"><button class=\"px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700\" type=\"submit\">Create</button></div>")
		_, _ = io.WriteString(w, "</form>")
		return nil
	}))
}

func ApplicationsIndexPage() templ.Component {
	return Layout("Applications", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<h2 class=\"text-xl font-semibold\">Volunteer Applications</h2>")
		_, _ = io.WriteString(w, "<div class=\"mt-4 text-gray-500\">No applications yet.</div>")
		return nil
	}))
}

func ApplicationShowPage() templ.Component {
	return Layout("Application", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<div class=\"space-y-4\">")
		_, _ = io.WriteString(w, "<h2 class=\"text-xl font-semibold\">Application</h2>")
		_, _ = io.WriteString(w, "<form method=\"post\" action=\"grade\" class=\"space-x-2\">")
		_, _ = io.WriteString(w, "<button name=\"decision\" value=\"approve\" class=\"px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700\">Approve</button>")
		_, _ = io.WriteString(w, "<button name=\"decision\" value=\"reject\" class=\"px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700\">Reject</button>")
		_, _ = io.WriteString(w, "</form>")
		_, _ = io.WriteString(w, "</div>")
		return nil
	}))
}
