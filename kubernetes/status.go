package kubernetes

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"kubor/common"
	"strings"
)

func IsReady(object runtime.Object) *bool {
	if us, ok := object.(*unstructured.Unstructured); ok {
		return NewAggregationFor(us).IsReady()
	}
	return nil
}

type Aggregation interface {
	Desired() *int32
	Ready() *int32
	UpToDate() *int32
	Available() *int32
	IsReady() *bool
}

func NewAggregationFor(object *unstructured.Unstructured) Aggregation {
	base := AnonymousAggregation{object}
	kind := object.GetObjectKind().GroupVersionKind()
	switch strings.ToLower(kind.Kind) {
	case "deployment":
		return DeploymentAggregation{base}
	case "daemonset":
		return DaemonSetAggregation{base}
	case "statefulset":
		return StatefulSetAggregation{base}
	}
	return base
}

type AnonymousAggregation struct {
	*unstructured.Unstructured
}

func (instance AnonymousAggregation) Desired() *int32 {
	return nil
}

func (instance AnonymousAggregation) Ready() *int32 {
	return nil
}

func (instance AnonymousAggregation) UpToDate() *int32 {
	return nil
}

func (instance AnonymousAggregation) Available() *int32 {
	return nil
}

func (instance AnonymousAggregation) IsReady() *bool {
	return nil
}

type DeploymentAggregation struct {
	AnonymousAggregation
}

func (instance DeploymentAggregation) Desired() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "spec", "replicas"))
}

func (instance DeploymentAggregation) Ready() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "spec", "readyReplicas"))
}

func (instance DeploymentAggregation) UpToDate() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "updatedReplicas"))
}

func (instance DeploymentAggregation) Available() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "availableReplicas"))
}

func (instance DeploymentAggregation) IsReady() *bool {
	desired := instance.Desired()
	available := instance.Available()
	if desired == nil || available == nil {
		return nil
	}
	result := *desired <= *available
	return &result
}

type DaemonSetAggregation struct {
	AnonymousAggregation
}

func (instance DaemonSetAggregation) Desired() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "desiredNumberScheduled"))
}

func (instance DaemonSetAggregation) Ready() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "numberReady"))
}

func (instance DaemonSetAggregation) UpToDate() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "updatedNumberScheduled"))
}

func (instance DaemonSetAggregation) Available() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "numberAvailable"))
}

func (instance DaemonSetAggregation) IsReady() *bool {
	desired := instance.Desired()
	available := instance.Available()
	if desired == nil || available == nil {
		return nil
	}
	result := *desired <= *available
	return &result
}

type StatefulSetAggregation struct {
	AnonymousAggregation
}

func (instance StatefulSetAggregation) Desired() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "spec", "replicas"))
}

func (instance StatefulSetAggregation) Ready() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "readyReplicas"))
}

func (instance StatefulSetAggregation) UpToDate() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "updatedReplicas"))
}

func (instance StatefulSetAggregation) Available() *int32 {
	return nil
}

func (instance StatefulSetAggregation) IsReady() *bool {
	desired := instance.Desired()
	ready := instance.Ready()
	if desired == nil || ready == nil {
		return nil
	}
	result := *desired <= *ready
	return &result
}
