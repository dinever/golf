package Golf

import (
	"fmt"
	"bytes"
	"testing"
)

func TestConfig(t *testing.T) {
	cases := []struct {
		key   string
		value interface{}
	}{
		{"foo", "bar"},
		{"foo", 123},
		{"foo", true},
		{"foo", 56.23},
		{"foo/bar", "bar"},
		{"123/foo/bar", "bar"},
		{"foo/bar/bar/bar", "bar"},
	}

	defaultValue := "None"

	for _, c := range cases {
		app := New()
		config := NewConfig(app)
		config.Set(c.key, c.value)

		value, err := config.Get(c.key, defaultValue)
		if err != nil {
			t.Error(err)
		}
		if value != c.value {
			t.Errorf("Value not match: %q != %q, key: %q", value, c.value, c.key)
		}
	}
}

func TestConfigWithMultipleEntires(t *testing.T) {
	settings := []struct {
		key, value string
	}{
		{"foo/bar/bar", "bar"},
		{"foo/bar/bar2", "bar2"},
		{"foo/bar3", "bar3"},
		{"foo2", "bar4"},
	}

	app := New()
	config := NewConfig(app)

	for _, c := range settings {
		config.Set(c.key, c.value)
	}

	for _, c := range settings {
		value, err := config.Get(c.key, "None")
		if err != nil {
			t.Error(err)
		}
		if value != c.value {
			t.Errorf("Value not match: %q != %q, key: %q", value, c.value, c.key)
		}
	}
}

func TestFromJSON(t *testing.T) {
	reader := bytes.NewReader([]byte(`{"cool" : {"foo" : "bar"}}`))
	config, err := ConfigFromJSON(reader)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	value, err := config.GetString("cool/foo", "")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if value != "bar" {
		t.Error(fmt.Sprintf("expected value to be abc but it was %v", value))
		t.Fail()
	}
}
