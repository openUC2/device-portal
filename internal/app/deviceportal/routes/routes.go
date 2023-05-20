// Package routes contains the route handlers for the Planktoscope Live server.
package routes

import (
	"github.com/sargassum-world/godest"

	"github.com/PlanktoScope/device-portal/internal/app/deviceportal/client"
	"github.com/PlanktoScope/device-portal/internal/app/deviceportal/routes/assets"
	"github.com/PlanktoScope/device-portal/internal/app/deviceportal/routes/home"
)

type Handlers struct {
	r       godest.TemplateRenderer
	globals *client.Globals
}

func New(r godest.TemplateRenderer, globals *client.Globals) *Handlers {
	return &Handlers{
		r:       r,
		globals: globals,
	}
}

func (h *Handlers) Register(er godest.EchoRouter, em godest.Embeds) {
	assets.RegisterStatic(er, em)
	assets.NewTemplated(h.r).Register(er)
	home.New(h.r, h.globals.Config.MachineName).Register(er)
}
