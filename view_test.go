package Golf

import (
	"testing"
)

func TestRenderFromString(t *testing.T) {
	cases := []struct {
		content string
		args    map[string]interface{}
		output  string
	}{
		{
			"This is a sample template without args.",
			map[string]interface{}{},
			"This is a sample template without args.",
		},
		{
			"Article: {{ .title }}",
			map[string]interface{}{
				"title": "Hello World",
			},
			"Article: Hello World",
		},
		{
			"Article: {{ .title }}, Count: {{ .count }}",
			map[string]interface{}{
				"title": "Hello World",
				"count": 5,
			},
			"Article: Hello World, Count: 5",
		},
	}

	for _, c := range cases {
		view := NewView(".")
		result, error := view.RenderFromString(c.content, c.args)
		if error != nil {
			panic(error)
		}
		if result != c.output {
			t.Errorf("Rendered content %q != %q", result, c.output)
		}
	}
}
