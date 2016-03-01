package Golf

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
var re_extends *regexp.Regexp = regexp.MustCompile(`{{ ?extends ["']?([^'"}']*)["']? ?}}`)
var re_include *regexp.Regexp = regexp.MustCompile(`{{ ?include ["']?([^"]*)["']? ?}}`)
var re_templateTag *regexp.Regexp = regexp.MustCompile(`{{ ?template \"([^"]*)" ?([^ ]*)? ?}}`)
var re_defineTag *regexp.Regexp = regexp.MustCompile(`{{ ?define "([^"]*)" ?"?([a-zA-Z0-9]*)?"? ?}}`)

type TemplateLoader interface {
	LoadTemplate(string) (string, error)
}

type FileSystemLoader struct {
	BaseDir string
}

func (loader *FileSystemLoader) LoadTemplate(name string) (string, error) {
	f, err := os.Open(path.Join(loader.BaseDir, name))
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(f)
	return string(b), err
}

type MapLoader map[string]string

func (loader *MapLoader) LoadTemplate(name string) (string, error) {
	if src, ok := (*loader)[name]; ok {
		return src, nil
	}
	return "", Errf("Could not find template " + name)
}

type TemplateManager struct {
	Loader  TemplateLoader
	FuncMap map[string]interface{} //template.FuncMap
}

type Template struct {
	Name string
	Src  string
}

func (t *TemplateManager) Render(w io.Writer, name string, data interface{}) error {
	stack := []*Template{}
	tplSrc, err := t.getSrc(name)
	err = t.push(&stack, tplSrc, name)
	tpl, err := t.assemble(stack)
	if err != nil {
		return err
	}

	if tpl == nil {
		return Errf("Nil template named %s", name)
	}

	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

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
	extendsMatches := re_extends.FindStringSubmatch(tplSrc)
	if len(extendsMatches) == 2 {
		src, err := t.getSrc(extendsMatches[1])
		t.push(stack, src, extendsMatches[1])
		if err != nil {
			return err
		}
		tplSrc = re_extends.ReplaceAllString(tplSrc, "")
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
		return "", Errf("Empty Template named %s", name)
	}
	return tplSrc, nil
}

func (t *TemplateManager) assemble(stack []*Template) (*template.Template, error) {
	blocks := map[string]string{}
	blockId := 0

	for _, Template := range stack {
		var errInReplace error = nil
		Template.Src = re_include.ReplaceAllStringFunc(Template.Src, func(raw string) string {
			parsed := re_include.FindStringSubmatch(raw)
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

		Template.Src = re_defineTag.ReplaceAllStringFunc(Template.Src, func(raw string) string {
			parsed := re_defineTag.FindStringSubmatch(raw)
			blockName := fmt.Sprintf("BLOCK_%d", blockId)
			blocks[parsed[1]] = blockName

			blockId += 1
			return "{{ define \"" + blockName + "\" }}"
		})
	}

	var rootTemplate *template.Template

	for i, Template := range stack {
		Template.Src = re_templateTag.ReplaceAllStringFunc(Template.Src, func(raw string) string {
			parsed := re_templateTag.FindStringSubmatch(raw)
			origName := parsed[1]
			replacedName, ok := blocks[origName]
			dot := "."
			if len(parsed) == 3 && len(parsed[2]) > 0 {
				dot = parsed[2]
			}
			if ok {
				return fmt.Sprintf(`{{ template "%s" %s }}`, replacedName, dot)
			} else {
				return ""
			}
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
