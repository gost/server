package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
	"strings"
)

// QueryBuilder can construct queries based on entities and QueryOptions
type QueryBuilder struct {
	maxTop int
	schema string
	tables map[entities.EntityType]string
}

// CreateQueryBuilder instantiates a new queryBuilder, the queryBuilder is used to create
// select queries based on the given entities, id en QueryOptions (ODATA)
// schema is the used database schema can be empty, maxTop is the maximum top the query should return
func CreateQueryBuilder(schema string, maxTop int) *QueryBuilder {
	qb := &QueryBuilder{
		schema: schema,
		maxTop: maxTop,
		tables: createTableMappings(schema),
	}

	return qb
}

// removeSchema removes the prefix in front of a table
func (qb *QueryBuilder) removeSchema(table string) string {
	i := strings.Index(table, ".")
	if i == -1 {
		return table
	}
	return table[i+1:]
}

// getLimit returns the max entities to retrieve, this number is set by ODATA's
// $top, if not provided use the global value
func (qb *QueryBuilder) getLimit(qo *odata.QueryOptions) string {
	if qo != nil && !qo.QueryTop.IsNil() {
		return fmt.Sprintf("%v", qo.QueryTop.Limit)
	}
	return fmt.Sprintf("%v", qb.maxTop)
}

// getOffset returns the offset, this number is set by ODATA's
// $skip, if not provided do not skip anything = return "0"
func (qb *QueryBuilder) getOffset(qo *odata.QueryOptions) string {
	if qo != nil && !qo.QuerySkip.IsNil() {
		return fmt.Sprintf("%v", qo.QuerySkip.Index)
	}
	return "0"
}

// getOrderBy returns the string that needs to be placed after ORDER BY, this is set using
// ODATA's $orderby if not given use the default ORDER BY "table".id DESC
func (qb *QueryBuilder) getOrderBy(et entities.EntityType, qo *odata.QueryOptions) string {
	if qo != nil && !qo.QueryOrderBy.IsNil() {
		return fmt.Sprintf("%v %v", selectMappings[et][strings.ToLower(qo.QueryOrderBy.Property)], strings.ToUpper(qo.QueryOrderBy.Suffix))
	}

	return fmt.Sprintf("%s DESC", selectMappings[et]["id"])
}

// getSelect return the select string that needs to be placed after SELECT in the query
// select is set by ODATA's $select, if not set get all properties for the given entity (return all)
// addID to true if it needs to be added and isn't in QuerySelect.Params, addAs to true if a field needs to be
// outputted with AS [name]
func (qb *QueryBuilder) getSelect(et entities.Entity, qo *odata.QueryOptions, qpi *QueryParseInfo, addID bool, addAs bool, fromAs bool, isExpand bool, selectString string) string {
	var properties []string
	if qo == nil || qo.QuerySelect == nil || len(qo.QuerySelect.Params) == 0 {
		properties = et.GetPropertyNames()
	} else {
		idAdded := false
		for _, p := range qo.QuerySelect.Params {
			if p == "id" {
				idAdded = true
			}
			for _, pn := range et.GetPropertyNames() {
				if strings.ToLower(p) == strings.ToLower(pn) {
					if p == "id" {
						properties = append([]string{"id"}, properties...)
					} else {
						properties = append(properties, pn)
					}
				}
			}
		}
		if addID && !idAdded {
			properties = append([]string{"id"}, properties...)
		}
	}

	// ToDo: this is a fix for supporting $expand=Observations/FeatureOfInterest, try to add observationFeatureOfInterestID in a different way
	if isExpand && et.GetEntityType() == entities.EntityTypeObservation {
		properties = append([]string{observationFeatureOfInterestID}, properties...)
	}

	for _, p := range properties {
		toAdd := ""
		if len(selectString) > 0 {
			toAdd += ", "
		}
		entityType := et.GetEntityType()

		field := ""
		if fromAs {
			field = qb.addAsPrefix(qpi, fmt.Sprintf("%s.%s", tableMappings[entityType], asMappings[entityType][strings.ToLower(p)]))
		} else {
			field = selectMappings[entityType][strings.ToLower(p)]
		}

		if addAs {
			if !isExpand {
				selectString += fmt.Sprintf("%s%s AS %s", toAdd, field, qb.addAsPrefix(qpi, asMappings[entityType][strings.ToLower(p)]))
			} else {
				selectString += fmt.Sprintf("%s%s AS %s", toAdd, field, asMappings[entityType][strings.ToLower(p)])
			}
		} else {
			selectString += fmt.Sprintf("%s%s", toAdd, field)
		}
	}

	if qpi != nil && len(qpi.SubEntities) > 0 {
		for _, subQPI := range qpi.SubEntities {
			selectString = qb.getSelect(subQPI.Entity, subQPI.ExpandOperation.QueryOptions, subQPI, true, true, true, false, selectString)
		}
	}

	return selectString
}

