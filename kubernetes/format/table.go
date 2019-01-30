package format

import (
	"fmt"
	"github.com/levertonai/kubor/common"
	"sort"
	"strconv"
	"time"
)

type Cell interface {
	Formatted() *string
	IsBlank() bool
	Interface() interface{}
	LessThan(Cell) bool
}

type Column interface {
	GetLabel() string
	GetWeight() int
}

type Columns []Column

func (instance Columns) Labels() []string {
	result := make([]string, len(instance))
	for i, column := range instance {
		result[i] = column.GetLabel()
	}
	return result
}

func (instance Columns) Len() int      { return len(instance) }
func (instance Columns) Swap(i, j int) { instance[i], instance[j] = instance[j], instance[i] }
func (instance Columns) Less(i, j int) bool {
	return instance[i].GetWeight() < instance[j].GetWeight()
}

type Row map[Column]Cell

func (instance Row) UsedColumns() Columns {
	var result Columns
	for column, value := range instance {
		if !value.IsBlank() {
			result = append(result, column)
		}
	}
	return result
}

func (instance Row) LessThan(what Row, on ...Column) bool {
	for _, column := range on {
		if instance.Cell(column).LessThan(what.Cell(column)) {
			return true
		}
	}
	return false
}

func (instance Row) Cell(column Column) Cell {
	if instance == nil {
		return EmptyCell{}
	}
	if cell, ok := instance[column]; ok {
		return cell
	}
	return EmptyCell{}
}

func (instance Row) Strings(of ...Column) []string {
	result := make([]string, len(of))
	for i, column := range of {
		cell := instance.Cell(column)
		formatted := cell.Formatted()
		if formatted != nil {
			result[i] = *formatted
		} else {
			result[i] = ""
		}
	}
	return result
}

type Rows []*Row

func (instance Rows) UsedColumns() Columns {
	columns := make(map[Column]bool)
	for _, row := range instance {
		for _, column := range row.UsedColumns() {
			columns[column] = true
		}
	}
	result := make(Columns, len(columns))
	var i int
	for column := range columns {
		result[i] = column
		i++
	}
	return result
}

func (instance Rows) Strings(of ...Column) [][]string {
	result := make([][]string, len(instance))
	for i, row := range instance {
		result[i] = row.Strings(of...)
	}
	return result
}

func (instance *Rows) Sort(by ...Column) {
	sort.Sort(byColumnsRows{
		Rows:      instance,
		byColumns: by,
	})
}

type byColumnsRows struct {
	*Rows
	byColumns Columns
}

func (instance byColumnsRows) Len() int { return len(*instance.Rows) }
func (instance byColumnsRows) Swap(i, j int) {
	(*instance.Rows)[i], (*instance.Rows)[j] = (*instance.Rows)[j], (*instance.Rows)[i]
}
func (instance byColumnsRows) Less(i, j int) bool {
	return (*instance.Rows)[i].LessThan(*(*instance.Rows)[j], instance.byColumns...)
}

type EmptyCell struct{}

func (instance EmptyCell) Formatted() *string {
	return nil
}

func (instance EmptyCell) IsBlank() bool {
	return true
}

func (instance EmptyCell) Interface() interface{} {
	return nil
}

func (instance EmptyCell) LessThan(what Cell) bool {
	if _, ok := what.(*EmptyCell); ok {
		return false
	}
	if _, ok := what.(EmptyCell); ok {
		return false
	}
	return what.Interface() != nil
}

type StringCell struct {
	Content *string
}

func (instance StringCell) Formatted() *string {
	return instance.Content
}

func (instance StringCell) IsBlank() bool {
	return instance.Content == nil
}

func (instance StringCell) Interface() interface{} {
	return instance.Content
}

func (instance StringCell) LessThan(what Cell) bool {
	if sv, ok := what.(*StringCell); ok {
		what = *sv
	}
	if sv, ok := what.(StringCell); !ok {
		return false
	} else if instance.Content == nil && sv.Content == nil {
		return false
	} else if instance.Content == nil {
		return true
	} else if sv.Content == nil {
		return false
	} else {
		return *instance.Content < *sv.Content
	}
}

