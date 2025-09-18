package server

import (
	"io/fs"

	assets "github.com/Taiterbase/vtrips/apps/frontend/web/static"
	"github.com/labstack/echo/v4"
)

func registerStatic(e *echo.Echo) {
	// Expose embedded static files at /static
	var root fs.FS = assets.Static
	e.StaticFS("/static", root)
	// Common aliases for convenience when porting
	if sub, err := fs.Sub(root, "css"); err == nil {
		e.StaticFS("/css", sub)
	}
	if sub, err := fs.Sub(root, "js"); err == nil {
		e.StaticFS("/js", sub)
	}
	if sub, err := fs.Sub(root, "img"); err == nil {
		e.StaticFS("/img", sub)
	}
	if sub, err := fs.Sub(root, "svg"); err == nil {
		e.StaticFS("/svg", sub)
	}
}
