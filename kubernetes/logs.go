package kubernetes

import (
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"strings"
	"time"
)

var (
	errLogsTemporaryProblem = fmt.Errorf("temporary problem with logs")
)

func LogProviderFor(runtime Runtime, object kruntime.Object, container string) LogProvider {
	if us, ok := object.(*unstructured.Unstructured); ok {
		gvk := us.GetObjectKind().GroupVersionKind()
		switch strings.ToLower(gvk.Kind) {
		case "pod":
			return &PodLogProvider{runtime, us, container}
		}
		return nil
	}
	return nil
}

type WriterProvider func() (io.WriteCloser, error)

func PrintLogs(ctx context.Context, using LogProvider, to WriterProvider) (err error) {
	writer, tErr := to()
	if tErr != nil {
		return fmt.Errorf("cannot open target to write logs to: %w", tErr)
	}
	defer func() { _ = writer.Close() }()

	var cr io.ReadCloser
	var closed bool
	go func() {
		<-ctx.Done()
		closed = true
	}()
	step := func() (err error) {
		if cr, err = using.Open(); err == errLogsTemporaryProblem {
			return nil
		} else if err != nil {
			return
		}
		defer func() { _ = cr.Close() }()
		if _, err = io.Copy(writer, cr); err != nil {
			return
		}
		return
	}
	for !closed {
		err = step()
		if err != nil {
			return
		}
		if !closed {
			time.Sleep(time.Millisecond * 500)
		}
	}
	return
}

type LogProvider interface {
	Open() (io.ReadCloser, error)
}

type PodLogProvider struct {
	runtime   Runtime
	object    *unstructured.Unstructured
	container string
}

func (instance *PodLogProvider) Open() (io.ReadCloser, error) {
	client, err := instance.runtime.NewRestClient(instance.object.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	req := client.Get().Namespace(instance.object.GetNamespace()).
		Name(instance.object.GetName()).
		Resource("pods").
		SubResource("log").
		VersionedParams(&v1.PodLogOptions{
			Follow:    true,
			Container: instance.container,
		}, scheme.ParameterCodec)

	result, err := req.Stream()
	if err != nil {
		if ass, ok := err.(errors.APIStatus); ok {
			status := ass.Status()
			if status.Code == 400 &&
				status.Details == nil &&
				strings.HasSuffix(status.Message, "is waiting to start: ContainerCreating") {
				return nil, errLogsTemporaryProblem
			}
		}
		return nil, err
	}
	return result, err
}
