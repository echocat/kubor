package model

import (
	"fmt"
	"github.com/echocat/kubor/common"
	"github.com/echocat/kubor/log"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	// Values set using Load() method.
	GroupId           string              `yaml:"groupId,omitempty" json:"groupId,omitempty"`
	ArtifactId        string              `yaml:"artifactId" json:"artifactId"`
	Release           string              `yaml:"release,omitempty" json:"release,omitempty"`
	Stages            Stages              `yaml:"stages,omitempty" json:"stages,omitempty"`
	Templating        Templating          `yaml:"templating,omitempty" json:"templating,omitempty"`
	ConditionalValues []ConditionalValues `yaml:"values,omitempty" json:"values,omitempty"`
	Labels            Labels              `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations       Annotations         `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	Validation        Validation          `yaml:"validation,omitempty" json:"validation,omitempty"`

	// Values set using implicitly.
	Source  string            `yaml:"-" json:"-"`
	Root    string            `yaml:"-" json:"-"`
	Values  Values            `yaml:"-" json:"-"`
	Env     map[string]string `yaml:"-" json:"-"`
	Context string            `yaml:"-" json:"-"`
}

func newProject() *Project {
	return &Project{
		Templating:        newTemplating(),
		ConditionalValues: []ConditionalValues{},
		Values:            Values{},
		Labels:            newLabels(),
		Annotations:       newAnnotations(),
	}
}

func (instance Project) Validate() error {
	if strings.TrimSpace(instance.ArtifactId) == "" {
		return fmt.Errorf("artifactId should not be empty")
	}
	return nil
}

func (instance *Project) Save() error {
	dir := filepath.Dir(instance.Source)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create parent directory for source file '%s': %v", instance.Source, err)
	} else if f, err := os.OpenFile(instance.Source, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return fmt.Errorf("cannot save source file '%s': %v", instance.Source, err)
	} else {
		//noinspection GoUnhandledErrorResult
		defer f.Close()
		encoder := yaml.NewEncoder(f)
		if err := encoder.Encode(instance); err != nil {
			return fmt.Errorf("cannot save source file '%s': %v", instance.Source, err)
		}
		return nil
	}
}

func (instance Project) RenderedTemplatesProvider() (ContentProvider, error) {
	return instance.Templating.RenderedTemplatesProvider(instance)
}

func (instance Project) RenderedTemplateFile(file string, writer io.Writer) error {
	return instance.Templating.RenderTemplateFile(file, instance, writer)
}

type ProjectFactory struct {
	source         string
	sourceRequired bool
	values         Values
	artifactId     string
	groupId        string
	release        string
}

func NewProjectFactory() *ProjectFactory {
	return &ProjectFactory{}
}

func (instance *ProjectFactory) Create(context string) (*Project, error) {
	result := newProject()
	result.Context = context

	if source, err := instance.resolveSource(); os.IsNotExist(err) {
		if instance.sourceRequired {
			return nil, fmt.Errorf("could not find source file '%s'", instance.source)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cannot open source file '%s': %v", instance.source, err)
	} else if f, err := os.Open(source); err != nil {
		return nil, fmt.Errorf("cannot open source file '%s': %v", source, err)
	} else {
		//noinspection GoUnhandledErrorResult
		defer f.Close()
		if err := yaml.NewDecoder(f).Decode(&result); err != nil {
			return nil, fmt.Errorf("cannot read source file '%s': %v", source, err)
		} else if err := result.Validate(); err != nil {
			return nil, fmt.Errorf("cannot read source file '%s': %v", source, err)
		}

		if result, err = instance.populateStage1(source, result); err != nil {
			return nil, err
		}
		if result, err = instance.populateStage2(result); err != nil {
			return nil, err
		}
	}

	if log.IsDebugEnabled() {
		name := result.ArtifactId
		l := log.
			WithField("source", result.Source).
			WithField("artifactId", result.ArtifactId)

		if result.GroupId != "" {
			name = fmt.Sprintf("%s:%s", result.GroupId, name)
			l = l.WithField("groupId", result.GroupId)
		}
		if result.Release != "" {
			name = fmt.Sprintf("%s:%s", name, result.Release)
			l = l.WithField("release", result.Release)
		}

		for k, v := range result.Values {
			l = l.WithField("value."+k, v)
		}

		l.Debug("Project %s", name)
	}

	return result, nil
}

func (instance *ProjectFactory) resolveSource() (string, error) {
	if _, err := os.Stat(instance.source); err == nil {
		return instance.source, nil
	} else if os.IsNotExist(err) {
		var alternative string
		for i := len(instance.source) - 1; i >= 0 && !os.IsPathSeparator(instance.source[i]); i-- {
			if instance.source[i] == '.' {
				switch instance.source[i:] {
				case ".yml":
					alternative = instance.source[:i] + ".yaml"
				case ".yaml":
					alternative = instance.source[:i] + ".yml"
				}
				break
			}
		}
		if alternative == "" {
			return "", err
		}
		if _, aErr := os.Stat(alternative); aErr == nil {
			return alternative, nil
		} else if os.IsNotExist(aErr) {
			return "", err
		} else {
			return "", aErr
		}
	} else {
		return "", err
	}
}

func (instance *ProjectFactory) populateStage1(source string, input *Project) (*Project, error) {
	result := *input
	result.Source = source
	result.Root = filepath.Dir(result.Source)
	result.Values = Values{}
	for k, v := range instance.values {
		result.Values[k] = v
	}
	if instance.groupId != "" {
		result.GroupId = instance.groupId
	}
	if instance.artifactId != "" {
		result.ArtifactId = instance.artifactId
	}
	if instance.release != "" {
		result.Release = instance.release
	}
	if result.Release == "" {
		result.Release = "latest"
	}
	result.Env = common.Environ()
	return &result, nil
}

func (instance *ProjectFactory) populateStage2(input *Project) (*Project, error) {
	result := *input
	for _, candidate := range input.ConditionalValues {
		if ok, err := candidate.On.Matches(result); err != nil {
			return nil, err
		} else if ok {
			result.Values = result.Values.MergeWith(candidate.Values)
			result.Values = result.Values.MergeWith(instance.values)
		}
	}
	return &result, nil
}

func (instance *ProjectFactory) ConfigureFlags(hf common.HasFlags) {
	hf.Flag("source", "Specifies the location of the kubor source file.").
		Default(".kubor.yml").
		Envar("KUBOR_SOURCE").
		PlaceHolder("<source file>").
		StringVar(&instance.source)
	hf.Flag("groupId", "If set it will overrides groupId from source file.").
		Envar("KUBOR_GROUP_ID").
		PlaceHolder("<groupId>").
		StringVar(&instance.groupId)
	hf.Flag("artifactId", "If set it will overrides artifactId from source file.").
		Envar("KUBOR_ARTIFACT_ID").
		PlaceHolder("<artifactId>").
		StringVar(&instance.artifactId)
	hf.Flag("release", "If set it will overrides release from source file.").
		Envar("KUBOR_RELEASE").
		PlaceHolder("<release>").
		StringVar(&instance.release)
	hf.Flag("sourceRequired", "If set to true the source file has to exist if not the execution will fail.").
		Default(fmt.Sprint(instance.sourceRequired)).
		Envar("KUBOR_SOURCE_REQUIRED").
		BoolVar(&instance.sourceRequired)
	hf.Flag("value", "Specifies values which should be provided to the runtime.").
		Short('v').
		PlaceHolder("<name>=[<value>]").
		SetValue(&instance.values)
}
