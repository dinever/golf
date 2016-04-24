package golf

import (
	"testing"
)

func assertStringEqual(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func assertSliceEqual(t *testing.T, expected, got []string) {
	if len(expected) != len(got) {
		t.Errorf("Slice length not equal, expected: %v, got %v", expected, got)
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != got[i] {
			t.Errorf("Slice not equal, expected: %v, got %v", expected, got)
		}
	}
}

type route struct {
	method   string
	path     string
	testPath string
	params   map[string]string
}

var githubAPI = []route{
	// OAuth Authorizations
	{"GET", "/authorizations", "/authorizations", map[string]string{}},
	{"GET", "/auth", "/auth", map[string]string{}},
	{"GET", "/authorizations/:id", "/authorizations/12345", map[string]string{"id": "12345"}},
	{"POST", "/authorizations", "/authorizations", map[string]string{}},
	{"DELETE", "/authorizations/:id", "/authorizations/12345", map[string]string{"id": "12345"}},
	{"GET", "/applications/:client_id/tokens/:access_token", "/applications/12345/tokens/67890", map[string]string{"client_id": "12345", "access_token": "67890"}},
	{"DELETE", "/applications/:client_id/tokens", "/applications/12345/tokens", map[string]string{"client_id": "12345"}},
	{"DELETE", "/applications/:client_id/tokens/:access_token", "/applications/12345/tokens/67890", map[string]string{"client_id": "12345", "access_token": "67890"}},

	// Activity
	{"GET", "/events", "/events", nil},
	{"GET", "/repos/:owner/:repo/events", "/repos/dinever/golf/events", map[string]string{"owner": "dinever", "repo": "golf"}},
	{"GET", "/networks/:owner/:repo/events", "/networks/dinever/golf/events", map[string]string{"owner": "dinever", "repo": "golf"}},
	{"GET", "/orgs/:org/events", "/orgs/golf/events", map[string]string{"org": "golf"}},
	{"GET", "/users/:user/received_events", "/users/dinever/received_events", nil},
	{"GET", "/users/:user/received_events/public", "/users/dinever/received_events/public", nil},
}

func handler(ctx *Context) {
}

func TestRouter(t *testing.T) {
	router := newRouter()
	for _, route := range githubAPI {
		router.AddRoute(route.method, route.path, handler)
	}

	for _, route := range githubAPI {
		_, param, err := router.FindRoute(route.method, route.testPath)
		if err != nil {
			t.Errorf("Can not find route: %v", route.testPath)
		}

		for key, expected := range route.params {
			val, err := param.ByName(key)
			if err != nil {
				t.Errorf("Can not retrieve parameter from route %v: %v", route.testPath, key)
			} else {
				assertStringEqual(t, expected, val)
			}
			val, err = param.ByName(key)
			if err != nil {
				t.Errorf("Can not retrieve parameter from route %v: %v", route.testPath, key)
			} else {
				assertStringEqual(t, expected, val)
			}
		}
	}
}

func TestSplitURLPath(t *testing.T) {

	var table = map[string][2][]string{
		"/users/:name":                         {{"/users/", ":"}, {"name"}},
		"/users/:name/put":                     {{"/users/", ":", "/put"}, {"name"}},
		"/users/:name/put/:section":            {{"/users/", ":", "/put/", ":"}, {"name", "section"}},
		"/customers/:name/put/:section":        {{"/customers/", ":", "/put/", ":"}, {"name", "section"}},
		"/customers/groups/:name/put/:section": {{"/customers/groups/", ":", "/put/", ":"}, {"name", "section"}},
	}

	for path, result := range table {
		parts, _ := splitURLPath(path)
		assertSliceEqual(t, parts, result[0])
	}
}

func TestIncorrectPath(t *testing.T) {
	path := "/users/foo:name/"
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	router := newRouter()
	router.AddRoute("GET", path, handler)
	t.Errorf("Incorrect path should raise an error.")
}

func TestPathNotFound(t *testing.T) {
	path := []struct {
		method, path, incomingMethod, incomingPath string
	}{
		{"GET", "/users/name/", "GET", "/users/name/dinever/"},
		{"GET", "/dinever/repo/", "POST", "/dinever/repo/"},
	}
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	router := newRouter()
	for _, path := range path {
		router.AddRoute(path.method, path.path, handler)
		h, p, err := router.FindRoute(path.incomingMethod, path.incomingPath)
		if h != nil {
			t.Errorf("Should return nil handler when path not found.")
		}
		if p.Len() != 0 {
			t.Errorf("Should return nil parameter when path not found.")
		}
		if err == nil {
			t.Errorf("Should rasie an error when path not found.")
		}
	}
}
