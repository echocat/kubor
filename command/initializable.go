package command

import "github.com/levertonai/kubor/model"

type Initializable interface {
	Init(pf *model.ProjectFactory) error
}

var (
	initializables []Initializable
)

func RegisterInitializable(initializable Initializable) Initializable {
	initializables = append(initializables, initializable)
	return initializable
}

func Init(pf *model.ProjectFactory) error {
	for _, initializable := range initializables {
		if err := initializable.Init(pf); err != nil {
			return err
		}
	}
	return nil
}