// addAsPrefix adds a prefix in front of the current as for example A_, B_ to be able to
// distinguish the different results if multiple tables are requested of the same type
func (qb *QueryBuilder) addAsPrefix(qpi *QueryParseInfo, as string) string {
	if qpi == nil {
		return as
	} else {
		return fmt.Sprintf("%v_%s", qpi.AsPrefix, as)
	}
}

func (qb *QueryBuilder) createJoin(e1 entities.Entity, e2 entities.Entity, isExpand bool, qo *odata.QueryOptions, qpi *QueryParseInfo, joinString string) string {
	if e2 != nil {
		nqo := qo
		et2 := e2.GetEntityType()

		asPrefix := ""
		if qpi != nil && qpi.Parent != nil && qpi.Parent.QueryIndex != 0 {
			asPrefix = qpi.Parent.AsPrefix
		}

		if !isExpand {
			nqo = &odata.QueryOptions{
				QuerySelect: &odata.QuerySelect{Params: []string{"id"}},
			}
			joinString = fmt.Sprintf("%s"+
				"INNER JOIN LATERAL ("+
				"SELECT %s FROM %s %s "+
				"%s) "+
				"AS %s on true ", joinString,
				qb.getSelect(e2, nqo, nil, true, true, false, false, ""),
				qb.tables[et2],
				getJoin(qb.tables, et2, e1.GetEntityType(), asPrefix),
				qb.getFilterQueryString(et2, nqo, true),
				tableMappings[et2])
		} else {
			joinString = fmt.Sprintf("%s"+
				"LEFT JOIN LATERAL ("+
				"SELECT %s FROM %s %s "+
				"%s"+
				"ORDER BY %s "+
				"LIMIT %s OFFSET %s) AS %s on true ", joinString,
				qb.getSelect(e2, nqo, qpi, true, true, false, true, ""),
				qb.tables[et2],
				getJoin(qb.tables, et2, e1.GetEntityType(), asPrefix),
				qb.getFilterQueryString(et2, nqo, true),
				qb.getOrderBy(et2, nqo),
				qb.getLimit(nqo),
				qb.getOffset(nqo),
				qb.addAsPrefix(qpi, tableMappings[et2]))
		}
	}

	if qpi != nil && len(qpi.SubEntities) > 0 {
		for _, subQPI := range qpi.SubEntities {
			joinString = qb.createJoin(subQPI.Parent.Entity, subQPI.Entity, true, subQPI.ExpandOperation.QueryOptions, subQPI, joinString)
		}
	}

	return joinString
}

func (qb *QueryBuilder) constructQueryParseInfo(operation *odata.ExpandOperation, main *QueryParseInfo, from *QueryParseInfo) {
	nQPI := &QueryParseInfo{}
	nQPI.Init(operation.Entity.GetEntityType(), main.GetNextQueryIndex(), from, operation)
	main.SubEntities = append(main.SubEntities, nQPI)

	if operation.ExpandOperation != nil {
		qb.constructQueryParseInfo(operation.ExpandOperation, main, nQPI)
	}
}

// createFilterQueryString converts an OData query string found in odata.QueryOptions.QueryFilter to a PostgreSQL query string
// ParamFactory is used for converting SensorThings parameter names to postgres field names
// Convert receives a name such as phenomenonTime and returns "data ->> 'id'" true, returns
// false if parameter cannot be converted
func (qb *QueryBuilder) getFilterQueryString(et entities.EntityType, qo *odata.QueryOptions, addWhere bool) string {
	q := ""
	if qo != nil && !qo.QueryFilter.IsNil() {
		if addWhere {
			q += " WHERE "
		}
		ps, ops := qo.QueryFilter.Predicate.Split()
		for i, p := range ps {
			operator, _ := qb.odataOperatorToPostgreSQL(p.Operator)
			q += fmt.Sprintf("%v %v %v", selectMappings[et][strings.ToLower(fmt.Sprintf("%v", p.Left))], operator, p.Right)
			if len(ops)-1 >= i {
				q += fmt.Sprintf(" %v ", ops[i])
			}
		}
		q += " "
	}

	return q
}

