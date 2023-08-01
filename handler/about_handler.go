package handler

import (
  "encoding/json"

  "io/ioutil"
  "net/http"

  "github.com/labstack/echo"

  "SampleApplicationGo/model"
)

func AboutHandler(c echo.Context) error {

  tr := &http.Transport{}
  client := &http.Client{Transport: tr}

  // Call the api
  resp, err := client.Get(
    "http://localhost:1323/api/get-full-name?first_name=Abhimanyu&last_name=kumar",
  )

  if err != nil {
    return err
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  // Unmarshal the response into a ExampleResponse struct
  var exampleResponse model.ExampleResponse
  if err = json.Unmarshal(body, &exampleResponse); err != nil {
    return err
  }

  // Please note the the second parameter "about.html" is the template name and should
  // be equal to one of the keys in the TemplateRegistry array defined in main.go
  return c.Render(http.StatusOK, "about.html", map[string]interface{}{
    "name": "About",
    "msg":  exampleResponse.Msg,
  })
}
