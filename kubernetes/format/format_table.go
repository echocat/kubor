package format

import (
	"github.com/levertonai/kubor/common"
	"github.com/levertonai/kubor/kubernetes"
	"github.com/olekukonko/tablewriter"
	"io"
	"strconv"
)

var _ = MustRegister(VariantTable, NewTableFormat())

func NewTableFormat() *TableFormat {
	return &TableFormat{
		Columns: []*TableColumn{{
			Label:         "Kind",
			CellFormatter: AggregationToGroupVersionKind,
		}, {
			Label:         "Name",
			CellFormatter: AggregationToName,
		}, {
			Label:         "Namespace",
			CellFormatter: AggregationToNamespace,
		}, {
			Label:         "Desired",
			CellFormatter: AggregationToDesired,
		}, {
			Label:         "Ready",
			CellFormatter: AggregationToReady,
		}, {
			Label:         "UpToDate",
			CellFormatter: AggregationToUpToDate,
		}, {
			Label:         "Available",
			CellFormatter: AggregationToAvailable,
		}, {
			Label:         "Is ready",
			CellFormatter: AggregationToIsReady,
		}},
	}
}

type TableCellFormatter func(kubernetes.Aggregation) (*string, error)
type TableFooterCellFormatter func(...kubernetes.Aggregation) (*string, error)

type TableFormat struct {
	Columns []*TableColumn
}

type TableColumn struct {
	Label           string
	CellFormatter   TableCellFormatter
	FooterFormatter TableFooterCellFormatter
}

func (instance *TableFormat) Format(to io.Writer, supplier ObjectSupplier) error {
	columns := make(map[*TableColumn]bool)
	var rows []*tableFormatRow
	for {
		if object, err := supplier(); err != nil {
			return err
		} else if object == nil {
			break
		} else if row, err := instance.newRowFor(object); err != nil {
			return err
		} else {
			rows = append(rows, row)
			for column := range row.cells {
				columns[column] = true
			}
		}
	}
	t := tablewriter.NewWriter(to)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeader(instance.toHeader(columns))
	for _, row := range rows {
		if row, err := instance.formatRow(row); err != nil {
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

func (instance *TableFormat) newRowFor(object kubernetes.Object) (*tableFormatRow, error) {
	row := &tableFormatRow{
		TableFormat: instance,
		cells:       map[*TableColumn]*string{},
	}
	aggregation := kubernetes.NewAggregationFor(object)
	for _, column := range instance.Columns {
		if formatted, err := column.CellFormatter(aggregation); err != nil {
			return nil, err
		} else {
			row.cells[column] = formatted
		}
	}
	return row, nil
}

type tableFormatRow struct {
	*TableFormat
	cells map[*TableColumn]*string
}

func (instance TableFormat) toHeader(columns map[*TableColumn]bool) []string {
	result := make([]string, len(columns))
	var i int
	for _, column := range instance.Columns {
		if columns[column] {
			result[i] = column.Label
			i++
		}
	}
	return result
}

func (instance TableFormat) toFooter(objects ...kubernetes.Object) ([]string, error) {
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

func (instance TableFormat) formatRow(object kubernetes.Object) ([]string, error) {
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

func AggregationToGroupVersionKind(a kubernetes.Aggregation) (*string, error) {
	return common.Pstring(kubernetes.FormatGroupVersionKind(a.GroupVersionKind())), nil
}

func AggregationToName(a kubernetes.Aggregation) (*string, error) {
	return common.Pstring(a.GetName()), nil
}

func AggregationToNamespace(a kubernetes.Aggregation) (*string, error) {
	return common.Pstring(a.GetNamespace()), nil
}

func AggregationToDesired(in kubernetes.Aggregation) (*string, error) {
	return Pint32ToString(in.GetDesired)
}

func AggregationToReady(in kubernetes.Aggregation) (*string, error) {
	return Pint32ToString(in.GetReady)
}

func AggregationToUpToDate(in kubernetes.Aggregation) (*string, error) {
	return Pint32ToString(in.GetUpToDate)
}

func AggregationToAvailable(in kubernetes.Aggregation) (*string, error) {
	return Pint32ToString(in.GetAvailable)
}

func AggregationToIsReady(in kubernetes.Aggregation) (*string, error) {
	return PboolToString(in.IsReady)
}

func Pint32ToString(getter func() *int32) (*string, error) {
	if val := getter(); val != nil {
		return common.Pstring(strconv.FormatInt(int64(*val), 10)), nil
	} else {
		return nil, nil
	}
}

func PboolToString(getter func() *bool) (*string, error) {
	if val := getter(); val == nil {
		return nil, nil
	} else if *val {
		return common.Pstring("Yes"), nil
	} else {
		return common.Pstring("No"), nil
	}
}
