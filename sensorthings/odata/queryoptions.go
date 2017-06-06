package odata

import (
	"github.com/gost/godata"
	"net/url"
)

// SupportedExpandParameters contains a list of endpoints with their supported expand parameters
var SupportedExpandParameters map[string][]string

// SupportedSelectParameters contains a list of endpoints with their supported select parameters
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

// ExpandParametersSupported returns if the QueryOptions expand request is supported by the endpoints
// todo: implement check
func (q *QueryOptions) ExpandParametersSupported() bool {
	return true
}

// SelectParametersSupported returns if the QueryOptions select request is supported by the endpoints
// todo: implement check
func (q *QueryOptions) SelectParametersSupported() bool {
	return true
}

// GoDataValueQuery true when $value is requested false if not
type GoDataValueQuery bool

// GoDataRefQuery true when $ref is requested false if not
type GoDataRefQuery bool

// ExpandItemToQueryOptions converts an ExpandItem into QueryOptions
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

// ParseURLQuery parses an incoming url query into QueryOptions
func ParseURLQuery(query url.Values) (*QueryOptions, error) {
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

	value = query.Get("$ref")
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
