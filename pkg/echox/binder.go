package echox

import (
	"github.com/labstack/echo/v4"
)

type Binder struct{}

func (Binder) Bind(i interface{}, c echo.Context) error {
	var binder echo.DefaultBinder

	err := binder.Bind(i, c)
	if err != nil {
		return err
	}

	return c.Validate(i)
}
