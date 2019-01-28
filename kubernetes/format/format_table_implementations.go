package format

import (
	appV1 "k8s.io/api/apps/v1"
	appV1beta1 "k8s.io/api/apps/v1beta1"
	appV1beta2 "k8s.io/api/apps/v1beta2"
	extensionsV1beta2 "k8s.io/api/extensions/v1beta1"
)

var _ = MustRegisterTableBasedObjectFormatOf("Deployments",
	appV1beta1.SchemeGroupVersion.WithKind("deployment"),
	appV1beta2.SchemeGroupVersion.WithKind("deployment"),
	appV1.SchemeGroupVersion.WithKind("deployment"),
	extensionsV1beta2.SchemeGroupVersion.WithKind("deployment"),
).
	WithColumn("Namespace", ObjectPathFormatter("metadata", "namespace")).
	WithColumn("Name", ObjectPathFormatter("metadata", "name")).
	WithColumn("Is ready", AggregationFormatter(AggregationToIsReady)).
	WithColumn("Desired", AggregationFormatter(AggregationToDesired)).
	WithColumn("Ready", AggregationFormatter(AggregationToReady)).
	WithColumn("Up to date", AggregationFormatter(AggregationToUpToDate)).
	WithColumn("Available", AggregationFormatter(AggregationToAvailable))

var _ = MustRegisterTableBasedObjectFormatOf("DaemonSets",
	appV1beta2.SchemeGroupVersion.WithKind("daemonset"),
	appV1.SchemeGroupVersion.WithKind("daemonset"),
	extensionsV1beta2.SchemeGroupVersion.WithKind("daemonset"),
).
	WithColumn("Namespace", ObjectPathFormatter("metadata", "namespace")).
	WithColumn("Name", ObjectPathFormatter("metadata", "name")).
	WithColumn("Is ready", AggregationFormatter(AggregationToIsReady)).
	WithColumn("Desired", AggregationFormatter(AggregationToDesired)).
	WithColumn("Ready", AggregationFormatter(AggregationToReady)).
	WithColumn("Up to date", AggregationFormatter(AggregationToUpToDate)).
	WithColumn("Available", AggregationFormatter(AggregationToAvailable))

var _ = MustRegisterTableBasedObjectFormatOf("StatefulSets",
	appV1beta1.SchemeGroupVersion.WithKind("statefulset"),
	appV1beta2.SchemeGroupVersion.WithKind("statefulset"),
	appV1.SchemeGroupVersion.WithKind("statefulset"),
	extensionsV1beta2.SchemeGroupVersion.WithKind("statefulset"),
).
	WithColumn("Namespace", ObjectPathFormatter("metadata", "namespace")).
	WithColumn("Name", ObjectPathFormatter("metadata", "name")).
	WithColumn("Is ready", AggregationFormatter(AggregationToIsReady)).
	WithColumn("Desired", AggregationFormatter(AggregationToDesired)).
	WithColumn("Ready", AggregationFormatter(AggregationToReady)).
	WithColumn("Up to date", AggregationFormatter(AggregationToUpToDate))

var _ = MustRegisterTableBasedObjectFormatOf("").
	WithColumn("Namespace", ObjectPathFormatter("metadata", "namespace")).
	WithColumn("Name", ObjectPathFormatter("metadata", "name"))
