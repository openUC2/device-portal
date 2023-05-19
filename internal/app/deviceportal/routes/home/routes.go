// Package home contains the route handlers related to the app's home screen.
package home

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/sargassum-world/godest"
)

type Handlers struct {
	r godest.TemplateRenderer
}

func New(r godest.TemplateRenderer) *Handlers {
	return &Handlers{
		r: r,
	}
}

func (h *Handlers) Register(er godest.EchoRouter) {
	er.GET("/", h.HandleHomeGet())
}

type HomeViewData struct{}

func getHomeViewData(ctx context.Context) (vd HomeViewData, err error) {
	return vd, nil
}

func (h *Handlers) HandleHomeGet() echo.HandlerFunc {
	t := "home/home.page.tmpl"
	h.r.MustHave(t)
	return func(c echo.Context) error {
		// Run queries
		ctx := c.Request().Context()
		homeViewData, err := getHomeViewData(ctx)
		if err != nil {
			return err
		}
		// Produce output
		return h.r.CacheablePage(c.Response(), c.Request(), t, homeViewData, struct{}{})
	}
}
