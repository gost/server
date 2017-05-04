package odata

import (
	"github.com/gost/godata"
	"net/url"
)

var SupportedExpandParameters map[string][]string
var SupportedSelectParameters map[string][]string

// QueryOptions extents upon godata.GoDataQuery to implement extra
// odata functions not found in the godata package
type QueryOptions struct {
	godata.GoDataQuery
	Value      *GoDataValueQuery
	Ref        *GoDataRefQuery
	RawExpand  string
	RawFilter  string
	RawOrderBy string
}

func (q *QueryOptions) ExpandParametersSupported() bool {
	return true
}

func (q *QueryOptions) SelectParametersSupported() bool {
	return true
}

type GoDataValueQuery bool
type GoDataRefQuery bool

func ExpandItemToQueryOptions(ei *godata.ExpandItem) *QueryOptions {
	qo := QueryOptions{}
	qo.Top = ei.Top
	qo.Filter = ei.Filter
	qo.OrderBy = ei.OrderBy
	qo.Search = ei.Search
	qo.Select = ei.Select
	qo.Skip = ei.Skip

	return &qo
}

func ParseUrlQuery(query url.Values) (*QueryOptions, error) {
	if query == nil || len(query) == 0 {
		return nil, nil
	}

	qo, err := godata.ParseUrlQuery(query)
	if err != nil {
		return nil, err
	}

	result := &QueryOptions{}
	result.GoDataQuery = *qo

	value := query.Get("$value")

	val := GoDataValueQuery(false)
	if value != "" {
		val = GoDataValueQuery(true)
	}
	result.Value = &val

	ref := GoDataRefQuery(false)
	if value != "" {
		ref = GoDataRefQuery(true)
	}
	result.Ref = &ref

	//store raw queries
	result.RawExpand = query.Get("$expand")
	result.RawFilter = query.Get("$filter")
	result.RawOrderBy = query.Get("$orderby")

	return result, err
}
