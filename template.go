package golf

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

// Regular Expressions to find out the extension syntax.
var reExtends = regexp.MustCompile(`{{ ?extends ["']?([^'"}']*)["']? ?}}`)
var reInclude = regexp.MustCompile(`{{ ?include ["']?([^"]*)["']? ?}}`)
var reTemplateTag = regexp.MustCompile(`{{ ?template \"([^"]*)" ?([^ ]*)? ?}}`)
var reDefineTag = regexp.MustCompile(`{{ ?define "([^"]*)" ?"?([a-zA-Z0-9]*)?"? ?}}`)

// TemplateLoader is the loader interface for templates.
type TemplateLoader interface {
	LoadTemplate(string) (string, error)
}

// FileSystemLoader is an implementation of TemplateLoader that loads templates
// from file system.
type FileSystemLoader struct {
	BaseDir string
}

// LoadTemplate loads a template from a file.
func (loader *FileSystemLoader) LoadTemplate(name string) (string, error) {
	f, err := os.Open(path.Join(loader.BaseDir, name))
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

// MapLoader is a implementation of TemplateLoader that loads templates from a
// map data structure.
type MapLoader map[string]string

// LoadTemplate loads a template from the map.
func (loader *MapLoader) LoadTemplate(name string) (string, error) {
	if src, ok := (*loader)[name]; ok {
		return src, nil
	}
	return "", Errorf("Could not find template " + name)
}

// TemplateManager handles the template loader and stores the function map.
type TemplateManager struct {
	Loader  TemplateLoader
	FuncMap map[string]interface{} //template.FuncMap
}

// Template type stands for a template.
type Template struct {
	Name string
	Src  string
}

// Render renders a template and writes it to the io.Writer interface.
func (t *TemplateManager) Render(w io.Writer, name string, data interface{}) error {
	stack := []*Template{}
	tplSrc, err := t.getSrc(name)
	if err != nil {
		return err
	}
	err = t.push(&stack, tplSrc, name)
	if err != nil {
		return err
	}
	tpl, err := t.assemble(stack)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if tpl == nil {
		return Errorf("Nil template named %s", name)
	}

	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

// RenderFromString renders a template from string.
func (t *TemplateManager) RenderFromString(w io.Writer, tplSrc string, data interface{}) error {
	stack := []*Template{}
	t.push(&stack, tplSrc, "_fromString")
	tpl, err := t.assemble(stack)
	if err != nil {
		return err
	}

	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (t *TemplateManager) push(stack *[]*Template, tplSrc string, name string) error {
	extendsMatches := reExtends.FindStringSubmatch(tplSrc)
	if len(extendsMatches) == 2 {
		src, err := t.getSrc(extendsMatches[1])
		t.push(stack, src, extendsMatches[1])
		if err != nil {
			return err
		}
		tplSrc = reExtends.ReplaceAllString(tplSrc, "")
	}
	Template := &Template{
		Name: name,
		Src:  tplSrc,
	}
	*stack = append((*stack), Template)
	return nil
}

func (t *TemplateManager) getSrc(name string) (string, error) {
	tplSrc, err := t.Loader.LoadTemplate(name)
	if err != nil {
		return "", err
	}

	if len(tplSrc) < 1 {
		return "", Errorf("Empty Template named %s", name)
	}
	return tplSrc, nil
}

func (t *TemplateManager) assemble(stack []*Template) (*template.Template, error) {
	blocks := map[string]string{}
	blockID := 0

	for _, Template := range stack {
		var errInReplace error
		Template.Src = reInclude.ReplaceAllStringFunc(Template.Src, func(raw string) string {
			parsed := reInclude.FindStringSubmatch(raw)
			templatePath := parsed[1]

			subTpl, err := t.Loader.LoadTemplate(templatePath)
			if err != nil {
				errInReplace = err
				return "[error]"
			}

			return subTpl
		})
		if errInReplace != nil {
			return nil, errInReplace
		}
	}

	for _, Template := range stack {

		Template.Src = reDefineTag.ReplaceAllStringFunc(Template.Src, func(raw string) string {
			parsed := reDefineTag.FindStringSubmatch(raw)
			blockName := fmt.Sprintf("BLOCK_%d", blockID)
			blocks[parsed[1]] = blockName

			blockID++
			return "{{ define \"" + blockName + "\" }}"
		})
	}

	var rootTemplate *template.Template

	for i, Template := range stack {
		Template.Src = reTemplateTag.ReplaceAllStringFunc(Template.Src, func(raw string) string {
			parsed := reTemplateTag.FindStringSubmatch(raw)
			origName := parsed[1]
			replacedName, ok := blocks[origName]
			dot := "."
			if len(parsed) == 3 && len(parsed[2]) > 0 {
				dot = parsed[2]
			}
			if ok {
				return fmt.Sprintf(`{{ template "%s" %s }}`, replacedName, dot)
			}
			return ""
		})
		var thisTemplate *template.Template

		if i == 0 {
			thisTemplate = template.New(Template.Name)
			rootTemplate = thisTemplate
		} else {
			thisTemplate = rootTemplate.New(Template.Name)
		}
		thisTemplate.Funcs(t.FuncMap)
		_, err := thisTemplate.Parse(Template.Src)
		if err != nil {
			return nil, err
		}
	}
	return rootTemplate, nil
}
