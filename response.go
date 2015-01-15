package Golf

import (
	"encoding/json"
	"net/http"
	"time"
)

// A wrapper of http.ResponseWriter
type Response struct {
	http.ResponseWriter
	app        *Application
	StatusCode int
	Body       []byte
}

func NewResponse(res http.ResponseWriter, app *Application) *Response {
	response := new(Response)
	response.ResponseWriter = res
	response.app = app
	response.Header().Set("Content-Type", "text/html;charset=UTF-8")
	return response
}

func (res *Response) Send(str string) {
	res.Body = []byte(str)
}

func (res *Response) JSON(obj interface{}) {
	js, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	res.Body = js
	res.Header().Set("Content-Type", "application/json")
}

func (res *Response) Status(code int) {
	res.WriteHeader(code)
}

func (res *Response) Redirect(url string, code int) {
	res.Header().Set("Location", url)
	res.WriteHeader(code)
	res.StatusCode = code
}

func (res *Response) SetCookie(key string, value string, expire int) {
	now := time.Now()
	expireTime := now.Add(time.Duration(expire) * time.Second)
	cookie := &http.Cookie{
		Name:    key,
		Value:   value,
		Path:    "/",
		MaxAge:  expire,
		Expires: expireTime,
	}
	http.SetCookie(res, cookie)
}

func (res *Response) Render(file_path string, arguments map[string]interface{}) {
	result, e := res.app.view.Render(file_path, arguments)
	if e != nil {
		panic(e)
	}
	res.Send(result)
}

// Sets the response’s HTTP header field to val.
func (res *Response) Set(key, val string) {
	res.Header().Set(key, val)
}

func (res *Response) Get(key string) string {
	return res.Header().Get(key)
}
