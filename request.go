package Yafg

import (
  "net/http"
)

// A wrapper of http.Request
type Request struct {
  Request *http.Request
	Params map[string]string
	URL string
}

func NewRequest(req *http.Request) *Request {
  request := new(Request)
  request.Request = req
  return request
}

