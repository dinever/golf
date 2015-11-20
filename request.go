package Yafg

import (
  "strings"
  "net/http"
)

// A wrapper of http.Request
type Request struct {
  *http.Request
	Params map[string]string
	IP string
}

func NewRequest(req *http.Request) *Request {
  request := new(Request)
  request.Request = req
	request.IP = strings.Split(req.RemoteAddr, ":")[0]
  return request
}

func (req *Request) GetCookie(key string) string {
  cookie, err := req.Request.Cookie(key)
  if err != nil {
    return ""
  } else {
    return cookie.Value
  }
}
