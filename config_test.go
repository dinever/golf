package golf

import (
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
		config := NewConfig()
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

	config := NewConfig()

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
		t.Errorf("expected value to be abc but it was %v", value)
	}
}

func TestGetStringException(t *testing.T) {
	defaultValue := "None"

	config := NewConfig()
	config.Set("foo", 123)
	val, err := config.GetString("foo", defaultValue)
	if err == nil {
		t.Errorf("Should have raised an type error when getting a non-string value by GetString.")
	}
	if val != defaultValue {
		t.Errorf("Should have used the default value when raising an error.")
	}

	val, err = config.GetString("bar", defaultValue)
	if err == nil {
		t.Errorf("Should have raised an type error when getting a non-existed value by GetString.")
	}
	if val != defaultValue {
		t.Errorf("Should have used the default value when raising an error.")
	}
}

func TestGetIntegerException(t *testing.T) {
	defaultValue := 123

	config := NewConfig()
	config.Set("foo", "bar")
	val, err := config.GetInt("foo", defaultValue)
	if err == nil {
		t.Errorf("Should have raised an type error when getting a non-string value by GetString.")
	}
	if val != defaultValue {
		t.Errorf("Should have used the default value when raising an error.")
	}

	val, err = config.GetInt("bar", defaultValue)
	if err == nil {
		t.Errorf("Should have raised an type error when getting a non-existed value by GetString.")
	}
	if val != defaultValue {
		t.Errorf("Should have used the default value when raising an error.")
	}
}

func TestGetBoolException(t *testing.T) {
	defaultValue := false

	config := NewConfig()
	config.Set("foo", "bar")
	val, err := config.GetBool("foo", defaultValue)
	if err == nil {
		t.Errorf("Should have raised an type error when getting a non-string value by GetString.")
	}
	if val != defaultValue {
		t.Errorf("Should have used the default value when raising an error.")
	}

	val, err = config.GetBool("bar", defaultValue)
	if err == nil {
		t.Errorf("Should have raised an type error when getting a non-existed value by GetString.")
	}
	if val != defaultValue {
		t.Errorf("Should have used the default value when raising an error.")
	}
}

func TestGetFloat64Exception(t *testing.T) {
	defaultValue := 0.5

	config := NewConfig()
	config.Set("foo", "bar")
	val, err := config.GetFloat("foo", defaultValue)
	if err == nil {
		t.Errorf("Should have raised an type error when getting a non-string value by GetString.")
	}
	if val != defaultValue {
		t.Errorf("Should have used the default value when raising an error.")
	}

	val, err = config.GetFloat("bar", defaultValue)
	if err == nil {
		t.Errorf("Should have raised an type error when getting a non-existed value by GetString.")
	}
	if val != defaultValue {
		t.Errorf("Should have used the default value when raising an error.")
	}
}
