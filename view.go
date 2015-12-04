package Golf

import (
	"bytes"
	"html/template"
)

type View struct {
	BaseDir string
	FuncMap template.FuncMap
}

func NewView(baseDir string) *View {
	view := new(View)
	view.BaseDir = baseDir
	view.FuncMap = make(template.FuncMap)
	view.FuncMap["Html"] = func(s string) template.HTML { return template.HTML(s) }
	return view
}

func (view *View) getTemplateFromPath(filePath string) (*template.Template, error) {
	t := template.New(filePath)
	t.Funcs(view.FuncMap)
	t, e := t.ParseFiles(filePath)
	if e != nil {
		return nil, e
	}
	return t, nil
}

func (view *View) getTemplateFromString(content string) (*template.Template, error) {
	t, e := template.New("").Parse(content)
	if e != nil {
		return nil, e
	}
	return t, nil
}

func (view *View) Render(filePath string, data map[string]interface{}) (string, error) {
	t, e := view.getTemplateFromPath(filePath)
	if e != nil {
		return "", e
	}
	var buf bytes.Buffer
	e = t.Execute(&buf, data)
	if e != nil {
		return "", e
	}
	return buf.String(), nil
}

func (view *View) RenderFromString(content string, data map[string]interface{}) (string, error) {
	t, e := view.getTemplateFromString(content)
	if e != nil {
		return "", e
	}
	var buf bytes.Buffer
	e = t.Execute(&buf, data)
	if e != nil {
		return "", e
	}
	return buf.String(), nil
}
