package Golf

import (
	"errors"
	"net/http"
	"strings"
)

// A wrapper of http.Request
type Request struct {
	*http.Request
	Params map[string]string
	IP     string
}

func NewRequest(req *http.Request) *Request {
	request := new(Request)
	request.Request = req
	request.IP = strings.Split(req.RemoteAddr, ":")[0]
	return request
}

// Query returns query data by the query key.
func (req *Request) Query(key string, index ...int) (string, error) {
	req.ParseForm()
	if val, ok := req.Form[key]; ok {
		if len(index) == 1 {
			return val[index[0]], nil
		} else {
			return val[0], nil
		}
	} else {
		return "", errors.New("Query key not found.")
	}
}

// Cookie returns request cookie item string by a given key.
func (req *Request) Cookie(key string) string {
	cookie, err := req.Request.Cookie(key)
	if err != nil {
		return ""
	} else {
		return cookie.Value
	}
}

// Protocol returns the request protocol string
func (req *Request) Protocol() string {
	return req.Proto
}