type Int32Cell struct {
	Content *int32
}

func (instance Int32Cell) Formatted() *string {
	if instance.Content == nil {
		return nil
	}
	return common.Pstring(strconv.FormatInt(int64(*instance.Content), 10))
}

func (instance Int32Cell) IsBlank() bool {
	return instance.Content == nil
}

func (instance Int32Cell) Interface() interface{} {
	return instance.Content
}

func (instance Int32Cell) LessThan(what Cell) bool {
	if sv, ok := what.(*Int32Cell); ok {
		what = *sv
	}
	if sv, ok := what.(Int32Cell); !ok {
		return false
	} else if instance.Content == nil && sv.Content == nil {
		return false
	} else if instance.Content == nil {
		return true
	} else if sv.Content == nil {
		return false
	} else {
		return *instance.Content < *sv.Content
	}
}

type BoolCell struct {
	Content *bool
}

func (instance BoolCell) Formatted() *string {
	if instance.Content == nil {
		return nil
	}
	switch *instance.Content {
	case true:
		return common.Pstring("Yes")
	default:
		return common.Pstring("No")
	}
}

func (instance BoolCell) IsBlank() bool {
	return instance.Content == nil
}

func (instance BoolCell) Interface() interface{} {
	return instance.Content
}

func (instance BoolCell) LessThan(what Cell) bool {
	if sv, ok := what.(*BoolCell); ok {
		what = *sv
	}
	if sv, ok := what.(BoolCell); !ok {
		return false
	} else if instance.Content == nil && sv.Content == nil {
		return false
	} else if instance.Content == nil {
		return true
	} else if sv.Content == nil {
		return false
	} else if *instance.Content == *sv.Content {
		return false
	} else if *sv.Content {
		return true
	} else {
		return false
	}
}

type DurationCell struct {
	Content *time.Duration
}

func (instance DurationCell) Formatted() *string {
	if instance.Content == nil {
		return nil
	}
	var result string
	d := *instance.Content
	if d < time.Minute {
		result = fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		result = fmt.Sprintf("%.0fm", d.Minutes())
	} else if d < (time.Hour * 24) {
		result = fmt.Sprintf("%.0fh", d.Hours())
	} else {
		days := d.Hours() / 24.0
		result = fmt.Sprintf("%.0fd", days)
	}
	return common.Pstring(result)
}

func (instance DurationCell) IsBlank() bool {
	return instance.Content == nil
}

func (instance DurationCell) Interface() interface{} {
	return instance.Content
}

func (instance DurationCell) LessThan(what Cell) bool {
	if sv, ok := what.(*DurationCell); ok {
		what = *sv
	}
	if sv, ok := what.(DurationCell); !ok {
		return false
	} else if instance.Content == nil && sv.Content == nil {
		return false
	} else if instance.Content == nil {
		return true
	} else if sv.Content == nil {
		return false
	} else {
		return *instance.Content < *sv.Content
	}
}

type TimeCell struct {
	Content *time.Time
}

func (instance TimeCell) Formatted() *string {
	if instance.Content == nil {
		return nil
	}
	return common.Pstring(instance.Content.Format(time.RFC3339))
}

func (instance TimeCell) IsBlank() bool {
	return instance.Content == nil
}

func (instance TimeCell) Interface() interface{} {
	return instance.Content
}

func (instance TimeCell) LessThan(what Cell) bool {
	if sv, ok := what.(*TimeCell); ok {
		what = *sv
	}
	if sv, ok := what.(TimeCell); !ok {
		return false
	} else if instance.Content == nil && sv.Content == nil {
		return false
	} else if instance.Content == nil {
		return true
	} else if sv.Content == nil {
		return false
	} else {
		return instance.Content.Before(*sv.Content)
	}
}
