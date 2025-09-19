package views

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// VolunHeaderAssets renders the shared head assets (CSS and JS).
func VolunHeaderAssets() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<link href="/css/styles.css" rel="stylesheet" />`)
		_, _ = io.WriteString(w, `<script src="/js/htmx.min.js"></script>`)
		return nil
	})
}

// VolunNavbar ports the Voluntrips root navbar.
func VolunNavbar(userLoggedIn bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<nav class="z-40 flex-shrink-0">
  <div class="flex flex-nowrap overflow-hidden items-stretch h-24 bg-neutral-850 shadow-black shadow-sm">
    <div class="flex-grow flex-shrink-[2] w-full flex items-stretch justify-start flex-nowrap">
      <a href="/" class="decoration-[none] flex-shrink-0">
        <div class="p-2 px-6 inline-flex h-full items-center">
          <div class="inline-flex">
            <img class="h-16 w-auto" src="/img/logo.png" alt="logo" />
          </div>
        </div>
      </a>
      <div>
        <div class="flex justify-between flex-row h-full">`)
		if userLoggedIn {
			_, _ = io.WriteString(w, `<div class="flex-col h-full px-2 md:px-4 lg:px-6 flex">
            <div class="flex self-center h-full">
              <button id="nav-item" hx-get="/likes" hx-swap="outerHTML" hx-target="#partial" hx-trigger="click"
                hx-push-url="true"
                class="text-xl lg:text-2xl px-2 md:px-0 mt-2 border-b-4 border-b-transparent hover:text-indigo-500 font-bold items-center text-center whitespace-nowrap text-neutral-100">
                <span class="hidden md:block"> Likes </span>
                <span class="block md:hidden">
                  <svg type="color-fill-current" width="20px" height="20px" viewBox="0 0 20 20" class="fill-current">
                    <g>
                      <path fill-rule="evenodd"
                        d="M9.171 4.171A4 4 0 006.343 3H6a4 4 0 00-4 4v.343a4 4 0 001.172 2.829L10 17l6.828-6.828A4 4 0 0018 7.343V7a4 4 0 00-4-4h-.343a4 4 0 00-2.829 1.172L10 5l-.829-.829zm.829 10l5.414-5.414A2 2 0 0016 7.343V7a2 2 0 00-2-2h-.343a2 2 0 00-1.414.586L10 7.828 7.757 5.586A2 2 0 006.343 5H6a2 2 0 00-2 2v.343a2 2 0 00.586 1.414L10 14.172z"
                        clip-rule="evenodd"></path>
                    </g>
                  </svg>
                </span>
              </button>
            </div>
          </div>`)
		}
		_, _ = io.WriteString(w, `<div class="flex-col h-full px-2 md:px-4 lg:px-6 flex">
            <div class="flex self-center h-full">
              <button id="nav-item" hx-get="/browse" hx-swap="outerHTML" hx-target="#partial" hx-trigger="click"
                hx-push-url="true"
                class="text-xl lg:text-2xl px-2 md:px-0 mt-2 border-b-4 border-b-transparent hover:text-indigo-500 font-bold items-center text-center whitespace-nowrap text-neutral-100">
                <span class="hidden md:block"> Browse </span>
                <span class="block md:hidden">
                  <svg type="color-fill-current" width="20px" height="20px" viewBox="0 0 20 20" class="fill-current">
                    <g>
                      <path d="M5 2a2 2 0 00-2 2v8a2 2 0 002 2V4h8a2 2 0 00-2-2H5z"></path>
                      <path fill-rule="evenodd"
                        d="M7 8a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H9a2 2 0 01-2-2V8zm2 0h6v8H9V8z"
                        clip-rule="evenodd"></path>
                    </g>
                  </svg>
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
      <div class="flex h-full items-center px-6 mr-4">
        <div class="inline-block relative">
          <div class="relative">
            <button
              class="inline-flex items-center justify-center select-none rounded-md h-10 w-10 bg-transparent text-neutral-100 relative align-middle overflow-hidden decoration-[none] whitespace-nowrap font-semibold text-xl"
              aria-label="More" id="more-drop-down-button">
              <div class="pointer-events-none w-10 h-10">
                <div class="inline-flex items-center w-full h-full fill-current">
                  <svg width="100%" height="100%" viewBox="0 0 20 20" focusable="false" aria-hidden="true">
                    <path d="M10 18a2 2 0 1 1 0-4 2 2 0 0 1 0 4zm0-6a2 2 0 1 1 0-4 2 2 0 0 1 0 4zM8 4a2 2 0 1 0 4 0 2 2 0 0 0-4 0z"></path>
                  </svg>
                </div>
              </div>
            </button>
            <div id="more-drop-down-portal" class="fixed z-50"></div>
          </div>
        </div>
      </div>
    </div>
    <div class="flex flex-grow-[2] flex-shrink w-full items-center justify-center">
      <div class="relative basis-[40rem] mx-8">
        <div class="max-w-[40rem] relative z-10 sm:block hidden">
          <div class="p-2">
            <div class="flex w-full">
              <div class="flex-grow">
                <div class="transform-gpu">
                  <div class="relative">
                    <input autofocus id="search-input" aria-label="Search Input" aria-haspopup="grid" type="search"
                      class="h-14 [-webkit-filter: invert(100%)] border-neutral-500 text-neutral-100 placeholder-neutral-200 bg-transparent text-xl p-4 pr-2 flex w-full rounded-l-lg hover:ring-2 hover:ring-neutral-500 hover:ring-inset hover:ring-offset-neutral-500 focus:ring-2 focus:ring-inset focus:outline-none focus:border-indigo-600 focus:ring-indigo-600 focus:ring-offset-indigo-600"
                      placeholder="Search" autocapitalize="off" autocorrect="off" maxlength="70" spellcheck="false" />
                  </div>
                </div>
              </div>
              <button
                class="bg-neutral-700 flex-shrink-0 p-2 rounded-r-md inline-flex relative items-center align-middle justify-center overflow-hidden whitespace-nowrap select-none">
                <div class="flex items-center">
                  <div class="inline-flex items-center w-10">
                    <div class="inline-flex items-center h-full w-full">
                      <div class="relative w-full overflow-hidden">
                        <svg fill="#fff" width="100%" height="100%" viewBox="0 0 20 20" class="ScIconSVG-sc-1q25cff-1 jpczqG">
                          <g>
                            <path fill-rule="evenodd"
                              d="M13.192 14.606a7 7 0 111.414-1.414l3.101 3.1-1.414 1.415-3.1-3.1zM14 9A5 5 0 114 9a5 5 0 0110 0z"
                              clip-rule="evenodd"></path>
                          </g>
                        </svg>
                      </div>
                    </div>
                  </div>
                </div>
              </button>
            </div>
          </div>
        </div>
        <!-- <div class="top-0 bg-neutral-950 rounded-md block shadow-md w-full p-2 absolute"></div> -->
      </div>
    </div>
    <div class="flex flex-grow flex-shrink-[2] w-full items-center justify-end">
      <div class="flex-nowrap flex sm:hidden">
        <div class="flex px-2 md:px-4 lg:px-6 ">
          <button id="search-drop-down-button"
            class="hover:bg-neutral-700 p-2 flex-shrink-0 rounded-md inline-flex relative items-center align-middle justify-center overflow-hidden whitespace-nowrap select-none">
            <div class="flex items-center">
              <div class="inline-flex items-center w-8">
                <div class="inline-flex items-center h-full w-full">
                  <div class="relative w-full overflow-hidden">
                    <svg fill="#fff" width="100%" height="100%" viewBox="0 0 20 20" class="ScIconSVG-sc-1q25cff-1 jpczqG">
                      <g>
                        <path fill-rule="evenodd"
                          d="M13.192 14.606a7 7 0 111.414-1.414l3.101 3.1-1.414 1.415-3.1-3.1zM14 9A5 5 0 114 9a5 5 0 0110 0z"
                          clip-rule="evenodd"></path>
                      </g>
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </button>
        </div>
        <div id="search-drop-down-portal" class="fixed z-50"></div>
      </div>`)
		if !userLoggedIn {
			_, _ = io.WriteString(w, `<div class="py-4 h-full flex-row sm:flex hidden">
        <div class="flex-nowrap flex">
          <div class="flex flex-col flex-nowrap h-full px-2 md:px-4 items-center justify-center">
            <button hx-get="/modal/login" hx-target="#modal-portal" hx-swap="innerHTML" hx-trigger="click"
              hx-push-url="false"
              class="bg-neutral-500 hover:bg-neutral-600 whitespace-nowrap text-neutral-100 py-2 text-xl lg:text-2xl font-semibold rounded-md ">
              <div class="flex-col h-full px-4 flex">
                <div class="flex self-center h-full">
                  <span class="hidden md:block"> Log In </span>
                  <span class="block md:hidden"> Log In </span>
                </div>
              </div>
            </button>
          </div>
        </div>
        <div class="flex-nowrap flex">
          <div class="flex flex-col flex-nowrap h-full px-2 md:px-4 items-center justify-center">
            <button hx-get="/modal/sign-up" hx-target="#modal-portal" hx-swap="innerHTML" hx-trigger="click"
              hx-push-url="false"
              class="bg-indigo-500 hover:bg-indigo-600 whitespace-nowrap text-neutral-100 py-2 text-xl lg:text-2xl font-semibold rounded-md ">
              <div class="flex-col h-full px-4 flex">
                <div class="flex self-center h-full">
                  <span class="hidden md:block"> Sign Up </span>
                  <span class="block md:hidden"> Sign Up </span>
                </div>
              </div>
            </button>
          </div>
        </div>
      </div>`)
		} else {
			_, _ = io.WriteString(w, `<div class="flex flex-col flex-nowrap h-full px-2 md:px-4 items-center justify-center">
        <a href="/trips"
          class="bg-indigo-500 hover:bg-indigo-600 whitespace-nowrap text-neutral-100 py-2 text-xl lg:text-2xl font-semibold rounded-md ">
          <div class="flex-col h-full px-4 flex">
            <div class="flex self-center h-full">
              <span class="hidden md:block"> Manage Trips </span>
              <span class="block md:hidden"> Manage Trips </span>
            </div>
          </div>
        </a>
      </div>`)
		}
		_, _ = io.WriteString(w, `<div class="h-full mr-4 flex px-2">
        <div class="flex-nowrap flex h-full">
          <div class="pl-2 relative flex flex-grow items-stretch h-full">
            <div class="inline-block relative">
              <button id="user-drop-down-button" class="bg-none h-full border-none border-r-0 text-inherit">
                <div class="relative max-h-full w-16 h-16">
                  <img class="p-2 w-full h-full block border-none max-w-full align-top" src="/img/avatar.svg"
                    alt="profile avatar" />
                </div>
              </button>
              <div id="user-drop-down-portal" class="fixed z-50"></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</nav>
<script>
  document
    .getElementById("user-drop-down-button")
    .addEventListener("click", function () {
      if (!document.getElementById("user-drop-down-overlay")) {
        htmx.ajax("GET", "/modal/user?layout=root", {
          target: "#user-drop-down-portal",
          swap: "innerHTML",
          pushUrl: "false",
        });
      }
    });

  document
    .getElementById("more-drop-down-button")
    .addEventListener("click", function () {
      if (!document.getElementById("more-drop-down-overlay")) {
        htmx.ajax("GET", "/modal/more", {
          target: "#more-drop-down-portal",
          swap: "innerHTML",
          pushUrl: "false",
        });
      }
    });

  document
    .getElementById("search-drop-down-button")
    .addEventListener("click", function () {
      if (!document.getElementById("search-drop-down-overlay")) {
        htmx.ajax("GET", "/modal/search", {
          target: "#search-drop-down-portal",
          swap: "innerHTML",
          pushUrl: "false",
        });
      }
    });

  function updateActiveNavItem() {
    const currentPath = window.location.pathname;
    document.querySelectorAll("#nav-item").forEach((item) => {
      if (item.getAttribute("hx-get") === currentPath) {
        item.classList.add("text-indigo-400", "border-b-indigo-400");
        item.classList.remove("text-neutral-100");
        item.classList.remove("border-b-transparent");
      } else {
        item.classList.remove("text-indigo-400", "border-b-indigo-400");
        item.classList.add("text-neutral-100");
        item.classList.add("border-b-transparent");
      }
    });
  }

  updateActiveNavItem();
  document.body.addEventListener("htmx:afterSwap", updateActiveNavItem);
</script>`)
		return nil
	})
}

