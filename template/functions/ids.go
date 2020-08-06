package functions

import (
	"github.com/google/uuid"
)

var FuncNewUuid = Function{
	Description: "Creates a new, random UUID",
}.MustWithFunc(func() uuid.UUID {
	return uuid.New()
})

var FuncsIds = Functions{
	"newUuid": FuncNewUuid,
}
var CategoryIds = Category{
	Functions: FuncsIds,
}
