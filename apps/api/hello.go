package api

import (
	"github.com/labstack/echo"
	"net/http"
)

type HelloAPI struct {

}

func (api HelloAPI) GetHello(c echo.Context) error  {

	return c.HTML(http.StatusOK, "Hello World")
}