// VolunRootLayout composes the page with header assets, navbar, content, and modal portal.
func VolunRootLayout(title string, userLoggedIn bool, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<!DOCTYPE html><html class="h-full bg-neutral-950"><head><title>`)
		_, _ = io.WriteString(w, templ.EscapeString(title))
		_, _ = io.WriteString(w, `</title>`)
		if err := VolunHeaderAssets().Render(ctx, w); err != nil {
			return err
		}
		_, _ = io.WriteString(w, `</head><body class="h-full overflow-hidden"><div class="flex flex-col h-full">`)
		if err := VolunNavbar(userLoggedIn).Render(ctx, w); err != nil {
			return err
		}
		if err := content.Render(ctx, w); err != nil {
			return err
		}
		_, _ = io.WriteString(w, `</div><div id="modal-portal" class="modal-portal"></div></body></html>`)
		return nil
	})
}

// VolunHomeContent renders the Home page content.
func VolunHomeContent() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial" class="flex flex-grow flex-nowrap flex-col overflow-hidden relative h-full">
  <main class="h-full w-full bg-neutral-950 relative flex-grow flex flex-col z-[1] max-w-8xl mx-auto">
    <div class="relative flex h-full w-full py-8">
      <div class="flex flex-col h-full w-full items-center">
        <div class="flex flex-col min-w-full overflow-x-hidden overflow-y-scroll">
          <div class="flex flex-col px-4 sm:px-4 md:px-6 lg:px-8 flex-grow-0 flex-nowrap">
            <div class="flex-nowrap whitespace-nowrap w-[30rem]">
              <h1 class="text-6xl font-bold text-neutral-100">Home</h1>
              <hr class="w-full my-2" />
            </div>
          </div>
          <div class="flex flex-grow min-h-[40rem] max-h-full w-full mt-8">
            <div class="grid flex-grow grid-cols-2 md:grid-cols-3 xl:grid-cols-4 grid-flow-row w-full">
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
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
	})
}

// VolunBrowseContent renders the Browse page content.
func VolunBrowseContent() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial" class="flex flex-grow flex-nowrap flex-col overflow-x-hidden overflow-y-scroll relative h-full">
  <main class="h-full w-full bg-neutral-950 relative flex-grow flex flex-col z-[1] max-w-8xl mx-auto">
    <div class="relative flex h-full w-full py-8">
      <div class="flex flex-col h-full w-full items-center">
        <div class="flex flex-col min-w-full ">
          <div class="flex flex-col px-4 sm:px-4 md:px-6 lg:px-8 flex-grow-0 flex-nowrap">
            <div class="flex-nowrap whitespace-nowrap w-[30rem]">
              <h1 class="text-6xl font-bold text-neutral-100">Browse</h1>
              <hr class="w-full my-2" />
            </div>
          </div>
          <div class="flex flex-grow min-h-[40rem] max-h-full w-full mt-8">
            <div class="grid flex-grow grid-cols-2 md:grid-cols-3 xl:grid-cols-4 grid-flow-row w-full">
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
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
	})
}

// VolunLikesContent renders the Likes page content.
func VolunLikesContent() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="partial" class="flex flex-grow flex-nowrap flex-col overflow-hidden relative h-full">
  <main class="h-full w-full bg-neutral-950 relative flex-grow flex flex-col z-[1] max-w-8xl mx-auto">
    <div class="relative flex h-full w-full py-8">
      <div class="flex flex-col h-full w-full items-center">
        <div class="flex flex-col min-w-full overflow-x-hidden overflow-y-scroll">
          <div class="flex flex-col px-4 sm:px-4 md:px-6 lg:px-8 flex-grow-0 flex-nowrap">
            <div class="flex-nowrap whitespace-nowrap w-[30rem]">
              <h1 class="text-6xl font-bold text-neutral-100">Likes</h1>
              <hr class="w-full my-2" />
            </div>
          </div>
          <div class="flex flex-grow min-h-[40rem] max-h-full w-full mt-8">
            <div class="grid flex-grow grid-cols-2 md:grid-cols-3 xl:grid-cols-4 grid-flow-row w-full">
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
                </div>
              </div>
              <div class="w-full h-full min-h-[40rem]">
                <div class="p-2 w-full h-full">
                  <div class="bg-neutral-700 w-full h-full rounded-md"></div>
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
	})
}

// Page wrappers using the Volun layout and navbar.
func VolunHomePage(userLoggedIn bool) templ.Component {
	return VolunRootLayout("Home", userLoggedIn, VolunHomeContent())
}

func VolunBrowsePage(userLoggedIn bool) templ.Component {
	return VolunRootLayout("Browse", userLoggedIn, VolunBrowseContent())
}

func VolunLikesPage(userLoggedIn bool) templ.Component {
	return VolunRootLayout("Likes", userLoggedIn, VolunLikesContent())
}

// Modals and dropdowns
func ModalLogin() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, ``) // filled below by Login modal content
		_, _ = io.WriteString(w, `<div id="modal-overlay" onclick="if (event.target.id === this.id) this?.remove();"
  class="bg-opacity-80 bg-black text-neutral-100 items-start flex fixed bottom-0 left-0 right-0 top-0 z-50 overflow-auto justify-center">
  <div id="modal-content" class="flex flex-grow-0 h-full w-full outline-none justify-center pointer-events-none">
    <div class="p-4 pointer-events-none relative flex w-full h-full items-start justify-center">
      <div class="my-auto outline-none flex flex-grow-0 justify-center relative w-full pointer-events-none">
        <div class="pointer-events-auto max-w-full relative block shadow-md">
          <div id="modal-form" class="flex rounded-md overflow-hidden">
            <div class="w-[50rem] overflow-auto block">
              <div class="flex flex-col bg-neutral-850 p-12">
                <div class="box-border block">
                  <div class="flex flex-col mt-2">
                    <div class="inline-flex items-center justify-center">
                      <div class="inline-flex items-center fill-slate-100 flex-shrink-0">
                        <img class="w-auto h-[50px]" src="/img/logo.png" alt="Voluntrips logo" />
                      </div>
                      <div class="ml-4 text-center">
                        <h4 class="text-2xl font-semibold">
                          Log in to Voluntrips
                        </h4>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="mb-4">
                  <form novalidate>
                    <div class="w-full flex-col flex">
                      <div id="username" class="mt-8">
                        <div class="block align-baseline box-border">
                          <div class="mb-2 flex items-center">
                            <div class="flex-grow">
                              <label for="username-input"
                                class="text-neutral-100 font-semibold text-xl">Username</label>
                            </div>
                          </div>
                          <div class="transform-gpu">
                            <div class="relative">
                              <input autofocus id="username-input" name="username" aria-label="Username Input"
                                aria-haspopup="grid" type="text"
                                class="h-14 text-xl [-webkit-filter: invert(100%)] border-neutral-500 text-neutral-100 bg-transparent p-4 flex w-full rounded-lg hover:ring-2 hover:ring-neutral-500 hover:ring-inset hover:ring-offset-neutral-500 focus:ring-2 focus:ring-inset focus:outline-none focus:border-indigo-600 focus:ring-indigo-600 focus:ring-offset-indigo-600"
                                placeholder="" autocapitalize="off" autocorrect="off" autocomplete="username"
                                spellcheck="false" />
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="mt-8 block align-baseline">
                        <div class="holder-of-sorts">
                          <div class="another-holder-of-sorts">
                            <div class="min-h-[2rem] mb-4 flex items-center">
                              <div class="flex-grow">
                                <label for="password-input"
                                  class="text-neutral-100 font-semibold text-xl">Password</label>
                              </div>
                            </div>
                            <div class="relative">
                              <div class="transform-gpu">
                                <div data-a-target="login-password-input" class="relative">
                                  <input id="password-input" name="password" aria-label="Enter your password"
                                    type="password"
                                    class="h-14 text-xl [-webkit-filter: invert(100%)] border-neutral-500 text-neutral-100 bg-transparent p-4 pr-14 flex w-full rounded-lg hover:ring-2 hover:ring-neutral-500 hover:ring-inset hover:ring-offset-neutral-500 focus:ring-2 focus:ring-inset focus:outline-none focus:border-indigo-600 focus:ring-indigo-600 focus:ring-offset-indigo-600"
                                    autocapitalize="off" autocorrect="off" autocomplete="current-password"
                                    data-a-target="tw-input" spellcheck="false" value="" />
                                </div>
                              </div>
                              <div
                                class="absolute flex items-center text-neutral-100 bg-transparent right-0 top-0 bottom-0">
                                <button
                                  class="mr-[0.40rem] text-lg hover:bg-neutral-700 inline-flex items-center justify-center rounded-md h-11 w-11 bg-transparent text-neutral-100 select-none relative align-middle overflow-hidden no-underline whitespace-nowrap bg-none border-none"
                                  type="button" tabindex="-1" aria-label="Toggle password visibility"
                                  onclick="togglePasswordVisibility(this)">
                                  <div class="block pointer-events-none w-8 h-8">
                                    <div class="inline-flex items-center w-full h-full">
                                      <div class="relative w-full overflow-hidden">
                                        <div class="pb-[100%]"></div>
                                        <svg id="password-visibility-icon" width="100%" height="100%" viewBox="0 0 20 20"
                                          aria-hidden="true" focusable="false"
                                          class="absolute left-0 w-full min-h-full top-0 fill-current p-0 m-0 box-border">
                                          <g>
                                            <path d="M11.998 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                                            <path
                                              d="M16.175 7.567L18 10l-1.825 2.433a9.992 9.992 0 01-2.855 2.575l-.232.14a6 6 0 01-6.175 0l-.232-.14a9.992 9.992 0 01-2.855-2.575L2 10l1.825-2.433A9.992 9.992 0 016.68 4.992l.233-.14a6 6 0 016.175 0l.232.14a9.992 9.992 0 012.855 2.575zm-1.6 3.666a7.99 7.99 0 01-2.28 2.058l-.24.144a4 4 0 01-4.11 0 38.552 38.552 0 00-.239-.144 7.994 7.994 0 01-2.28-2.058L4.5 10l.925-1.233a7.992 7.992 0 012.28-2.058 37.9 37.9 0 00.24-.144 4 4 0 014.11 0l.239.144a7.996 7.996 0 012.28 2.058L15.5 10l-.925 1.233z"
                                              fill-rule="evenodd" clip-rule="evenodd"></path>
                                          </g>
                                        </svg>
                                      </div>
                                    </div>
                                  </div>
                                </button>
                              </div>
                            </div>
                            <div class="mt-4">
                              <a class="decoration-transparent text-indigo-400" rel="noopener noreferrer" target="_blank"
                                href="https://www.voluntrips.com/user/account-recovery">
                                <p class="text-lg">Trouble logging in?</p>
                              </a>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="mt-8">
                        <button disabled type="submit" hx-swap="outerHTML" hx-push-url="false" hx-target="#modal-portal"
                          hx-put="/auth/login" id="login-button"
                          class="cursor-not-allowed disabled:bg-neutral-500 bg-indigo-500 hover:bg-indigo-600 inline-flex relative items-center justify-center align-middle overflow-hidden whitespace-nowrap select-none w-full font-semibold rounded-md h-12 text-lg">
                          <div class="flex items-center flex-grow py-0 px-4">
                            <div class="flex-grow flex items-center justify-center">
                              Log in
                            </div>
                          </div>
                        </button>
                      </div>
                      <div id="sign-up-link" class="mt-8 justify-center">
                        <button type="button"
                          class="hover:bg-neutral-700 hover:text-neutral-100 bg-transparent text-indigo-400 inline-flex relative items-center justify-center align-middle overflow-hidden whitespace-nowrap select-none w-full font-semibold rounded-md text-lg h-12">
                          <div class="flex items-center flex-grow-0 px-4 py-0">
                            <div hx-get="/modal/sign-up" hx-target="#modal-portal" hx-swap="innerHTML"
                              hx-trigger="click" hx-push-url="false"
                              class="flex-grow-0 flex items-center justify-start">
                              Don't have an account? Sign up
                            </div>
                          </div>
                        </button>
                      </div>
                    </div>
                  </form>
                </div>
              </div>
            </div>
          </div>
          <div id="modal-close-button" class="left-auto right-4 top-4 absolute ml-2">
            <button
              class="h-12 w-12 rounded-md hover:bg-neutral-500 inline-flex items-center justify-center select-none bg-transparent text-neutral-100"
              aria-label="Close modal" onclick="document.getElementById('modal-overlay').remove();">
              <div class="pointer-events-none w-8 h-8">
                <div class="inline-flex items-center w-full h-full fill-current">
                  <svg fill="rgb(241 245 249)" width="100%" height="100%" viewBox="0 0 20 20" focusable="false"
                    aria-hidden="true">
                    <path
                      d="M8.5 10 4 5.5 5.5 4 10 8.5 14.5 4 16 5.5 11.5 10l4.5 4.5-1.5 1.5-4.5-4.5L5.5 16 4 14.5 8.5 10z">
                    </path>
                  </svg>
                </div>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
<script>
  function togglePasswordVisibility(button) {
    input = document.getElementById("password-input");
    svg = document.getElementById("password-visibility-icon");
    if (input.type === "password") {
      input.type = "text";
      button.setAttribute("aria-label", "Hide password");
      svg.innerHTML = '<g>\n          <path\n            d="M16.5 18l1.5-1.5-2.876-2.876a9.99 9.99 0 001.051-1.191L18 10l-1.825-2.433a9.992 9.992 0 00-2.855-2.575 35.993 35.993 0 01-.232-.14 6 6 0 00-6.175 0 35.993 35.993 0 01-.35.211L3.5 2 2 3.5 16.5 18zm-2.79-5.79a8 8 0 00.865-.977L15.5 10l-.924-1.233a7.996 7.996 0 00-2.281-2.058 37.22 37.22 0 01-.24-.144 4 4 0 00-4.034-.044l1.53 1.53a2 2 0 012.397 2.397l1.762 1.762z"\n            fill-rule="evenodd"\n            clip-rule="evenodd"\n          ></path>\n          <path d="M11.35 15.85l-1.883-1.883a3.996 3.996 0 01-1.522-.532 38.552 38.552 0 00-.239-.144 7.994 7.994 0 01-2.28-2.058L4.5 10l.428-.571L3.5 8 2 10l1.825 2.433a9.992 9.992 0 002.855 2.575c.077.045.155.092.233.14a6 6 0 004.437.702z"></path>\n        </g>';
    } else {
      input.type = "password";
      button.setAttribute("aria-label", "Show password");
      svg.innerHTML = '<g>\n          <path d="M11.998 10a2 2 0 11-4 0 2 2 0 014 0z"></path>\n          <path\n            d="M16.175 7.567L18 10l-1.825 2.433a9.992 9.992 0 01-2.855 2.575l-.232.14a6 6 0 01-6.175 0 35.993 35.993 0 00-.233-.14 9.992 9.992 0 01-2.855-2.575L2 10l1.825-2.433A9.992 9.992 0 016.68 4.992l.233-.14a6 6 0 016.175 0l.232.14a9.992 9.992 0 012.855 2.575zm-1.6 3.666a7.99 7.99 0 01-2.28 2.058l-.24.144a4 4 0 01-4.11 0 38.552 38.552 0 00-
          .239-.144 7.994 7.994 0 01-2.28-2.058L4.5 10l.925-1.233a7.992 7.992 0 012.28-2.058 37.9 37.9 0 00.24-.144 4 4 0 014.11 0l.239.144a7.996 7.996 0 012.28 2.058L15.5 10l-.925 1.233z"\n            fill-rule="evenodd"\n            clip-rule="evenodd"\n          ></path>\n        </g>';
    }
  }

  document
    .getElementById("username-input")
    .addEventListener("input", toggleButtonState);
  document
    .getElementById("password-input")
    .addEventListener("input", toggleButtonState);

  function toggleButtonState() {
    var usernameInput = document.getElementById("username-input");
    var passwordInput = document.getElementById("password-input");
    var button = document.getElementById("login-button");

    if (
      usernameInput.value.trim() !== "" &&
      passwordInput.value.trim() !== ""
    ) {
      button.removeAttribute("disabled");
      button.style.cursor = "pointer";
    } else {
      button.setAttribute("disabled", true);
      button.style.cursor = "not-allowed";
    }
  }
</script>`)
		return nil
	})
}

func ModalSignup() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		// Signup modal content based on voluntrips layout
		_, _ = io.WriteString(w, `<div id="modal-overlay" onclick="if(event.target.id === 'modal-overlay') this?.remove();"
  class="bg-opacity-80 bg-black text-neutral-100 items-start flex fixed bottom-0 left-0 right-0 top-0 z-50 overflow-auto justify-center">
  <div id="modal-content" class="flex flex-grow-0 h-full w-full outline-none justify-center pointer-events-none">
    <div class="p-4 pointer-events-none relative flex w-full h-full items-start justify-center">
      <div class="my-auto outline-none flex flex-grow-0 justify-center relative w-full pointer-events-none">
        <div class="pointer-events-auto max-w-full relative block shadow-md">
          <div id="modal-form" class="flex rounded-md overflow-hidden">
            <div class="w-[50rem] overflow-auto block">
              <div class="flex flex-col bg-neutral-850 p-12">
                <div class="box-border block">
                  <div class="flex flex-col mt-2">
                    <div class="inline-flex items-center justify-center">
                      <div class="inline-flex items-center fill-slate-100 flex-shrink-0">
                        <img class="w-auto h-[50px]" src="/img/logo.png" alt="Voluntrips logo" />
                      </div>
                      <div class="ml-4 text-center">
                        <h4 class="text-2xl font-semibold">
                          Join Voluntrips today
                        </h4>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="mb-4">
                  <form id="modal-form">
                    <div class="w-full flex-col flex">
                      <div class="mt-8">
                        <div class="block align-baseline box-border">
                          <div class="mb-2 flex items-center">
                            <div class="flex-grow">
                              <label for="username" class="text-neutral-100 font-semibold text-xl">Username</label>
                            </div>
                          </div>
                          <div class="transform-gpu">
                            <div class="relative">
                              <input autofocus id="username" name="username" aria-label="Username Input"
                                aria-haspopup="grid" type="text"
                                class="h-14 text-xl [-webkit-filter: invert(100%)] border-neutral-500 text-neutral-100 bg-transparent p-4 flex w-full rounded-lg hover:ring-2 hover:ring-neutral-500 hover:ring-inset hover:ring-offset-neutral-500 focus:ring-2 focus:ring-inset focus:outline-none focus:border-indigo-600 focus:ring-indigo-600 focus:ring-offset-indigo-600"
                                placeholder="" autocapitalize="off" autocorrect="off" autocomplete="username"
                                spellcheck="false" required />
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="mt-8 block align-baseline">
                        <div class="holder-of-sorts">
                          <div class="another-holder-of-sorts">
                            <div class="min-h-[2rem] mb-4 flex items-center">
                              <div class="flex-grow">
                                <label for="password" class="text-neutral-100 font-semibold text-xl">Password</label>
                              </div>
                            </div>
                            <div class="relative">
                              <div class="transform-gpu">
                                <div data-a-target="signup-password" class="relative">
                                  <input id="password" name="password" aria-label="Enter your password" type="password"
                                    class="h-14 text-xl [-webkit-filter: invert(100%)] border-neutral-500 text-neutral-100 bg-transparent p-4 pr-14 flex w-full rounded-lg hover:ring-2 hover:ring-neutral-500 hover:ring-inset hover:ring-offset-neutral-500 focus:ring-2 focus:ring-inset focus:outline-none focus:border-indigo-600 focus:ring-indigo-600 focus:ring-offset-indigo-600"
                                    autocapitalize="off" autocorrect="off" autocomplete="current-password"
                                    data-a-target="tw-input" spellcheck="false" required value="" />
                                </div>
                              </div>
                              <div
                                class="absolute flex items-center text-neutral-100 bg-transparent right-0 top-0 bottom-0">
                                <button
                                  class="mr-[0.40rem] text-lg hover:bg-neutral-700 inline-flex items-center justify-center rounded-md h-11 w-11 bg-transparent text-neutral-100 select-none relative align-middle overflow-hidden no-underline whitespace-nowrap bg-none border-none"
                                  type="button" tabindex="-1" aria-label="Toggle password visibility"
                                  onclick="togglePasswordVisibility(this)">
                                  <div class="block pointer-events-none w-8 h-8">
                                    <div class="inline-flex items-center w-full h-full">
                                      <div class="relative w-full overflow-hidden">
                                        <div class="pb-[100%]"></div>
                                        <svg id="password-visibility-icon" width="100%" height="100%" viewBox="0 0 20 20"
                                          aria-hidden="true" focusable="false"
                                          class="absolute left-0 w-full min-h-full top-0 fill-current p-0 m-0 box-border">
                                          <g>
                                            <path d="M11.998 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                                            <path
                                              d="M16.175 7.567L18 10l-1.825 2.433a9.992 9.992 0 01-2.855 2.575l-.232.14a6 6 0 01-6.175 0l-.239-.144a7.996 7.996 0 01-2.28-2.058L4.5 10l.925-1.233a7.992 7.992 0 012.28-2.058 37.9 37.9 0 00.24-.144 4 4 0 014.11 0l.239.144a7.996 7.996 0 012.28 2.058L15.5 10l-.925 1.233z"
                                              fill-rule="evenodd" clip-rule="evenodd"></path>
                                          </g>
                                        </svg>
                                      </div>
                                    </div>
                                  </div>
                                </button>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="mt-8">
                        <div class="block align-baseline box-border">
                          <div class="mb-2 flex items-center">
                            <div class="flex-grow">
                              <label for="contact" class="text-neutral-100 font-semibold text-xl">Email Address</label>
                            </div>
                          </div>
                          <div class="transform-gpu">
                            <div class="relative">
                              <input id="contact" name="contact" aria-label="Contact Input" aria-haspopup="grid"
                                type="email"
                                class="h-14 text-xl [-webkit-filter: invert(100%)] border-neutral-500 text-neutral-100 bg-transparent p-4 flex w-full rounded-lg hover:ring-2 hover:ring-neutral-500 hover:ring-inset hover:ring-offset-neutral-500 focus:ring-2 focus:ring-inset focus:outline-none focus:border-indigo-600 focus:ring-indigo-600 focus:ring-offset-indigo-600"
                                placeholder="" autocapitalize="off" autocorrect="off" autocomplete="email"
                                spellcheck="false" required />
                            </div>
                          </div>
                        </div>
                      </div>
                      <div class="mt-8">
                        <button disabled type="submit" hx-swap="outerHTML" hx-push-url="false" hx-post="/auth/sign-up"
                          hx-target="#modal-portal" id="signup-button"
                          class="cursor-not-allowed disabled:bg-neutral-500 bg-indigo-500 hover:bg-indigo-600 inline-flex relative items-center justify-center align-middle overflow-hidden whitespace-nowrap select-none w-full font-semibold rounded-md h-12 text-lg">
                          <div class="flex items-center flex-grow py-0 px-4">
                            <div class="flex-grow flex items-center justify-center">
                              Sign up
                            </div>
                          </div>
                        </button>
                      </div>
                      <div id="sign-up-link" class="mt-8 justify-center">
                        <button type="button"
                          class="hover:bg-neutral-700 hover:text-neutral-100 bg-transparent text-indigo-400 inline-flex relative items-center justify-center align-middle overflow-hidden whitespace-nowrap select-none w-full font-semibold rounded-md text-lg h-12">
                          <div class="flex items-center flex-grow-0 px-4 py-0">
                            <div hx-get="/modal/login" hx-target="#modal-portal" hx-swap="innerHTML" hx-trigger="click"
                              hx-push-url="false" class="flex-grow-0 flex items-center justify-start">
                              Already have an account? Log in
                            </div>
                          </div>
                        </button>
                      </div>
                    </div>
                  </form>
                </div>
              </div>
            </div>
          </div>
          <div id="modal-close-button" class="left-auto right-4 top-4 absolute ml-2">
            <button
              class="h-12 w-12 rounded-md hover:bg-neutral-500 inline-flex items-center justify-center select-none bg-transparent text-neutral-100"
              type="button" onclick="document.getElementById('modal-overlay').remove();">
              <div class="pointer-events-none w-8 h-8">
                <div class="inline-flex items-center w-full h-full fill-current">
                  <svg fill="rgb(241 245 249)" width="100%" height="100%" viewBox="0 0 20 20" focusable="false"
                    aria-hidden="true">
                    <path
                      d="M8.5 10 4 5.5 5.5 4 10 8.5 14.5 4 16 5.5 11.5 10l4.5 4.5-1.5 1.5-4.5-4.5L5.5 16 4 14.5 8.5 10z">
                    </path>
                  </svg>
                </div>
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
<script>
  function togglePasswordVisibility(button) {
    input = document.getElementById("password");
    svg = document.getElementById("password-visibility-icon");
    if (input.type === "password") {
      input.type = "text";
      button.setAttribute("aria-label", "Hide password");
      svg.innerHTML = '<g>\n          <path\n            d="M16.5 18l1.5-1.5-2.876-2.876a9.99 9.99 0 001.051-1.191L18 10l-1.825-2.433a9.992 9.992 0 00-2.855-2.575 35.993 35.993 0 01-.232-.14 6 6 0 00-6.175 0 35.993 35.993 0 01-.35.211L3.5 2 2 3.5 16.5 18zm-2.79-5.79a8 8 0 00.865-.977L15.5 10l-.924-1.233a7.996 7.996 0 00-2.281-2.058 37.22 37.22 0 01-.24-.144 4 4 0 00-4.034-.044l1.53 1.53a2 2 0 012.397 2.397l1.762 1.762z"\n            fill-rule="evenodd"\n            clip-rule="evenodd"\n          ></path>\n          <path d="M11.35 15.85l-1.883-1.883a3.996 3.996 0 01-1.522-.532 38.552 38.552 0 00-.239-.144 7.994 7.994 0 01-2.28-2.058L4.5 10l.428-.571L3.5 8 2 10l1.825 2.433a9.992 9.992 0 002.855 2.575c.077.045.155.092.233.14a6 6 0 004.437.702z"></path>\n        </g>';
    } else {
      input.type = "password";
      button.setAttribute("aria-label", "Show password");
      svg.innerHTML = '<g>\n          <path d="M11.998 10a2 2 0 11-4 0 2 2 0 014 0z"></path>\n          <path\n            d="M16.175 7.567L18 10l-1.825 2.433a9.992 9.992 0 01-2.855 2.575l-.232.14a6 6 0 01-6.175 0 35.993 35.993 0 00-.233-.14 9.992 9.992 0 01-2.855-2.575L2 10l1.825-2.433A9.992 9.992 0 016.68 4.992l.233-.14a6 6 0 016.175 0l.232.14a9.992 9.992 0 012.855 2.575zm-1.6 3.666a7.99 7.99 0 01-2.28 2.058l-.24.144a4 4 0 01-4.11 0 38.552 38.552 0 00-.239-.144 7.994 7.994 0 01-2.28-2.058L4.5 10l.925-1.233a7.992 7.992 0 012.28-2.058 37.9 37.9 0 00.24-.144 4 4 0 014.11 0l.239.144a7.996 7.996 0 012.28 2.058L15.5 10l-.925 1.233z"\n            fill-rule="evenodd"\n            clip-rule="evenodd"\n          ></path>\n        </g>';
    }
  }

  document
    .getElementById("username")
    .addEventListener("input", toggleButtonState);
  document
    .getElementById("password")
    .addEventListener("input", toggleButtonState);
  document
    .getElementById("contact")
    .addEventListener("input", toggleButtonState);

  function toggleButtonState() {
    var usernameInput = document.getElementById("username");
    var passwordInput = document.getElementById("password");
    var contactInput = document.getElementById("contact");
    var button = document.getElementById("signup-button");

    if (
      usernameInput.value.trim() !== "" &&
      passwordInput.value.trim() !== "" &&
      contactInput.value.trim() !== ""
    ) {
      button.removeAttribute("disabled");
      button.style.cursor = "pointer";
    } else {
      button.setAttribute("disabled", true);
      button.style.cursor = "not-allowed";
    }
  }
</script>`)
		return nil
	})
}

func DropdownSearch() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div
  id="search-drop-down-overlay"
  class="bg-neutral-800 text-neutral-100 absolute flex left-auto bottom-auto top-14 -right-40 top rounded-md mt-0 drop-shadow-lg"
>
  <div id="drop-down-content" class="pointer-events-none p-1">
    <div class="pointer-events-none flex">
      <div>
        <div class="pointer-events-auto max-w-full flex flex-row">
          <div class="w-[25rem]">
            <div class="relative">
              <input
                autofocus
                id="search-input"
                aria-label="Search Input"
                aria-haspopup="grid"
                type="search"
                class="h-14 [-webkit-filter: invert(100%)] border-neutral-500 text-neutral-100 placeholder-neutral-200 bg-transparent text-xl p-4 pr-2 flex w-full rounded-l-lg hover:ring-2 hover:ring-neutral-500 hover:ring-inset hover:ring-offset-neutral-500 focus:ring-2 focus:ring-inset focus:outline-none focus:border-indigo-600 focus:ring-indigo-600 focus:ring-offset-indigo-600"
                placeholder="Search"
                autocapitalize="off"
                autocorrect="off"
                maxlength="70"
                spellcheck="false"
              />
            </div>
          </div>
          <button
            class="bg-neutral-700 flex-shrink-0 p-2 rounded-r-md inline-flex relative items-center align-middle justify-center overflow-hidden whitespace-nowrap select-none"
          >
            <div class="flex items-center">
              <div class="inline-flex items-center w-10">
                <div class="inline-flex items-center h-full w-full">
                  <div class="relative w-full overflow-hidden">
                    <svg
                      fill="#fff"
                      width="100%"
                      height="100%"
                      viewBox="0 0 20 20"
                      class="ScIconSVG-sc-1q25cff-1 jpczqG"
                    >
                      <g>
                        <path
                          fill-rule="evenodd"
                          d="M13.192 14.606a7 7 0 111.414-1.414l3.101 3.1-1.414 1.415-3.1-3.1zM14 9A5 5 0 114 9a5 5 0 0110 0z"
                          clip-rule="evenodd"
                        ></path>
                      </g>
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </button>
        </div>
      </div>
    </div>
  </div>

  <script>
    // if the document is clicked, and it isn't the search-drop-down-overlay, remove it
    document.addEventListener("click", function (event) {
      if (
        !event.target.closest("#drop-down-content") &&
        !event.target.closest("#drop-down-portal") &&
        !event.target.closest("#search-drop-down-overlay") &&
        document.getElementById("search-drop-down-overlay")
      ) {
        document.getElementById("search-drop-down-overlay").remove();
      }
    });
  </script>
</div>`)
		return nil
	})
}

