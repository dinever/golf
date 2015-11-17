package Yafg

import (
  "net/http"
)

// A wrapper of http.ResponseWriter
type Response struct {
  Response http.ResponseWriter
}

func NewResponse(res http.ResponseWriter) *Response {
  response := new(Response)
  response.Response = res
  return response
}
