package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
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
	Desired() *int64
	Ready() *int64
	UpToDate() *int64
	Available() *int64
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

func (instance AnonymousAggregation) Desired() *int64 {
	return nil
}

func (instance AnonymousAggregation) Ready() *int64 {
	return nil
}

func (instance AnonymousAggregation) UpToDate() *int64 {
	return nil
}

func (instance AnonymousAggregation) Available() *int64 {
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

func (instance DeploymentAggregation) Desired() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "replicas"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance DeploymentAggregation) Ready() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "readyReplicas"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance DeploymentAggregation) UpToDate() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "updatedReplicas"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance DeploymentAggregation) Available() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "availableReplicas"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
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

func (instance DaemonSetAggregation) Desired() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "desiredNumberScheduled"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance DaemonSetAggregation) Ready() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "numberReady"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance DaemonSetAggregation) UpToDate() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "updatedNumberScheduled"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance DaemonSetAggregation) Available() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "numberAvailable"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
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

func (instance StatefulSetAggregation) Desired() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "replicas"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance StatefulSetAggregation) Ready() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "readyReplicas"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance StatefulSetAggregation) UpToDate() *int64 {
	if v, ok, err := unstructured.NestedInt64(instance.Object, "status", "updatedReplicas"); err != nil || !ok {
		return nil
	} else {
		return &v
	}
}

func (instance StatefulSetAggregation) Available() *int64 {
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

func (instance PodAggregation) Desired() *int64 {
	statuses, ok, err := unstructured.NestedSlice(instance.Object, "status", "containerStatuses")
	if err != nil || !ok || len(statuses) == 0 {
		return nil
	}
	l := int64(len(statuses))
	return &l
}

func (instance PodAggregation) Ready() *int64 {
	statuses, ok, err := unstructured.NestedSlice(instance.Object, "status", "containerStatuses")
	if err != nil || !ok || len(statuses) == 0 {
		return nil
	}
	var ready int64
	for _, candidate := range statuses {
		if candidate, ok := candidate.(map[string]interface{}); ok {
			if v, ok, err := unstructured.NestedBool(candidate, "ready"); err == nil && ok && v {
				ready++
			}
		}
	}
	return &ready
}

func (instance PodAggregation) UpToDate() *int64 {
	return nil
}

func (instance PodAggregation) Available() *int64 {
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
