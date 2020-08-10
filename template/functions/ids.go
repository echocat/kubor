package functions

import (
	"github.com/google/uuid"
)

var FuncUuid = Function{
	Description: "Creates a new, random UUID",
}.MustWithFunc(func() uuid.UUID {
	return uuid.New()
})

var FuncsIds = Functions{
	"uuid": FuncUuid,
}
var CategoryIds = Category{
	Functions: FuncsIds,
}
