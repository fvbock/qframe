// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package bcolumn

import (
	"fmt"
	"reflect"

	"github.com/tobgu/qframe/internal/column"
	"github.com/tobgu/qframe/internal/index"
)

// Code generated from template/column.go DO NOT EDIT

type Column struct {
	data []bool
}

func New(d []bool) Column {
	return Column{data: d}
}

func NewConst(val bool, count int) Column {
	var nullVal bool
	data := make([]bool, count)
	if val != nullVal {
		for i := range data {
			data[i] = val
		}
	}

	return Column{data: data}
}

// Apply single argument function. The result may be a column
// of a different type than the current column.
func (c Column) Apply1(fn interface{}, ix index.Int) (interface{}, error) {
	switch t := fn.(type) {
	case func(bool) int:
		result := make([]int, len(c.data))
		for _, i := range ix {
			result[i] = t(c.data[i])
		}
		return result, nil
	case func(bool) float64:
		result := make([]float64, len(c.data))
		for _, i := range ix {
			result[i] = t(c.data[i])
		}
		return result, nil
	case func(bool) bool:
		result := make([]bool, len(c.data))
		for _, i := range ix {
			result[i] = t(c.data[i])
		}
		return result, nil
	case func(bool) *string:
		result := make([]*string, len(c.data))
		for _, i := range ix {
			result[i] = t(c.data[i])
		}
		return result, nil
	default:
		return nil, fmt.Errorf("cannot apply type %#v to column", fn)
	}
}

// Apply double argument function to two columns. Both columns must have the
// same type. The resulting column will have the same type as this column.
func (c Column) Apply2(fn interface{}, s2 column.Column, ix index.Int) (column.Column, error) {
	ss2, ok := s2.(Column)
	var typ bool
	if !ok {
		return Column{}, fmt.Errorf("%v.apply2: invalid column type: %#v", reflect.TypeOf(typ), s2)
	}

	t, ok := fn.(func(bool, bool) bool)
	if !ok {
		return Column{}, fmt.Errorf("%v.apply2: invalid function type: %#v", reflect.TypeOf(typ), fn)
	}

	result := make([]bool, len(c.data))
	for _, i := range ix {
		result[i] = t(c.data[i], ss2.data[i])
	}

	return New(result), nil
}

func (c Column) subset(index index.Int) Column {
	data := make([]bool, len(index))
	for i, ix := range index {
		data[i] = c.data[ix]
	}

	return Column{data: data}
}

func (c Column) Subset(index index.Int) column.Column {
	return c.subset(index)
}

func (c Column) Comparable(reverse bool) column.Comparable {
	if reverse {
		return Comparable{data: c.data, ltValue: column.GreaterThan, gtValue: column.LessThan}
	}

	return Comparable{data: c.data, ltValue: column.LessThan, gtValue: column.GreaterThan}
}

func (c Column) String() string {
	return fmt.Sprintf("%v", c.data)
}

func (c Column) Len() int {
	return len(c.data)
}

func (c Column) Aggregate(indices []index.Int, fn interface{}) (column.Column, error) {
	var actualFn func([]bool) bool
	var ok bool

	switch t := fn.(type) {
	case string:
		actualFn, ok = aggregations[t]
		if !ok {
			return nil, fmt.Errorf("aggregation function %c is not defined for column", fn)
		}
	case func([]bool) bool:
		actualFn = t
	default:
		// TODO: Genny is buggy and won't let you use your own errors package.
		//       We use a standard error here for now.
		return nil, fmt.Errorf("invalid aggregation function type: %v", t)
	}

	data := make([]bool, 0, len(indices))
	for _, ix := range indices {
		subS := c.subset(ix)
		data = append(data, actualFn(subS.data))
	}

	return Column{data: data}, nil
}

func (c Column) View(ix index.Int) View {
	return View{data: c.data, index: ix}
}

type Comparable struct {
	data    []bool
	ltValue column.CompareResult
	gtValue column.CompareResult
}

type View struct {
	data  []bool
	index index.Int
}

func (v View) ItemAt(i int) bool {
	return v.data[v.index[i]]
}

func (v View) Len() int {
	return len(v.index)
}

// TODO: This forces an alloc, as an alternative a slice could be taken
//       as input that can be (re)used by the client. Are there use cases
//       where this would actually make sense?
func (v View) Slice() []bool {
	result := make([]bool, v.Len())
	for i, j := range v.index {
		result[i] = v.data[j]
	}
	return result
}
