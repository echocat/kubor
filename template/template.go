package template

import (
	"bytes"
	"io"
	nt "text/template"
)

type Template interface {
	WithSourceFile(sourceFile string) (Template, error)

	Execute(data interface{}, target io.Writer) error
	ExecuteToString(data interface{}) (string, error)
	MustExecute(data interface{}, target io.Writer)
	MustExecuteToString(data interface{}) string

	GetSource() string
	GetSourceName() string
	GetSourceFile() *string
}

func newDelegate(name string, code string, functionProvider FunctionProvider) (*nt.Template, error) {
	if functions, err := functionProvider.GetFunctions(); err != nil {
		return nil, err
	} else if funcMap, err := functions.CreateDummyFuncMap(); err != nil {
		return nil, err
	} else {
		return nt.New(name).
			Funcs(funcMap).
			Option("missingkey=error").
			Parse(code)
	}
}

type Impl struct {
	sourceName string
	sourceFile *string
	sourceCode string

	factory          Factory
	functionProvider FunctionProvider
	delegate         *nt.Template
}

func (instance *Impl) WithSourceFile(sourceFile string) (Template, error) {
	return &Impl{
		sourceName:       instance.sourceName,
		sourceFile:       &sourceFile,
		sourceCode:       instance.sourceCode,
		functionProvider: instance.functionProvider,
		factory:          instance.factory,
		delegate:         instance.delegate,
	}, nil
}

func (instance *Impl) Execute(data interface{}, target io.Writer) error {
	if context, err := instance.newExecutionContext(data); err != nil {
		return err
	} else if clone, err := instance.delegate.Clone(); err != nil {
		return err
	} else if functions, err := instance.functionProvider.GetFunctions(); err != nil {
		return err
	} else if funcMap, err := functions.CreateFuncMap(context); err != nil {
		return err
	} else {
		return clone.
			Funcs(funcMap).
			Execute(target, data)
	}
}

func (instance *Impl) ExecuteToString(data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := instance.Execute(data, buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (instance *Impl) MustExecute(data interface{}, target io.Writer) {
	if err := instance.Execute(data, target); err != nil {
		panic(err)
	}
}

func (instance *Impl) MustExecuteToString(data interface{}) string {
	if result, err := instance.ExecuteToString(data); err != nil {
		panic(err)
	} else {
		return result
	}
}

func (instance *Impl) GetSourceName() string {
	return instance.sourceName
}

func (instance *Impl) GetSourceFile() *string {
	return instance.sourceFile
}

func (instance *Impl) GetSource() string {
	if instance.sourceFile != nil {
		return *instance.sourceFile
	}
	return instance.sourceName
}

func (instance *Impl) newExecutionContext(data interface{}) (ExecutionContext, error) {
	return &ExecutionContextImpl{
		Template: instance,
		Factory:  instance.factory,
		Data:     data,
	}, nil
}
