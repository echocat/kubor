package kubernetes

import (
	"github.com/echocat/kubor/common"
	v1 "k8s.io/api/core/v1"
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

func StateOf(object runtime.Object) *State {
	if us, ok := object.(*unstructured.Unstructured); ok {
		return NewAggregationFor(us).State()
	}
	return nil
}

type Aggregation interface {
	Desired() *int32
	Ready() *int32
	UpToDate() *int32
	Available() *int32
	IsReady() *bool
	State() *State
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

func (instance AnonymousAggregation) State() *State {
	return nil
}

type DeploymentAggregation struct {
	AnonymousAggregation
}

func (instance DeploymentAggregation) Desired() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "replicas"))
}

func (instance DeploymentAggregation) Ready() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "readyReplicas"))
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

func (instance DeploymentAggregation) State() *State {
	return nil
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

func (instance DaemonSetAggregation) State() *State {
	return nil
}

type StatefulSetAggregation struct {
	AnonymousAggregation
}

func (instance StatefulSetAggregation) Desired() *int32 {
	return TryCastToInt32(common.GetObjectPathValue(instance.Object, "status", "replicas"))
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

func (instance StatefulSetAggregation) State() *State {
	return nil
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

func (instance PodAggregation) State() *State {
	plain, found, err := unstructured.NestedString(instance.Object, "status", "phase")
	if !found || err != nil {
		return PState(StateUnknown)
	}
	switch v1.PodPhase(plain) {
	case v1.PodPending:
		return PState(StatePending)
	case v1.PodRunning:
		return PState(StateRunning)
	case v1.PodSucceeded:
		return PState(StateSucceeded)
	case v1.PodFailed:
		return PState(StateFailed)
	default:
		return PState(StateUnknown)
	}
}

type State uint8

const (
	StateUnknown   = State(0)
	StatePending   = State(1)
	StateRunning   = State(2)
	StateSucceeded = State(3)
	StateFailed    = State(4)
)

func PState(in State) *State {
	return &in
}

func (instance State) IsActive() bool {
	switch instance {
	case StatePending, StateRunning:
		return true
	default:
		return false
	}
}

func (instance State) IsDone() bool {
	switch instance {
	case StateSucceeded, StateFailed:
		return true
	default:
		return false
	}
}
