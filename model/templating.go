package model

import (
	"bytes"
	"fmt"
	"github.com/levertonai/kubor/template/functions"
	"io"
	"path/filepath"
	"strings"
)

type Templating struct {
	TemplateFilePattern []string `yaml:"templateFilePattern,omitempty" json:"templateFilePattern,omitempty"`
}

func newTemplating() Templating {
	return Templating{
		TemplateFilePattern: []string{
			"?{{ .Root }}/kubernetes/templates/*.yml",
			"?{{ .Root }}/kubernetes/templates/*.yaml",
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
			if err := instance.RenderTemplateFile(file, data, buf); err != nil {
				return file, nil, err
			}
			return file, buf.Bytes(), nil
		}, nil
	}
}

func (instance Templating) RenderTemplateFile(file string, data interface{}, writer io.Writer) error {
	if tmpl, err := functions.DefaultTemplateFactory().NewFromFile(file); err != nil {
		return fmt.Errorf("cannot parse template file '%s': %v", file, err)
	} else if err := tmpl.Execute(data, writer); err != nil {
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
		if tmpl, err := functions.DefaultTemplateFactory().New(pattern, pattern); err != nil {
			return nil, fmt.Errorf("cannot handle %s pattern '%s': %v", name, pattern, err)
		} else if rendered, err := tmpl.ExecuteToString(data); err != nil {
			return nil, fmt.Errorf("cannot handle %s pattern: %v", name, err)
		} else if rendered == "" {
			// Ignore ... could happen if we use {{ if }} clauses
		} else if matches, err := filepath.Glob(rendered); err != nil {
			return nil, fmt.Errorf("cannot handle %s pattern '%s' => '%s': %v", name, pattern, rendered, err)
		} else {
			if len(matches) <= 0 && atLeastOneMatchExpected {
				return nil, fmt.Errorf("there does not at least one %s file exist that matches '%s' => '%s'", name, pattern, rendered)
			}
			result = append(result, matches...)
		}
	}
	return result, nil
}