func DropdownMore() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="more-drop-down-overlay"
  class="bg-neutral-800 text-neutral-100 absolute flex left-auto bottom-auto items-center top rounded-md mt-0 drop-shadow-lg">
  <div id="drop-down-content" class="pointer-events-none p-4">
    <div class="pointer-events-none flex">
      <div>
        <div class="pointer-events-auto max-w-full block">
          <div class="relative">
            <div class="flex flex-col">
              <div class="flex-nowrap flex">
                <h5
                  class="flex items-center flex-grow-0 px-2 my-2 text-xl text-gray-400 text-semibold whitespace-nowrap">
                  GENERAL
                </h5>
              </div>
              <div onclick="document.getElementById('more-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/about" hx-get="/about" hx-target="#partial" hx-swap="outerHTML " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-lg font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0">
                      <div class="flex-grow-0 flex items-center justify-start">
                        About
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <div onclick="document.getElementById('more-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/developers" hx-get="/developers" hx-target="#partial" hx-swap="outerHTML "
                    hx-trigger="click" hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-lg font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Developers
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <hr class="opacity-20 my-4" />
              <div class="flex-nowrap flex">
                <h5
                  class="flex items-center flex-grow-0 px-2 my-2 text-xl text-gray-400 text-semibold whitespace-nowrap">
                  HELP & LEGAL
                </h5>
              </div>
              <div onclick="document.getElementById('more-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/tos" hx-get="/tos" hx-target="#partial" hx-swap="outerHTML " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-lg font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Terms of Service
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <div onclick="document.getElementById('more-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/policy" hx-get="/policy" hx-target="#partial" hx-swap="outerHTML " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-lg font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Privacy Policy
                      </div>
                    </div>
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <script>
    // if the document is clicked, and it isn't the more-drop-down-overlay, remove it
    document.addEventListener("click", function (event) {
      if (
        !event.target.closest("#drop-down-content") &&
        !event.target.closest("#drop-down-portal") &&
        !event.target.closest("#more-drop-down-overlay") &&
        document.getElementById("more-drop-down-overlay")
      ) {
        document.getElementById("more-drop-down-overlay").remove();
      }
    });
  </script>
