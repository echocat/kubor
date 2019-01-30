package format

import (
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/kubernetes"
	"github.com/olekukonko/tablewriter"
	"io"
	"sort"
	"time"
)

var _ = MustRegister(VariantTable, NewTableFormat())

func NewTableFormat() *TableFormat {
	return &TableFormat{
		Columns: []*TableColumn{{
			Label:        "Kind",
			CellProvider: AggregationToGroupVersionKind,
			Weight:       1000,
		}, {
			Label:        "Namespace",
			CellProvider: AggregationToNamespace,
			Weight:       2000,
		}, {
			Label:        "Name",
			CellProvider: AggregationToName,
			Weight:       3000,
		}, {
			Label:        "Desired",
			CellProvider: AggregationToDesired,
			Weight:       4000,
		}, {
			Label:        "Ready",
			CellProvider: AggregationToReady,
			Weight:       5000,
		}, {
			Label:        "UpToDate",
			CellProvider: AggregationToUpToDate,
			Weight:       6000,
		}, {
			Label:        "Available",
			CellProvider: AggregationToAvailable,
			Weight:       7000,
		}, {
			Label:        "Status",
			CellProvider: AggregationToStatus,
			Weight:       8000,
		}, {
			Label:        "Age",
			CellProvider: AggregationToAge,
			Weight:       9000,
		}},
	}
}

type CellProvider func(kubernetes.Aggregation) (Cell, error)

type TableColumn struct {
	Label        string
	Weight       int
	CellProvider CellProvider
}

func (instance *TableColumn) GetLabel() string {
	return instance.Label
}

func (instance *TableColumn) GetWeight() int {
	return instance.Weight
}

type TableFormat struct {
	Columns []*TableColumn
}

func (instance *TableFormat) Format(to io.Writer) (Task, error) {
	return &tableFormatTask{
		TableFormat: instance,
		to:          to,
		rows:        Rows{},
	}, nil
}

type tableFormatTask struct {
	*TableFormat
	to   io.Writer
	rows Rows
}

func (instance *tableFormatTask) Next(object kubernetes.Object) error {
	row := make(Row)
	aggregation := kubernetes.NewAggregationFor(object)
	for _, column := range instance.Columns {
		if formatted, err := column.CellProvider(aggregation); err != nil {
			return err
		} else {
			row[column] = formatted
		}
	}
	instance.rows = append(instance.rows, &row)
	return nil
}

func (instance *tableFormatTask) Close() error {
	instance.rows.Sort(instance.Columns[1], instance.Columns[2])
	columns := instance.rows.UsedColumns()
	sort.Sort(columns)

	t := tablewriter.NewWriter(instance.to)
	t.SetHeader(columns.Labels())
	t.AppendBulk(instance.rows.Strings(columns...))
	t.Render()
	_, err := instance.to.Write([]byte("\n"))
	return err
}

func AggregationToGroupVersionKind(a kubernetes.Aggregation) (Cell, error) {
	return StringToCell(kubernetes.FormatGroupVersionKind(a.GroupVersionKind()))
}

func AggregationToName(a kubernetes.Aggregation) (Cell, error) {
	return StringToCell(a.GetName())
}

func AggregationToNamespace(a kubernetes.Aggregation) (Cell, error) {
	return StringToCell(a.GetNamespace())
}

func AggregationToDesired(in kubernetes.Aggregation) (Cell, error) {
	return Pint32ToCell(in.GetDesired)
}

func AggregationToReady(in kubernetes.Aggregation) (Cell, error) {
	return Pint32ToCell(in.GetReady)
}

func AggregationToUpToDate(in kubernetes.Aggregation) (Cell, error) {
	return Pint32ToCell(in.GetUpToDate)
}

func AggregationToAvailable(in kubernetes.Aggregation) (Cell, error) {
	return Pint32ToCell(in.GetAvailable)
}

func AggregationToIsReady(in kubernetes.Aggregation) (Cell, error) {
	return PboolToCell(in.IsReady)
}

func AggregationToStatus(in kubernetes.Aggregation) (Cell, error) {
	return PstringToCell(in.GetStatus)
}

func AggregationToAge(in kubernetes.Aggregation) (Cell, error) {
	return PdurationToCell(in.GetAge)
}

func StringToCell(in string) (Cell, error) {
	return StringCell{common.Pstring(in)}, nil
}

func Pint32ToCell(getter func() *int32) (Cell, error) {
	return Int32Cell{getter()}, nil
}

func PboolToCell(getter func() *bool) (Cell, error) {
	return BoolCell{getter()}, nil
}

func PstringToCell(getter func() *string) (Cell, error) {
	return StringCell{getter()}, nil
}

func PdurationToCell(getter func() *time.Duration) (Cell, error) {
	return DurationCell{getter()}, nil
}
