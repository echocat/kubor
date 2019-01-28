package format

import (
	"errors"
	"fmt"
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/kubernetes"
	"github.com/olekukonko/tablewriter"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strconv"
)

var (
	ErrUnsupportedObject = errors.New("unsupported object")
)

func MustRegisterTableBasedObjectFormatOf(title string, kinds ...schema.GroupVersionKind) *TableFormatter {
	result := MustNewTemplateBasedObjectFormatOf(title, kinds...)
	Provider.MustRegister(VariantTable, result)
	return result
}

func MustNewTemplateBasedObjectFormatOf(title string, kinds ...schema.GroupVersionKind) *TableFormatter {
	if result, err := NewTemplateBasedObjectFormatOf(title, kinds...); err != nil {
		panic(err)
	} else {
		return result
	}
}

func NewTemplateBasedObjectFormatOf(title string, kinds ...schema.GroupVersionKind) (*TableFormatter, error) {
	return &TableFormatter{
		Title:         title,
		SupportedGvks: NormalizeGroupVersionKinds(kinds),
		Columns:       []TableColumn{},
	}, nil
}

type TableCellFormatter func(runtime.Object) (string, error)
type TableFooterCellFormatter func(...runtime.Object) (string, error)

type TableFormatter struct {
	Title         string
	Columns       []TableColumn
	SupportedGvks []schema.GroupVersionKind
}

func (instance *TableFormatter) WithColumn(label string, formatter TableCellFormatter) *TableFormatter {
	instance.Columns = append(instance.Columns, TableColumn{
		Label:         label,
		CellFormatter: formatter,
	})
	return instance
}

type TableColumn struct {
	Label           string
	CellFormatter   TableCellFormatter
	FooterFormatter TableFooterCellFormatter
}

func (instance TableFormatter) Supports(gvks ...schema.GroupVersionKind) bool {
	if len(instance.SupportedGvks) == 0 {
		return true
	}
	for _, gvk := range gvks {
		if !instance.supports(gvk) {
			return false
		}
	}
	return true
}

func (instance TableFormatter) supports(gvk schema.GroupVersionKind) bool {
	for _, candidate := range instance.SupportedGvks {
		if NormalizeGroupVersionKind(candidate) == gvk {
			return true
		}
	}
	return false
}

func (instance TableFormatter) Format(to io.Writer, objects ...runtime.Object) error {
	if len(objects) == 0 {
		return nil
	}
	title := instance.Title
	if title == "" {
		title = FormatGroupVersionKind(objects[0].GetObjectKind().GroupVersionKind())
	}
	if _, err := fmt.Fprintf(to, "%s\n", title); err != nil {
		return err
	}
	t := tablewriter.NewWriter(to)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader(instance.toHeader())
	for _, object := range objects {
		if row, err := instance.toRow(object); err != nil {
			return err
		} else {
			t.Append(row)
		}
	}
	if footer, err := instance.toFooter(objects...); err != nil {
		return err
	} else if footer != nil {
		t.SetFooter(footer)
	}
	t.Render()
	_, err := to.Write([]byte("\n"))
	return err
}

func (instance TableFormatter) toHeader() []string {
	result := make([]string, len(instance.Columns))
	for i, column := range instance.Columns {
		result[i] = column.Label
	}
	return result
}

func (instance TableFormatter) toFooter(objects ...runtime.Object) ([]string, error) {
	atLeastOneFooterFound := false
	result := make([]string, len(instance.Columns))
	for i, column := range instance.Columns {
		if column.FooterFormatter == nil {
			result[i] = ""
		} else if formatted, err := column.FooterFormatter(objects...); err != nil {
			return nil, err
		} else {
			result[i] = formatted
			atLeastOneFooterFound = true
		}
	}
	if !atLeastOneFooterFound {
		return nil, nil
	}
	return result, nil
}

func (instance TableFormatter) toRow(object runtime.Object) ([]string, error) {
	result := make([]string, len(instance.Columns))
	var i int
	for _, column := range instance.Columns {
		if formatted, err := column.CellFormatter(object); err != nil {
			return []string{}, err
		} else {
			result[i] = formatted
		}
		i++
	}
	return result, nil
}

func (instance TableFormatter) ShouldBeCombined() bool {
	return true
}

func ObjectPathFormatter(path ...string) TableCellFormatter {
	return func(object runtime.Object) (string, error) {
		var val interface{}
		if u, ok := object.(runtime.Unstructured); ok {
			val = common.GetObjectPathValue(u.UnstructuredContent(), path...)
		} else {
			val = common.GetObjectPathValue(object, path...)
		}
		return fmt.Sprint(val), nil
	}
}

func AggregationFormatter(mapper func(kubernetes.Aggregation) (string, error)) TableCellFormatter {
	return func(object runtime.Object) (string, error) {
		if u, ok := object.(*unstructured.Unstructured); ok {
			aggregation := kubernetes.NewAggregationFor(u)
			return mapper(aggregation)
		}
		return "", ErrUnsupportedObject
	}
}

func AggregationToDesired(in kubernetes.Aggregation) (string, error) {
	return Pint32ToString(in.Desired)
}

func AggregationToReady(in kubernetes.Aggregation) (string, error) {
	return Pint32ToString(in.Ready)
}

func AggregationToUpToDate(in kubernetes.Aggregation) (string, error) {
	return Pint32ToString(in.UpToDate)
}

func AggregationToAvailable(in kubernetes.Aggregation) (string, error) {
	return Pint32ToString(in.Available)
}

func AggregationToIsReady(in kubernetes.Aggregation) (string, error) {
	return PboolToString(in.IsReady)
}

func Pint32ToString(getter func() *int32) (string, error) {
	if val := getter(); val != nil {
		return strconv.FormatInt(int64(*val), 10), nil
	} else {
		return "n/a", nil
	}
}

func PboolToString(getter func() *bool) (string, error) {
	if val := getter(); val == nil {
		return "n/a", nil
	} else if *val {
		return "Yes", nil
	} else {
		return "No", nil
	}
}
