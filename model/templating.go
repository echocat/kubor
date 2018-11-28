package model

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"kubor/common"
	"path/filepath"
	"strings"
	"text/template"
)

type Templating struct {
	TemplateFilePattern []string `yaml:"templateFilePattern,omitempty" json:"templateFilePattern,omitempty"`
}

func newTemplating() Templating {
	return Templating{
		TemplateFilePattern: []string{
			"{{ .Root }}/kubernetes/templates/*.yaml",
		},
	}
}

func (instance Templating) TemplateFiles(data interface{}) ([]string, error) {
	return instance.renderFiles(instance.TemplateFilePattern, "template", data)
}

func (instance Templating) RenderedTemplatesProvider(data interface{}) (ContentProvider, error) {
	if files, err := instance.TemplateFiles(data); err != nil {
		return nil, err
	} else {
		i := 0
		return func() (string, []byte, error) {
			if i >= len(files) {
				return "", nil, io.EOF
			}
			buf := new(bytes.Buffer)
			file := files[i]
			i++
			if err := instance.renderFile(file, data, buf); err != nil {
				return file, nil, err
			}
			return file, buf.Bytes(), nil
		}, nil
	}
}

func (instance Templating) renderFile(file string, data interface{}, writer io.Writer) error {
	if content, err := ioutil.ReadFile(file); err != nil {
		return fmt.Errorf("cannot read template file '%s': %v", file, err)
	} else if tmpl, err := instance.newTemplate(file, string(content)); err != nil {
		return fmt.Errorf("cannot parse template file '%s': %v", file, err)
	} else if err := tmpl.Execute(writer, data); err != nil {
		return fmt.Errorf("cannot render template file '%s': %v", file, err)
	} else {
		return nil
	}
}

func (instance Templating) renderFiles(patterns []string, name string, data interface{}) ([]string, error) {
	var result []string
	for _, pattern := range patterns {
		atLeastOneMatchExpected := true
		if strings.HasPrefix(pattern, "?") {
			pattern = pattern[1:]
			atLeastOneMatchExpected = false
		}
		buf := new(bytes.Buffer)
		if tmpl, err := instance.newTemplate(pattern, pattern); err != nil {
			return nil, fmt.Errorf("cannot handle %s pattern '%s': %v", name, pattern, err)
		} else if err := tmpl.Execute(buf, data); err != nil {
			return nil, fmt.Errorf("cannot handle %s pattern: %v", name, err)
		} else if buf.Len() == 0 {
			// Ignore ... could happen if we use {{ if  }} clauses
		} else if matches, err := filepath.Glob(buf.String()); err != nil {
			return nil, fmt.Errorf("cannot handle %s pattern '%s' => '%s': %v", name, pattern, buf.String(), err)
		} else {
			if len(matches) <= 0 && atLeastOneMatchExpected {
				return nil, fmt.Errorf("there does not at least one %s file exist that matches '%s' => '%s'", name, pattern, buf.String())
			}
			result = append(result, matches...)
		}
	}
	return result, nil
}

func (instance Templating) newTemplate(name string, plain string) (*template.Template, error) {
	return common.NewTemplate(name, plain)
}
