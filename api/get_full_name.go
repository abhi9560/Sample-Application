package api

import (
  "fmt"
  "net/http"

  "github.com/labstack/echo"

  "SampleApplicationGo/model"
)

func GetFullName(c echo.Context) error {
  // Bind the input data to ExampleRequest
  exampleRequest := new(model.ExampleRequest)
  if err := c.Bind(exampleRequest); err != nil {
    return err
  }

  // Manipulate the input data
  greeting := exampleRequest.FirstName + " " + exampleRequest.LastName

  return c.JSONBlob(
    http.StatusOK,
    []byte(
      fmt.Sprintf(`{
        "first_name": %q,
        "last_name": %q,
        "msg": "Hello %s"
      }`, exampleRequest.FirstName, exampleRequest.LastName, greeting),
    ),
  )
}