// OdataOperatorToPostgreSQL converts an odata.OdataOperator to a PostgreSQL string representation
func (qb *QueryBuilder) odataOperatorToPostgreSQL(o odata.Operator) (string, error) {
	switch o {
	case odata.And:
		return "AND", nil
	case odata.Or:
		return "OR", nil
	case odata.Not:
		return "NOT", nil
	case odata.Equals:
		return "=", nil
	case odata.NotEquals:
		return "!=", nil
	case odata.GreaterThan:
		return ">", nil
	case odata.GreaterThanOrEquals:
		return ">=", nil
	case odata.LessThan:
		return "<", nil
	case odata.LessThanOrEquals:
		return "<=", nil
	case odata.IsNull:
		return "IS NULL", nil
	}

	return "", fmt.Errorf("Operator %v not implemented", o.ToString())
}

// CreateQuery creates a new query based on given input
//   e1: entity to get
//   e2: from entity
//   id: e2 == nil: where e1.id = ... | e2 != nil: where e2.id = ...
// example: Datastreams(1)/Thing = CreateQuery(&entities.Thing, &entities.Datastream, 1, nil)
func (qb *QueryBuilder) CreateQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, qo *odata.QueryOptions) (string, *QueryParseInfo) {
	et1 := e1.GetEntityType()
	et2 := e1.GetEntityType()
	if e2 != nil { // 2nd entity is given, this means get e1 by e2
		et2 = e2.GetEntityType()
	}

	eo := &odata.ExpandOperation{
		QueryOptions: qo,
	}
	qpi := &QueryParseInfo{}
	qpi.Init(et1, 0, nil, eo)

	if qo != nil && !qo.QueryExpand.IsNil() {
		qpi.SubEntities = make([]*QueryParseInfo, 0)
		for _, qe := range qo.QueryExpand.Operations {
			qb.constructQueryParseInfo(qe, qpi, qpi)
		}
	}

	queryString := fmt.Sprintf("SELECT %s FROM %s %s", qb.getSelect(e1, qo, qpi, true, true, false, false, ""), qb.tables[et1], qb.createJoin(e1, e2, false, qo, qpi, ""))
	if id != nil {
		if e2 == nil {
			queryString = fmt.Sprintf("%s WHERE %s = %v", queryString, selectMappings[et2][idField], id)
		} else {
			queryString = fmt.Sprintf("%s WHERE %s.%s = %v", queryString, tableMappings[et2], asMappings[et2][idField], id)
		}
	}

	if qo != nil && !qo.QueryFilter.IsNil() {
		if id != nil {
			queryString = fmt.Sprintf("%s AND %s", queryString, qb.getFilterQueryString(et1, qo, false))
		} else {
			queryString = fmt.Sprintf("%s %s", queryString, qb.getFilterQueryString(et1, qo, true))
		}
	}

	queryString = fmt.Sprintf("%s ORDER BY %s", queryString, qb.getOrderBy(et1, qo))
	queryString = fmt.Sprintf("%s LIMIT %s OFFSET %s", queryString, qb.getLimit(qo), qb.getOffset(qo))

	return queryString, qpi
}

func (qb *QueryBuilder) CreateCountQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, qo *odata.QueryOptions) string {
	et1 := e1.GetEntityType()
	et2 := e1.GetEntityType()
	if e2 != nil { // 2nd entity is given, this means get e1 by e2
		et2 = e2.GetEntityType()
	}

	queryString := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", qb.tables[et1], qb.createJoin(e1, e2, false, nil, nil, ""))
	if id != nil {
		queryString = fmt.Sprintf("%s WHERE %s.%s = %v", queryString, tableMappings[et2], asMappings[et2][idField], id)
	}

	if qo != nil && !qo.QueryFilter.IsNil() {
		if id != nil {
			queryString = fmt.Sprintf("%s AND %s", queryString, qb.getFilterQueryString(et1, qo, false))
		} else {
			queryString = fmt.Sprintf("%s %s", queryString, qb.getFilterQueryString(et1, qo, true))
		}
	}

	return queryString
}
