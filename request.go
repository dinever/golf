package Yafg

import (
  "net/http"
)

// A wrapper of http.Request
type Request struct {
  *http.Request
	Params map[string]string
}

func NewRequest(req *http.Request) *Request {
  request := new(Request)
  request.Request = req
  return request
}

