package kubernetes

import (
	"github.com/echocat/kubor/common"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
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
	case "pod":
		return PodAggregation{base}
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
	return Pbool(desired != nil && available != nil && *desired <= *available)
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
	return Pbool(desired != nil && available != nil && *desired <= *available)
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
	return Pbool(desired != nil && ready != nil && *desired <= *ready)
}

type PodAggregation struct {
	AnonymousAggregation
}

func (instance PodAggregation) Desired() *int32 {
	statuses := common.GetObjectPathValue(instance.Object, "status", "containerStatuses")
	if statuses == nil {
		return nil
	}
	vStatuses := common.SimplifyValue(reflect.ValueOf(statuses))
	if vStatuses.Kind() != reflect.Slice {
		return nil
	}
	return Pint32(int32(vStatuses.Len()))
}

func (instance PodAggregation) Ready() *int32 {
	statuses := common.GetObjectPathValue(instance.Object, "status", "containerStatuses")
	if statuses == nil {
		return nil
	}
	vStatuses := common.SimplifyValue(reflect.ValueOf(statuses))
	if vStatuses.Kind() != reflect.Slice {
		return nil
	}
	var ready int32
	numberOfContainers := vStatuses.Len()
	for i := 0; i < numberOfContainers; i++ {
		readyValue := common.GetObjectPathValue(vStatuses.Index(i).Interface(), "ready")
		switch v := readyValue.(type) {
		case *bool:
			if *v {
				ready++
			}
		case bool:
			if v {
				ready++
			}
		}
	}
	return Pint32(ready)
}

func (instance PodAggregation) UpToDate() *int32 {
	return nil
}

func (instance PodAggregation) Available() *int32 {
	return nil
}

func (instance PodAggregation) IsReady() *bool {
	desired := instance.Desired()
	ready := instance.Ready()
	return Pbool(desired != nil && ready != nil && *desired <= *ready)
}
