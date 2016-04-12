package Golf

import (
	"testing"
)

func TestStringRendering(t *testing.T) {
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
		view := NewView()
		result, err := view.RenderFromString("", c.content, c.args)
		if err != nil {
			t.Errorf("Can not render from string")
		}
		if result != c.output {
			t.Errorf("Rendered content %q != %q", result, c.output)
		}
	}
}

func TestSetTemplateLoader(t *testing.T) {
	view := NewView()
	view.SetTemplateLoader("test", "/test/path/")
	if view.templateLoader["test"] == nil {
		t.Errorf("Could not set template loader for view")
	}
}
