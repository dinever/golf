package Golf

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
func (view *View) Render(loaderName string, filePath string, data interface{}) (string, error) {
	loader, ok := view.templateLoader[loaderName]
	if !ok {
		panic(errors.New("Template loader not found: " + loaderName))
	}
	var buf bytes.Buffer
	e := loader.Render(&buf, filePath, data)
	if e != nil {
		return "", e
	}
	return buf.String(), nil
}

// RenderFromString does the same thing as render but renders a template by
// indicating the template source from tplSrc.
func (view *View) RenderFromString(loaderName string, tplSrc string, data interface{}) (string, error) {
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
		panic(Errf("Template loader not fount: %s", loaderName))
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

// ErrorPage returns an error page.
func (view *View) ErrorPage(err error) string {
	text := `<!DOCTYPE HTML><html> <head> <meta http-equiv="content-type" content="text/html; charset=utf-8"> <title>{{.Title}}</title> <style type="text/css" media="screen">html,body{padding:0; margin:0; font-family: Tahoma; color: #34495e;}h1{color: #fff; margin: 0;}.container{max-width: 1220px; margin: 0 auto; padding: 0 20px;}#header{display: block; background-color: #3498db; height: 120px; width: 100%;}#title{padding: 40px 0;}</style> </head> <body style="background-color: #E4F9F5"> <div id="header"> <div class="container"> <div id="title"> <h1>{{.Title}}</h1> </div></div></div><div class="container"> <p>{{.Content}}</p></div></body></html>`
	tmpl := template.New("error")
	tmpl.Parse(text)
	var buf bytes.Buffer
	tmpl.Execute(&buf, map[string]interface{}{"Title": "Server Error", "Content": err.Error()})
	return buf.String()
}
