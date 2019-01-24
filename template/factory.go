package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Factory interface {
	New(name string, code string) (Template, error)
	NewFromReader(name string, reader io.Reader) (Template, error)
	NewFromFile(file string) (Template, error)

	Must(name string, code string) Template
	MustFromReader(name string, reader io.Reader) Template
	MustFromFile(file string) Template
}

type FactoryImpl struct {
	FunctionProvider FunctionProvider
}

func (instance *FactoryImpl) new(name string, file *string, code string) (Template, error) {
	if delegate, err := newDelegate(name, code, instance.FunctionProvider); err != nil {
		return nil, err
	} else {
		return &Impl{
			sourceName:       name,
			sourceFile:       file,
			sourceCode:       code,
			functionProvider: instance.FunctionProvider,
			factory:          instance,
			delegate:         delegate,
		}, nil
	}
}

func (instance *FactoryImpl) New(name string, code string) (Template, error) {
	return instance.new(name, nil, code)
}

func (instance *FactoryImpl) NewFromReader(name string, reader io.Reader) (Template, error) {
	if content, err := ioutil.ReadAll(reader); err != nil {
		return nil, err
	} else {
		return instance.new(name, nil, string(content))
	}
}

func (instance *FactoryImpl) NewFromFile(file string) (Template, error) {
	if f, err := os.Open(file); os.IsNotExist(err) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("cannot read template from %s: %v", file, err)
	} else {
		//noinspection GoUnhandledErrorResult
		defer f.Close()
		if content, err := ioutil.ReadAll(f); err != nil {
			return nil, fmt.Errorf("cannot read template from %s: %v", file, err)
		} else {
			return instance.new(file, &file, string(content))
		}
	}
}

func (instance *FactoryImpl) Must(name string, code string) Template {
	if result, err := instance.New(name, code); err != nil {
		panic(err)
	} else {
		return result
	}
}

func (instance *FactoryImpl) MustFromReader(name string, reader io.Reader) Template {
	if result, err := instance.NewFromReader(name, reader); err != nil {
		panic(err)
	} else {
		return result
	}
}

func (instance *FactoryImpl) MustFromFile(file string) Template {
	if result, err := instance.NewFromFile(file); err != nil {
		panic(err)
	} else {
		return result
	}
}
