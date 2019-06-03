package core

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"path/filepath"
)

type (
	Template struct {
		LayoutPath  string
		IncludePath string

		//Templates *template.Template
		templates map[string]*template.Template
	}
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return errors.New("template doesn't exist")
	}
	//_ = tmpl.ExecuteTemplate(os.Stdout, "main", data)
	return tmpl.ExecuteTemplate(w, "main", data)
}

func NewTemplate(includePath string) *Template {
	return &Template{
		IncludePath: includePath,
		templates: map[string]*template.Template{},
	}
}

func (t *Template) WithLayoutPath(layoutPath string) *Template {
	t.LayoutPath = layoutPath
	return t
}

func (t *Template) Parse(e *echo.Echo) *Template {
	layoutFiles, err := filepath.Glob(t.LayoutPath + "*.html")
	if err != nil {
		e.Logger.Fatal(err)
	}

	includeFiles, err := filepath.Glob(t.IncludePath + "*.html")
	if err != nil {
		e.Logger.Fatal(err)
	}

	mainTemplate := template.New("main")

	mainTemplate, err = mainTemplate.Parse(`{{define "main" }}{{ template "base.html" . }}{{ end }}`)
	if err != nil {
		e.Logger.Fatal(err)
	}
	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		t.templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			e.Logger.Fatal(err)
		}
		t.templates[fileName] = template.Must(t.templates[fileName].ParseFiles(files...))
	}
	return t
}
