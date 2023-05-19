package deviceportal

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sargassum-world/godest"
	"github.com/sargassum-world/godest/httperr"
)

type ErrorData struct {
	Code     int
	Error    httperr.DescriptiveError
	Messages []string
}

func NewHTTPErrorHandler(tr godest.TemplateRenderer) echo.HTTPErrorHandler {
	tr.MustHave("app/httperr.page.tmpl")
	return func(err error, c echo.Context) {
		c.Logger().Error(err)

		// Process error code
		code := http.StatusInternalServerError
		if herr, ok := err.(*echo.HTTPError); ok {
			code = herr.Code
		}
		errorData := ErrorData{
			Code:  code,
			Error: httperr.Describe(code),
		}

		// Produce output
		if perr := tr.Page(
			c.Response(), c.Request(), code, "app/httperr.page.tmpl", errorData, struct{}{},
			godest.WithUncacheable(),
		); perr != nil {
			c.Logger().Error(errors.Wrap(perr, "couldn't render error page in error handler"))
		}
	}
}
