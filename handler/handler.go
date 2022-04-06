package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
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
