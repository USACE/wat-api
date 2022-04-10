package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/usace/wat-api/wat"
)

type WatHandler struct {
}

const version = "2.0.1 Development"

func CreateWatHandler() *WatHandler {
	wh := WatHandler{}
	return &wh
}
func (wh *WatHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("WAT API Version %s", version))
}
func (wh *WatHandler) Plugins(c echo.Context) error {
	//ping the network to figure out what plugins are active?
	plugins := make([]wat.Plugin, 2)
	plugins[0] = wat.Plugin{Name: "plugin a"}
	plugins[1] = wat.Plugin{Name: "plugin b"}
	return c.JSON(http.StatusOK, plugins)
}
