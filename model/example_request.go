package model

type ExampleRequest struct {
  FirstName string `json:"first_name" form:"first_name" query:"first_name"`
  LastName  string `json:"last_name" form:"last_name" query:"last_name"`
}
