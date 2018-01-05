package postgis

import (
	"fmt"
	"strings"

	entities "github.com/gost/core"
	"github.com/gost/godata"
)

var filterToStringMap map[int]func(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string

// init is used to work around the initialization loop error (circular reference)
func init() {
	filterToStringMap = map[int]func(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string{
		godata.FilterTokenNav:       filterNavToString,
		godata.FilterTokenLogical:   filterLogicalToString,
		godata.FilterTokenFunc:      filterFuncToString,
		godata.FilterTokenOp:        filterOpToString,
		godata.FilterTokenGeography: filterGeographyToString,
		godata.FilterTokenLiteral:   filterLiteralToString,
		godata.FilterTokenLambda:    filterDefaultToString,
		godata.FilterTokenNull:      filterDefaultToString,
		godata.FilterTokenIt:        filterDefaultToString,
		godata.FilterTokenRoot:      filterDefaultToString,
		godata.FilterTokenFloat:     filterDefaultToString,
		godata.FilterTokenInteger:   filterDefaultToString,
		godata.FilterTokenString:    filterDefaultToString,
		godata.FilterTokenDate:      filterDefaultToString,
		godata.FilterTokenTime:      filterDefaultToString,
		godata.FilterTokenDateTime:  filterDefaultToString,
		godata.FilterTokenBoolean:   filterDefaultToString,
	}
}

func filterNavToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string {
	q := ""
	for i, part := range pn.Children {
		if i == 0 {
			q += fmt.Sprintf("%v ", strings.ToLower(qb.createFilter(et, part, false)))
			continue
		}

		arrow := "->"
		if i+1 == len(pn.Children) {
			arrow = "->>"
		}
		q += fmt.Sprintf("%v '%v'", arrow, part.Token.Value)
	}
	return q
}

func filterLogicalToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string {
	left := qb.createFilter(et, pn.Children[0], false)

	if len(pn.Children) == 1 && strings.ToLower(pn.Token.Value) == "not" {
		return fmt.Sprintf("%v %v", qb.odataLogicalOperatorToPostgreSQL(pn.Token.Value), left)
	}

	right := qb.createFilter(et, pn.Children[1], false)
	left, right = qb.prepareFilter(et, pn.Children[0].Token.Value, left, pn.Children[1].Token.Value, right)

	// Workaround for faulty OGC test
	result := "observation.data -> 'result'"
	if len(qb.odataLogicalOperatorToPostgreSQL(pn.Token.Value)) > 0 {
		if left == result {
			if strings.Index(right, "'") != 0 {
				left = qb.CastObservationResult(left, "double precision")
			} else {
				left = "observation.data ->> 'result'"
			}
		} else if right == result {
			if strings.Index(left, "'") != 0 {
				right = qb.CastObservationResult(right, "double precision")
			} else {
				right = "observation.data ->> 'result'"
			}
		}
	}
	// End workaround

	return fmt.Sprintf("%v %v %v", left, qb.odataLogicalOperatorToPostgreSQL(pn.Token.Value), right)
}

func filterFuncToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string {
	if convertFunction, ok := funcToStringMap[pn.Token.Value]; ok {
		return convertFunction(qb, pn, et)
	}

	return ""
}

func filterDefaultToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string {
	return fmt.Sprintf("%v", pn.Token.Value)
}

func filterOpToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string {
	if pn.Token.Value == "add" {
		return qb.createArithmetic(et, pn, "+", "double precision")
	} else if pn.Token.Value == "sub" {
		return qb.createArithmetic(et, pn, "-", "double precision")
	} else if pn.Token.Value == "mul" {
		return qb.createArithmetic(et, pn, "*", "double precision")
	} else if pn.Token.Value == "div" {
		return qb.createArithmetic(et, pn, "/", "double precision")
	} else if pn.Token.Value == "mod" {
		return qb.createArithmetic(et, pn, "%", "integer")
	}

	return ""
}

func filterGeographyToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string {
	return fmt.Sprintf("ST_GeomFromText(%v)", pn.Children[0].Token.Value)
}

func filterLiteralToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string {
	p := selectMappings[et][strings.ToLower(pn.Token.Value)]
	if p != "" {
		return p
	}

	if ignoreSelectAs && selectMappingsIgnore[et][pn.Token.Value] {
		return fmt.Sprintf("%s.%v", et.ToString(), pn.Token.Value)
	}

	return pn.Token.Value
}
