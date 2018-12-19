package functions

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"kubor/template"
	"os"
	"strings"
)

var _ = Register(Function{
	Name:        "decodeYaml",
	Category:    "serialization",
	Description: "Decodes YAML from given <source>.",
	Parameters: Parameters{{
		Name: "source",
	}},
	Returns: Return{
		Description: "Object which was decoded from <source>.",
	},
	Func: func(context template.ExecutionContext, source string) (interface{}, error) {
		reader := strings.NewReader(source)
		var result interface{}
		if err := yaml.NewDecoder(reader).Decode(&result); err != nil {
			return nil, fmt.Errorf("cannot decode yaml referenced in '%s': %v", context.GetTemplate().GetSource(), err)
		} else {
			return result, nil
		}
	},
}, Function{
	Name:        "decodeJson",
	Category:    "serialization",
	Description: "Decodes JSON from given <source>.",
	Parameters: Parameters{{
		Name: "source",
	}},
	Returns: Return{
		Description: "Object which was decoded from <source>.",
	},
	Func: func(context template.ExecutionContext, source string) (interface{}, error) {
		reader := strings.NewReader(source)
		var result interface{}
		if err := json.NewDecoder(reader).Decode(&result); err != nil {
			return nil, fmt.Errorf("cannot decode json referenced in '%s': %v", context.GetTemplate().GetSource(), err)
		} else {
			return result, nil
		}
	},
}, Function{
	Name:        "decodeYamlFromFile",
	Category:    "serialization",
	Description: "Decodes YAML from given <file>.",
	Parameters: Parameters{{
		Name: "file",
	}},
	Returns: Return{
		Description: "Object which was decoded from <file>.",
	},
	Func: func(context template.ExecutionContext, file string) (interface{}, error) {
		if resolved, err := resolvePathOfContext(context, file); err != nil {
			return nil, err
		} else if f, err := os.Open(resolved); os.IsNotExist(err) {
			return nil, fmt.Errorf("file '%s' referenced in '%s' does not exist", resolved, context.GetTemplate().GetSource())
		} else {
			//noinspection GoUnhandledErrorResult
			defer f.Close()
			var result interface{}
			if err := yaml.NewDecoder(f).Decode(&result); err != nil {
				return nil, fmt.Errorf("cannot decode yaml from '%s' referenced in '%s': %v", resolved, context.GetTemplate().GetSource(), err)
			} else {
				return result, nil
			}
		}
	},
}, Function{
	Name:        "decodeJsonFromFile",
	Category:    "serialization",
	Description: "Decodes JSON from given <file>.",
	Parameters: Parameters{{
		Name: "file",
	}},
	Returns: Return{
		Description: "Object which was decoded from <file>.",
	},
	Func: func(context template.ExecutionContext, file string) (interface{}, error) {
		if resolved, err := resolvePathOfContext(context, file); err != nil {
			return nil, err
		} else if f, err := os.Open(resolved); os.IsNotExist(err) {
			return nil, fmt.Errorf("file '%s' referenced in '%s' does not exist", resolved, context.GetTemplate().GetSource())
		} else {
			//noinspection GoUnhandledErrorResult
			defer f.Close()
			var result interface{}
			if err := json.NewDecoder(f).Decode(&result); err != nil {
				return nil, fmt.Errorf("cannot decode yaml from '%s' referenced in '%s': %v", resolved, context.GetTemplate().GetSource(), err)
			} else {
				return result, nil
			}
		}
	},
})
