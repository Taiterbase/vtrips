package views

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// UnauthenticatedBanner renders a simple banner prompting auth on protected layouts.
func UnauthenticatedBanner() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div class="p-6 bg-neutral-900 text-neutral-100 text-center border-b border-neutral-800">
  You must be logged in to access this page.
  <button hx-get="/modal/login" hx-target="#modal-portal" hx-swap="innerHTML" hx-trigger="click" class="ml-4 px-4 py-2 rounded-md bg-indigo-600 hover:bg-indigo-700">Log In</button>
</div>`)
		return nil
	})
}

// VolunCreateLayout mirrors the create layout which shows unauthenticated banner if not logged in.
func VolunCreateLayout(title string, userLoggedIn bool, page templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<!DOCTYPE html><html class="h-full bg-neutral-950"><head><title>`)
		_, _ = io.WriteString(w, templ.EscapeString(title))
		_, _ = io.WriteString(w, `</title>`)
		if err := VolunHeaderAssets().Render(ctx, w); err != nil {
			return err
		}
		_, _ = io.WriteString(w, `</head><body class="h-full overflow-hidden"><div class="flex flex-col h-full">`)
		if !userLoggedIn {
			if err := UnauthenticatedBanner().Render(ctx, w); err != nil {
				return err
			}
		} else {
			if err := VolunNavbar(true).Render(ctx, w); err != nil {
				return err
			}
		}
		if err := page.Render(ctx, w); err != nil {
			return err
		}
		_, _ = io.WriteString(w, `</div><div id="modal-portal" class="modal-portal"></div></body></html>`)
		return nil
	})
}

// VolunTripsPage is a port of pages/trips/page.html.
func VolunTripsPage(userName string) templ.Component {
	return VolunRootLayout("Trips", true, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial"
  class="flex flex-grow flex-nowrap flex-col items-center overflow-hidden relative h-full overflow-x-hidden overflow-y-scroll">
  <main
    class="h-full w-full mx-auto bg-neutral-950 relative flex-grow flex flex-col z-[1] px-8 md:px-10 lg:px-20 max-w-8xl">
    <div class="relative flex h-full w-full py-8">
      <div class="flex flex-col h-full w-full items-center">
        <div class="flex flex-col min-w-full justify-items-start items-center">
          <div class="w-full pt-4 sm:pt-8 md:pt-12 lg:pt-16">
            <div class="flex flex-row w-full justify-between items-center flex-nowrap">
              <div class="flex-nowrap whitespace-nowrap flex flex-grow">
                <h1 class="text-4xl md:text-5xl font-semibold text-neutral-100">
                  <span class="hidden md:block">
                    Welcome, `)
		_, _ = io.WriteString(w, templ.EscapeString(userName))
		_, _ = io.WriteString(w, `
                  </span>
                  <span class="block md:hidden"> Welcome </span>
                </h1>
              </div>
              <div class="flex-nowrap whitespace-nowrap flex flex-grow-0"></div>
            </div>
          </div>
          <div id="page_content" class="py-12 lg:py-16 flex flex-col w-full h-full">
            <div class="py-8 w-full h-full">
              <div class="flex flex-row justify-between items-center w-full">
                <h2 class="text-neutral-100 text-3xl md:text-4xl font-medium">
                  Trips this week
                </h2>
              </div>
              <div class="flex flex-col w-full h-full">
                <div class="flex flex-row w-full h-full flex-nowrap overflow-x-scroll justify-center items-center">
                  <div class="bg-neutral-700 h-96 w-full rounded-lg"></div>
                </div>
              </div>
            </div>
            <div class="py-8 w-full h-full">
              <div class="flex flex-row justify-between items-center w-full">
                <h2 class="text-neutral-100 text-3xl md:text-4xl font-medium">
                  Pending applications
                </h2>
              </div>
              <div class="flex flex-col w-full h-full">
                <div class="flex flex-row w-full h-full flex-nowrap overflow-x-scroll justify-center items-center">
                  <div class="bg-neutral-700 h-96 w-full rounded-lg"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </main>
</div>`)
		return nil
	}))
}

// Simple static pages for dropdown links.
func VolunAboutPage() templ.Component {
	return VolunRootLayout("About", false, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial" class="p-8 text-neutral-100">About page (placeholder)</div>`)
		return nil
	}))
}

func VolunDevelopersPage() templ.Component {
	return VolunRootLayout("Developers", false, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial" class="p-8 text-neutral-100">Developers page (placeholder)</div>`)
		return nil
	}))
}

func VolunTOSPage() templ.Component {
	return VolunRootLayout("Terms of Service", false, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial" class="p-8 text-neutral-100">Terms of Service (placeholder)</div>`)
		return nil
	}))
}

func VolunPolicyPage() templ.Component {
	return VolunRootLayout("Privacy Policy", false, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial" class="p-8 text-neutral-100">Privacy Policy (placeholder)</div>`)
		return nil
	}))
}

// VolunStaticContent returns only the inner partial content for HTMX swaps.
func VolunStaticContent(title, body string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial" class="p-8 text-neutral-100">`)
		_, _ = io.WriteString(w, templ.EscapeString(body))
		_, _ = io.WriteString(w, `</div>`)
		return nil
	})
}
