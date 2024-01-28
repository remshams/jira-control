package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/remshams/jira-control/web/utils"
	"github.com/remshams/jira-control/web/view"
)

func HandleLogWork(c echo.Context) error {
	return utils.Render(c, view.LogWork())
}
