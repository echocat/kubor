package functions

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/pkg/errors"
)

var FuncDecodeBase64 = Function{
	Description: "Decodes base64 encoded string from given <source>.",
	Parameters: Parameters{{
		Name: "source",
	}},
	Returns: Return{
		Description: "String which was decoded from <source>.",
	},
}.MustWithFunc(func(source string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(source)
	return string(b), err
})

var FuncDecodeBase64Advanced = Function{
	Description: "Decodes base64 encoded string from given <source>.",
	Parameters: Parameters{{
		Name:        "type",
		Description: "Can be either 'url' or 'standard'.",
	}, {
		Name:        "raw",
		Description: "If true raw base64 is assumed.",
	}, {
		Name: "source",
	}},
	Returns: Return{
		Description: "String which was decoded from <source>.",
	},
}.MustWithFunc(func(aType string, raw bool, source string) (string, error) {
	enc, err := base64EncodingFor(aType, raw)
	if err != nil {
		return "", err
	}
	b, err := enc.DecodeString(source)
	return string(b), err
})

var FuncEncodeBase64 = Function{
	Description: "Encodes the given <source> as base64 encoded string",
	Parameters: Parameters{{
		Name: "source",
	}},
	Returns: Return{
		Description: "String which was encoded from <source>.",
	},
}.MustWithFunc(func(source string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(source)), nil
})

var FuncEncodeBase64Advanced = Function{
	Description: "Encodes the given <source> as base64 encoded string.",
	Parameters: Parameters{{
		Name:        "type",
		Description: "Can be either 'url' or 'standard'.",
	}, {
		Name:        "raw",
		Description: "If true raw base64 is assumed.",
	}, {
		Name: "source",
	}},
	Returns: Return{
		Description: "String which was encoded from <source>.",
	},
}.MustWithFunc(func(aType string, raw bool, source string) (string, error) {
	enc, err := base64EncodingFor(aType, raw)
	if err != nil {
		return "", err
	}
	return enc.EncodeToString([]byte(source)), nil
})

var FuncDecodeHex = Function{
	Description: "Decodes hex encoded string from given <source>.",
	Parameters: Parameters{{
		Name: "source",
	}},
	Returns: Return{
		Description: "String which was decoded from <source>.",
	},
}.MustWithFunc(func(source string) (string, error) {
	b, err := hex.DecodeString(source)
	return string(b), err
})

var FuncEncodeHex = Function{
	Description: "Encodes the given <source> as hex encoded string",
	Parameters: Parameters{{
		Name: "source",
	}},
	Returns: Return{
		Description: "String which was encoded from <source>.",
	},
}.MustWithFunc(func(source string) (string, error) {
	return hex.EncodeToString([]byte(source)), nil
})

var FuncsCodecs = Functions{
	"decodeBase64":    FuncDecodeBase64,
	"decodeBase64Adv": FuncDecodeBase64Advanced,
	"encodeBase64":    FuncEncodeBase64,
	"encodeBase64Adv": FuncEncodeBase64Advanced,
	"decodeHex":       FuncDecodeHex,
	"encodeHex":       FuncEncodeHex,
}
var CategoryCodecs = Category{
	Functions: FuncsCodecs,
}

func base64EncodingFor(aType string, raw bool) (*base64.Encoding, error) {
	switch aType {
	case "", "standard":
		if raw {
			return base64.RawStdEncoding, nil
		} else {
			return base64.StdEncoding, nil
		}
	case "url":
		if raw {
			return base64.RawURLEncoding, nil
		} else {
			return base64.URLEncoding, nil
		}
	}
	return nil, errors.Errorf("unknown base64 type: %s", aType)
}
