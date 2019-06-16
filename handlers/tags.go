package handlers

import (
	"dou-parser/parser"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetTags(c echo.Context) error {
	err, tags := parser.ScrapEventTags()

	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, tags)
}
