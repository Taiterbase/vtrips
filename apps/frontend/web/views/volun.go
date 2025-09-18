package views

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// RootNavbar is a simplified port of the Voluntrips navbar.
func RootNavbar() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<nav class=\"z-40 flex-shrink-0\"><div class=\"flex flex-nowrap overflow-hidden items-stretch h-24 bg-neutral-850 shadow-black shadow-sm\">")
		// Brand
		_, _ = io.WriteString(w, "<a href=\"/\" class=\"decoration-[none] flex-shrink-0\"><div class=\"p-2 px-6 inline-flex h-full items-center\"><div class=\"inline-flex\"><img class=\"h-16 w-auto\" src=\"/img/avatar.svg\" alt=\"logo\"></div></div></a>")
		// Primary links
		_, _ = io.WriteString(w, "<div><div class=\"flex justify-between flex-row h-full\">")
		_, _ = io.WriteString(w, "<div class=\"flex-col h-full px-2 md:px-4 lg:px-6 flex\"><div class=\"flex self-center h-full\"><a id=\"nav-item\" href=\"/browse\" class=\"text-xl lg:text-2xl px-2 md:px-0 mt-2 border-b-4 border-b-transparent hover:text-indigo-500 font-bold items-center text-center whitespace-nowrap text-neutral-100\">Browse</a></div></div>")
		_, _ = io.WriteString(w, "</div></div>")
		// Right controls (placeholder)
		_, _ = io.WriteString(w, "<div class=\"flex h-full items-center px-6 mr-4\"></div>")
		_, _ = io.WriteString(w, "</div></nav>")
		return nil
	})
}

// RootLayout replicates dark theme layout with navbar and content slot.
func RootLayout(title string, page templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<!DOCTYPE html><html class=\"h-full bg-neutral-950\"><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
		_, _ = io.WriteString(w, "<title>")
		_, _ = io.WriteString(w, templ.EscapeString(title))
		_, _ = io.WriteString(w, "</title>")
		// Tailwind CDN + optional custom CSS
		_, _ = io.WriteString(w, "<script src=\"https://cdn.tailwindcss.com\"></script>")
		_, _ = io.WriteString(w, "<link rel=\"stylesheet\" href=\"/css/custom.css\">")
		// HTMX CDN for interactive porting
		_, _ = io.WriteString(w, "<script src=\"https://unpkg.com/htmx.org@1.9.12\"></script>")
		_, _ = io.WriteString(w, "</head><body class=\"h-full overflow-hidden\"><div class=\"flex flex-col h-full\">")
		if err := RootNavbar().Render(ctx, w); err != nil {
			return err
		}
		if err := page.Render(ctx, w); err != nil {
			return err
		}
		_, _ = io.WriteString(w, "</div></body></html>")
		return nil
	})
}

// HomePorted is a minimal port of Voluntrips home.
func HomePorted() templ.Component {
	return RootLayout("Home", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<div id=\"partial\" class=\"flex flex-grow flex-nowrap flex-col overflow-hidden relative h-full\">")
		_, _ = io.WriteString(w, "<main class=\"h-full w-full bg-neutral-950 relative flex-grow flex flex-col z-[1] max-w-8xl mx-auto\">")
		_, _ = io.WriteString(w, "<div class=\"relative flex h-full w-full py-8\"><div class=\"flex flex-col h-full w-full items-center\"><div class=\"flex flex-col min-w-full overflow-x-hidden overflow-y-scroll\">")
		_, _ = io.WriteString(w, "<div class=\"flex flex-grow min-h-[40rem] max-h-full w-full mt-8\"><div class=\"grid flex-grow grid-cols-2 md:grid-cols-3 xl:grid-cols-4 grid-flow-row w-full\">")
		for i := 0; i < 8; i++ {
			_, _ = io.WriteString(w, "<div class=\"w-full h-full min-h-[40rem]\"><div class=\"p-2 w-full h-full\"><div class=\"bg-neutral-700 w-full h-full rounded-md\"></div></div></div>")
		}
		_, _ = io.WriteString(w, "</div></div></div></div></div></main></div>")
		return nil
	}))
}

// BrowsePorted placeholder; can expand with filters later.
func BrowsePorted() templ.Component {
	return RootLayout("Browse", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, "<div class=\"p-6 text-neutral-100\">Browse trips (coming soon)</div>")
		return nil
	}))
}
