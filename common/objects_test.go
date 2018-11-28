package common

var (
	objectD = modelD{
		Value: "helloD",
	}
	objectC = modelC{
		Value: "helloC",
		D:     objectD,
	}
	objectB = modelB{
		Value: "helloB",
		C:     objectC,
	}
	objectA = modelA{
		Value: "helloA",
		AnInt: 666,
		B:     objectB,
	}
	mapD = map[string]interface{}{
		"Value": "helloD",
	}
	mapC = map[string]interface{}{
		"Value": "helloC",
		"D":     mapD,
	}
	mapB = map[string]interface{}{
		"Value": "helloB",
		"C":     mapC,
	}
	mapA = map[string]interface{}{
		"Value": "helloA",
		"B":     mapB,
		"AnInt": 666,
	}
)

type modelA struct {
	B     modelB
	Value string
	AnInt int
}

type modelB struct {
	C     modelC
	Value string
}

type modelC struct {
	D     modelD
	Value string
}

type modelD struct {
	Value string
}
