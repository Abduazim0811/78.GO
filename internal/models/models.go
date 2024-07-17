package models

type User struct {
    Name string `json:"name" xml:"name"`
    Age  int    `json:"age"  xml:"age"`
}