</div>`)
		return nil
	})
}

func DropdownUser(userLoggedIn bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, _ = io.WriteString(w, `<div id="user-drop-down-overlay"
  class="bg-neutral-800 text-neutral-100 absolute flex left-auto bottom-auto -right-16 -top-1 rounded-md mt-0 drop-shadow-lg shadow-lg">
  <div id="drop-down-content" class="pointer-events-none p-4">
    <div class="pointer-events-none flex">
      <div>
        <div class="pointer-events-auto max-w-full block">`)
		if userLoggedIn {
			_, _ = io.WriteString(w, `<div class="relative">
            <div class="flex flex-col">
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/profile" hx-get="/profile" hx-target="#partial" hx-swap="outerHTML  " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Profile
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/account" hx-get="/account" hx-target="#partial" hx-swap="outerHTML  " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Account
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <hr class="opacity-20 my-4" />
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/trips"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Manage Your Trips
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <hr class="opacity-20 my-4" />
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/tos" hx-get="/tos" hx-target="#partial" hx-swap="outerHTML  " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Terms of Service
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/policy" hx-get="/policy" hx-target="#partial" hx-swap="outerHTML  " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Privacy Policy
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <hr class="opacity-20 my-4" />
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <button hx-trigger="click" hx-push-url="false" hx-post="/auth/logout"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Logout
                      </div>
                    </div>
                  </button>
                </div>
              </div>
            </div>
          </div>`)
		} else {
			_, _ = io.WriteString(w, `<div class="relative">
            <div class="flex flex-col">
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/tos" hx-get="/tos" hx-target="#partial" hx-swap="outerHTML  " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Terms of Service
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <a href="/policy" hx-get="/policy" hx-target="#partial" hx-swap="outerHTML  " hx-trigger="click"
                    hx-push-url="true"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Privacy Policy
                      </div>
                    </div>
                  </a>
                </div>
              </div>
              <hr class="opacity-20 my-4" />
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <button hx-target="#modal-portal" hx-swap="innerHTML" hx-trigger="click" hx-push-url="false"
                    hx-get="/modal/login"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Sign Up
                      </div>
                    </div>
                  </button>
                </div>
              </div>
              <div onclick="document.getElementById('user-drop-down-overlay').remove();" class="flex-nowrap flex">
                <div class="flex w-full">
                  <button hx-target="#modal-portal" hx-swap="innerHTML" hx-trigger="click" hx-push-url="false"
                    hx-get="/modal/login"
                    class="p-2 w-full hover:bg-neutral-600 text-neutral-100 text-xl font-medium rounded-md relative items-center justify-start align-middle overflow-hidden whitespace-nowrap select-none">
                    <div class="flex items-center flex-grow-0 px-4">
                      <div class="flex-grow-0 flex items-center justify-start">
                        Login
                      </div>
                    </div>
                  </button>
                </div>
              </div>
            </div>
          </div>`)
		}
		_, _ = io.WriteString(w, `</div>
      </div>
    </div>
  </div>

  <script>
    // if the document is clicked, and it isn't the user-drop-down-overlay, remove it
    document.addEventListener("click", function (event) {
      if (
        !event.target.closest("#drop-down-content") &&
        !event.target.closest("#drop-down-portal") &&
        !event.target.closest("#user-drop-down-overlay") &&
        document.getElementById("user-drop-down-overlay")
      ) {
        document.getElementById("user-drop-down-overlay").remove();
      }
    });
  </script>
</div>`)
		return nil
	})
}
