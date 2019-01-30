package kubernetes

import (
	"github.com/levertonai/kubor/common"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
	"time"
)

func IsReady(object runtime.Object) *bool {
	if us, ok := object.(*unstructured.Unstructured); ok {
		return NewAggregationFor(us).IsReady()
	}
	return nil
}

type Aggregation interface {
	Object
	GetDesired() *int32
	GetReady() *int32
	GetUpToDate() *int32
	GetAvailable() *int32
	IsReady() *bool
	GetStatus() *string
	GetAge() *time.Duration
}

func NewAggregationFor(object Object) Aggregation {
	generic := GenericAggregation{object}
	kind := object.GroupVersionKind()
	switch strings.ToLower(kind.Kind) {
	case "deployment":
		return DeploymentAggregation{generic}
	case "daemonset":
		return DaemonSetAggregation{generic}
	case "statefulset":
		return StatefulSetAggregation{generic}
	case "pod":
		return PodAggregation{generic}
	case "configmap":
		return ConfigMapAggregation{generic}
	}
	return generic
}

type GenericAggregation struct {
	Object
}

type v1MetaData interface {
	GetCreationTimestamp() v1.Time
}

func (instance GenericAggregation) GroupVersionKind() schema.GroupVersionKind {
	return NormalizeGroupVersionKind(instance.Object.GroupVersionKind())
}

func (instance GenericAggregation) GetAge() *time.Duration {
	if v1md, ok := instance.Object.(v1MetaData); ok {
		return common.PtimeDuration(time.Now().Sub(v1md.GetCreationTimestamp().Time))
	}
	return nil
}

func (instance GenericAggregation) Interface(segments ...string) interface{} {
	if u, ok := instance.Object.(Unstructured); ok {
		return common.GetObjectPathValue(u.UnstructuredContent(), segments...)
	}
	return common.GetObjectPathValue(instance.Object, segments...)
}

func (instance GenericAggregation) TryInt32(segments ...string) *int32 {
	return TryCastToInt32(instance.Interface(segments...))
}

func (instance GenericAggregation) TryString(segments ...string) *string {
	return TryCastToString(instance.Interface(segments...))
}

func (instance GenericAggregation) GetDesired() *int32 {
	return nil
}

func (instance GenericAggregation) GetReady() *int32 {
	return nil
}

func (instance GenericAggregation) GetUpToDate() *int32 {
	return nil
}

func (instance GenericAggregation) GetAvailable() *int32 {
	return nil
}

func (instance GenericAggregation) IsReady() *bool {
	return nil
}

func (instance GenericAggregation) GetStatus() *string {
	return nil
}

type DeploymentAggregation struct {
	GenericAggregation
}

func (instance DeploymentAggregation) GetDesired() *int32 {
	return instance.TryInt32("spec", "replicas")
}

func (instance DeploymentAggregation) GetReady() *int32 {
	return instance.TryInt32("spec", "readyReplicas")
}

func (instance DeploymentAggregation) GetUpToDate() *int32 {
	return instance.TryInt32("status", "updatedReplicas")
}

func (instance DeploymentAggregation) GetAvailable() *int32 {
	return instance.TryInt32("status", "availableReplicas")
}

func (instance DeploymentAggregation) IsReady() *bool {
	desired := instance.GetDesired()
	available := instance.GetAvailable()
	if desired == nil || available == nil {
		return nil
	}
	result := *desired <= *available
	return &result
}

func (instance DeploymentAggregation) GetStatus() *string {
	if ready := instance.IsReady(); ready == nil {
		return nil
	} else if *ready {
		return common.Pstring("Ready")
	} else {
		return common.Pstring("Not ready")
	}
}

type DaemonSetAggregation struct {
	GenericAggregation
}

func (instance DaemonSetAggregation) GetDesired() *int32 {
	return instance.TryInt32("status", "desiredNumberScheduled")
}

func (instance DaemonSetAggregation) GetReady() *int32 {
	return instance.TryInt32("status", "numberReady")
}

func (instance DaemonSetAggregation) GetUpToDate() *int32 {
	return instance.TryInt32("status", "updatedNumberScheduled")
}

func (instance DaemonSetAggregation) GetAvailable() *int32 {
	return instance.TryInt32("status", "numberAvailable")
}

func (instance DaemonSetAggregation) IsReady() *bool {
	desired := instance.GetDesired()
	available := instance.GetAvailable()
	if desired == nil || available == nil {
		return nil
	}
	result := *desired <= *available
	return &result
}

func (instance DaemonSetAggregation) GetStatus() *string {
	if ready := instance.IsReady(); ready == nil {
		return nil
	} else if *ready {
		return common.Pstring("Ready")
	} else {
		return common.Pstring("Not ready")
	}
}

type StatefulSetAggregation struct {
	GenericAggregation
}

func (instance StatefulSetAggregation) GetDesired() *int32 {
	return instance.TryInt32("spec", "replicas")
}

func (instance StatefulSetAggregation) GetReady() *int32 {
	return instance.TryInt32("status", "readyReplicas")
}

func (instance StatefulSetAggregation) GetUpToDate() *int32 {
	return instance.TryInt32("status", "updatedReplicas")
}

func (instance StatefulSetAggregation) GetAvailable() *int32 {
	return nil
}

func (instance StatefulSetAggregation) IsReady() *bool {
	desired := instance.GetDesired()
	ready := instance.GetReady()
	if desired == nil || ready == nil {
		return nil
	}
	result := *desired <= *ready
	return &result
}

func (instance StatefulSetAggregation) GetStatus() *string {
	if ready := instance.IsReady(); ready == nil {
		return nil
	} else if *ready {
		return common.Pstring("Ready")
	} else {
		return common.Pstring("Not ready")
	}
}

type PodAggregation struct {
	GenericAggregation
}

func (instance PodAggregation) IsReady() *bool {
	phase := instance.TryString("status", "phase")
	return common.Pbool(phase != nil && strings.ToLower(*phase) == "running")
}

func (instance PodAggregation) GetStatus() *string {
	phase := instance.TryString("status", "phase")
	reason := instance.TryString("status", "reason")
	if phase != nil && *phase != "" && reason != nil && *reason != "" {
		return common.Pstring(*phase + ": " + *reason)
	} else if phase != nil && *phase != "" {
		return phase
	} else if reason != nil && *reason != "" {
		return reason
	} else {
		return nil
	}
}

type ConfigMapAggregation struct {
	GenericAggregation
}

func (instance ConfigMapAggregation) GetStatus() *string {
	return common.Pstring("Exists")
}
