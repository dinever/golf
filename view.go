package golf

import (
	"bytes"
	"errors"
	"html/template"
)

// View handles templates rendering
type View struct {
	FuncMap template.FuncMap

	// A view may have multiple template managers, e.g., one for the admin panel,
	// another one for the user end.
	templateLoader map[string]*TemplateManager
}

// NewView creates a new View instance.
func NewView() *View {
	view := new(View)
	view.templateLoader = make(map[string]*TemplateManager)
	view.FuncMap = make(template.FuncMap)
	view.FuncMap["Html"] = func(s string) template.HTML { return template.HTML(s) }
	return view
}

// Render renders a template by indicating the file path and the name of the template loader.
func (view *View) Render(loaderName string, filePath string, data map[string]interface{}) (string, error) {
	loader, ok := view.templateLoader[loaderName]
	if !ok {
		panic(errors.New("Template loader not found: " + loaderName))
	}
	var buf bytes.Buffer
	err := loader.Render(&buf, filePath, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderFromString does the same thing as render but renders a template by
// indicating the template source from tplSrc.
func (view *View) RenderFromString(loaderName string, tplSrc string, data map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	// If loaderName is not indicated, use the default template library of Go, no syntax like
	// `extends` or `include` will be supported.
	if loaderName == "" {
		tmpl := template.New("error")
		tmpl.Parse(tplSrc)
		e := tmpl.Execute(&buf, data)
		return buf.String(), e
	}
	loader, ok := view.templateLoader[loaderName]
	if !ok {
		panic(Errorf("Template loader not fount: %s", loaderName))
	}
	e := loader.Render(&buf, tplSrc, data)
	if e != nil {
		return "", e
	}
	return buf.String(), nil
}

// SetTemplateLoader sets the template loader by indicating the name and the path.
func (view *View) SetTemplateLoader(name string, path string) {
	loader := &TemplateManager{
		Loader: &FileSystemLoader{
			BaseDir: path,
		},
		FuncMap: view.FuncMap,
	}
	view.templateLoader[name] = loader
}
