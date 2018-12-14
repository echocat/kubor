package model

import (
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io"
	"kubor/common"
	"kubor/log"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	// Values set using Load() method.
	GroupId           string              `yaml:"groupId,omitempty" json:"groupId,omitempty"`
	ArtifactId        string              `yaml:"artifactId" json:"artifactId"`
	Release           string              `yaml:"release,omitempty" json:"release,omitempty"`
	Templating        Templating          `yaml:"templating,omitempty" json:"templating,omitempty"`
	ConditionalValues []ConditionalValues `yaml:"values,omitempty" json:"values,omitempty"`

	// Values set using implicitly.
	Source  string            `yaml:"-" json:"-"`
	Root    string            `yaml:"-" json:"-"`
	Values  Values            `yaml:"-" json:"-"`
	Env     map[string]string `yaml:"-" json:"-"`
	Context string            `yaml:"-" json:"-"`
}

func newProject() Project {
	return Project{
		Templating:        newTemplating(),
		ConditionalValues: []ConditionalValues{},
		Values:            Values{},
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

func (instance *ProjectFactory) Create(contextName string) (Project, error) {
	var err error
	result := newProject()
	result.Context = contextName

	if f, err := os.Open(instance.source); os.IsNotExist(err) {
		if instance.sourceRequired {
			return Project{}, fmt.Errorf("could not find source file '%s'", instance.source)
		}
	} else if err != nil {
		return Project{}, fmt.Errorf("cannot open source file '%s': %v", instance.source, err)
	} else {
		//noinspection GoUnhandledErrorResult
		defer f.Close()
		if err := yaml.NewDecoder(f).Decode(&result); err != nil {
			return Project{}, fmt.Errorf("cannot read source file '%s': %v", instance.source, err)
		} else if err := result.Validate(); err != nil {
			return Project{}, fmt.Errorf("cannot read source file '%s': %v", instance.source, err)
		}
	}

	if result, err = instance.populateStage1(result); err != nil {
		return Project{}, err
	}
	if result, err = instance.populateStage2(result); err != nil {
		return Project{}, err
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

func (instance *ProjectFactory) populateStage1(input Project) (Project, error) {
	result := input
	result.Source = instance.source
	result.Root = filepath.Dir(result.Source)
	result.Values = instance.values
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
	return result, nil
}

func (instance *ProjectFactory) populateStage2(input Project) (Project, error) {
	result := input
	for _, candidate := range input.ConditionalValues {
		if ok, err := candidate.On.Matches(result); err != nil {
			return Project{}, err
		} else if ok {
			result.Values = result.Values.MergeWith(candidate.Values)
		}
	}
	return result, nil
}

func (instance *ProjectFactory) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "source",
			Usage:       "Specifies the location of the kubor source file.",
			Value:       ".kubor.yaml",
			EnvVar:      "KUBOR_SOURCE",
			Destination: &instance.source,
		},
		&cli.StringFlag{
			Name:        "groupId",
			Usage:       "If set it will overrides groupId from source file.",
			EnvVar:      "KUBOR_GROUP_ID",
			Destination: &instance.groupId,
		},
		&cli.StringFlag{
			Name:        "artifactId",
			Usage:       "If set it will overrides artifactId from source file.",
			EnvVar:      "KUBOR_ARTIFACT_ID",
			Destination: &instance.artifactId,
		},
		&cli.StringFlag{
			Name:        "release",
			Usage:       "If set it will overrides release from source file.",
			EnvVar:      "KUBOR_RELEASE",
			Destination: &instance.release,
		},
		&cli.BoolTFlag{
			Name:        "sourceRequired",
			Usage:       "If set to true the source file has to exist if not the execution will fail.",
			EnvVar:      "KUBOR_SOURCE_REQUIRED",
			Destination: &instance.sourceRequired,
		},
		&cli.GenericFlag{
			Name:  "value, v",
			Usage: "Specifies values which should be provided to the runtime. Format <name>=[<value>].",
			Value: &instance.values,
		},
	}
}
