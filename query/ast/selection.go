package ast

import (
	"strconv"
	"strings"
)

// Selection represents a statement that aggregates data in a query.
type Selection struct {
	Name       string
	Dimensions []*VarRef
	Fields     Fields
}

func (s *Selection) node()      {}
func (s *Selection) statement() {}

// NewSelection creates a new Selection instance.
func NewSelection() *Selection {
	return &Selection{}
}

// HasAggregateFields returns true if there are any fields that use aggregation.
func (s *Selection) HasAggregateFields() bool {
	for _, field := range s.Fields {
		if field.IsAggregate() {
			return true
		}
	}
	return false
}

// HasNonAggregateFields returns true if there are any fields that do not use aggregation.
func (s *Selection) HasNonAggregateFields() bool {
	for _, field := range s.Fields {
		if !field.IsAggregate() {
			return true
		}
	}
	return false
}

// Path returns a period-delimited join of the name and dimension names.
func (s *Selection) Path() string {
	var names []string
	if s.Name != "" {
		names = append(names, s.Name)
	}
	for _, d := range s.Dimensions {
		names = append(names, d.Name)
	}
	return strings.Join(names, ".")
}

func (s *Selection) String() string {
	str := "SELECT "

	arr := []string{}
	for _, field := range s.Fields {
		arr = append(arr, field.String())
	}
	str += strings.Join(arr, ", ")

	if len(s.Dimensions) > 0 {
		var names []string
		for _, dimension := range s.Dimensions {
			names = append(names, dimension.Name)
		}
		str += " GROUP BY @" + strings.Join(names, ", @")
	}
	if s.Name != "" {
		str += " INTO " + strconv.Quote(s.Name)
	}
	return str
}
