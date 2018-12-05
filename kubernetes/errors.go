package kubernetes

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
)

func OptimizeError(input error) error {
	if input == nil {
		return nil
	}
	if sErr, ok := input.(*errors.StatusError); ok {
		sErr.ErrStatus.Message = fmt.Sprintf("%d - %s: message=%s, reason=%s", sErr.ErrStatus.Code, sErr.ErrStatus.Status, sErr.ErrStatus.Message, sErr.ErrStatus.Reason)
	}
	return input
}
