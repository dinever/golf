package Yafg

import (
  "net/http"
  "time"
)

// A wrapper of http.ResponseWriter
type Response struct {
  http.ResponseWriter
}

func NewResponse(res http.ResponseWriter) *Response {
  response := new(Response)
  response.ResponseWriter = res
  return response
}

func (res *Response) Send(str string) {
  res.Write([]byte(str))
}

func (res *Response) Redirect(url string, code int) {
  res.Header().Set("Location", url)
  res.WriteHeader(code)
}

func (res *Response) SetCookie(key string, value string, expire int) {
  now := time.Now()
  expireTime := now.Add(time.Duration(expire) * time.Second)
  cookie := &http.Cookie{
    Name: key,
    Value: value,
    Path: "/",
    MaxAge: expire,
    Expires: expireTime,
  }
  http.SetCookie(res, cookie)
}
