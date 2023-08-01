package handler

import (
  "net/http"

  "github.com/labstack/echo"
)

func HomeHandler(c echo.Context) error {
  // Please note the the second parameter "home.html" is the template name and should
  // be equal to one of the keys in the TemplateRegistry array defined in main.go
  return c.Render(http.StatusOK, "home.html", map[string]interface{}{
    "name": "HOME",
    "msg":  "Hello, GO!",
  })
}
