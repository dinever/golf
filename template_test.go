package Golf

import (
	"bytes"
	"testing"
)

type TemplateData struct {
	Title string
	Path  string
	User  interface{}
	Nav   map[string]string
	Data  map[string]interface{}
}

func TestTemplateRendering(t *testing.T) {
	st := &TemplateManager{
		Loader: &MapLoader{
			"test.html":    `<title>{{.Title}}</title> Key={{.Data}}`,
			"test2.html":   "Lorem ipsum dolor sit amet.",
			"test3.html":   `{{.Foo}}`,
			"extends.html": `header {{ template "content" }} footer`,
		},
	}

	cases := []struct {
		content string
		args    map[string]interface{}
		output  string
	}{
		{
			"Lorem ipsum dolor sit amet.",
			map[string]interface{}{},
			"Lorem ipsum dolor sit amet.",
		},
		{
			`Lorem ipsum {{ include "test2.html" }} sit amet.`,
			map[string]interface{}{},
			"Lorem ipsum Lorem ipsum dolor sit amet. sit amet.",
		},
		{
			`{{ extends "extends.html" }} {{ define "content" }}Lorem ipsum Lorem ipsum dolor sit amet. sit amet.{{ end }}`,
			map[string]interface{}{},
			`header Lorem ipsum Lorem ipsum dolor sit amet. sit amet. footer`,
		},
	}

	for _, c := range cases {
		var buf bytes.Buffer
		st.RenderFromString(&buf, c.content, c.args)
		if buf.String() != c.output {
			t.Errorf("String template rendering with loader failed: %q != %q", buf.String(), c.output)
		}
	}
}
