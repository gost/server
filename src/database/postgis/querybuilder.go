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
	joins  map[entities.EntityType]map[entities.EntityType]string
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

	qb.joins = createJoinMappings(qb.tables)
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
	if qo != nil && !qo.QueryTop.IsNil() {
		return fmt.Sprintf("%v", qo.QueryTop.Limit)
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
func (qb *QueryBuilder) getSelect(et entities.Entity, qo *odata.QueryOptions, addID bool, addAs bool, selectString string) string {
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
					properties = append(properties, pn)
				}
			}
		}
		if addID && !idAdded {
			properties = append(properties, "id")
		}
	}

	for _, p := range properties {
		toAdd := ""
		if len(selectString) > 0 {
			toAdd += ", "
		}
		entityType := et.GetEntityType()
		field := selectMappings[entityType][strings.ToLower(p)]
		if addAs {
			selectString += fmt.Sprintf("%s%s AS %s", toAdd, field, asMappings[entityType][strings.ToLower(p)])
		} else {
			selectString += fmt.Sprintf("%s%s", toAdd, field)
		}
	}

	if qo != nil && !qo.QueryExpand.IsNil() {
		for _, o := range qo.QueryExpand.Operations {
			selectString = qb.getSelect(o.Entity, o.QueryOptions, addID, addAs, selectString)
		}
	}

	return selectString
}

func (qb *QueryBuilder) createLateralJoin(e1 entities.Entity, e2 entities.Entity, isExpand bool, qo *odata.QueryOptions, joinString string) string {
	if e2 != nil {
		nqo := qo
		if !isExpand {
			nqo = &odata.QueryOptions{
				QuerySelect: &odata.QuerySelect{Params: []string{"id"}},
			}
		}

		et2 := e2.GetEntityType()
		joinString = fmt.Sprintf("%s"+
			"INNER JOIN LATERAL ("+
			"SELECT %s FROM %s %s "+
			"%s"+
			"ORDER BY %s "+
			"LIMIT %s OFFSET %s) AS %s on true", joinString,
			qb.getSelect(e2, nqo, true, false, ""),
			qb.tables[et2],
			qb.joins[e1.GetEntityType()][et2],
			qb.getFilterQueryString(et2, nqo, true),
			qb.getOrderBy(et2, nqo),
			qb.getLimit(nqo),
			qb.getOffset(nqo),
			qb.removeSchema(qb.tables[et2]))
	} else {
		if qo != nil && !qo.QueryExpand.IsNil() {
			for _, qe := range qo.QueryExpand.Operations {
				joinString = qb.createLateralJoin(e1, qe.Entity, true, qe.QueryOptions, joinString)
			}
		}
	}

	return joinString
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
//   e1 = entity to get
//   e2 = from entity
//   id = where e2
// example: Datastreams(1)/Thing = CreateQuery(&entities.Thing, &entities.Datastream, 1, nil)
func (qb *QueryBuilder) CreateQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, qo *odata.QueryOptions) (string, *QueryParseInfo, error) {
	et1 := e1.GetEntityType()
	et2 := e1.GetEntityType()
	if e2 != nil { // 2nd entity is given, this means get e1 by e2
		et2 = e2.GetEntityType()
	}

	queryString := fmt.Sprintf("SELECT %s FROM %s %s", qb.getSelect(e1, qo, false, true, ""), qb.tables[et1], qb.createLateralJoin(e1, e2, false, qo, ""))
	if id != nil {
		queryString = fmt.Sprintf("%s WHERE %s = %v", queryString, selectMappings[et2]["id"], id)
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

	qpi := &QueryParseInfo{}
	qpi.Init(entities.EntityTypeDatastream, 0, DatastreamParamFactory)

	return queryString, qpi, nil
}

// Test is a temporarily test function while developing
func (qb *QueryBuilder) Test() {
	fmt.Println("------------GET THINGS------------")
	sql1, _, _ := qb.CreateQuery(&entities.Thing{}, nil, nil, nil)
	fmt.Println(sql1)

	fmt.Println("------------GET THING WITH SELECT BY DATASTREAM------------")
	qo2 := &odata.QueryOptions{}
	qo2.QuerySelect = &odata.QuerySelect{}
	qo2.QuerySelect.Parse("name,description")
	sql2, _, _ := qb.CreateQuery(&entities.Thing{}, &entities.Datastream{}, 1, qo2)
	fmt.Println(sql2)

	fmt.Println("------------GET DATASTREAMS WITH SELECT, EXPAND THING WITH SELECT------------")
	qo31 := &odata.QueryOptions{}
	qo31.QuerySelect = &odata.QuerySelect{}
	qo31.QuerySelect.Parse("name,description")
	qo31.QueryExpand = &odata.QueryExpand{}
	qo31.QueryExpand.Parse("Thing($select=name)")
	sql3, _, _ := qb.CreateQuery(&entities.Datastream{}, nil, nil, qo31)
	fmt.Println(sql3)

	fmt.Println("------------GET THING BY LOCATION------------")
	sql4, _, _ := qb.CreateQuery(&entities.Thing{}, &entities.Location{}, 1, nil)
	fmt.Println(sql4)

	fmt.Println("------------GET HISTORICAL LOCATION BY THING ------------")
	sql5, _, _ := qb.CreateQuery(&entities.HistoricalLocation{}, &entities.Thing{}, 1, nil)
	fmt.Println(sql5)

	fmt.Println("------------GET LOCATION BY THING ------------")
	sql6, _, _ := qb.CreateQuery(&entities.Location{}, &entities.Thing{}, 1, nil)
	fmt.Println(sql6)

	fmt.Println("------------GET HISTORICAL LOCATION BY LOCATION ------------")
	sql7, _, _ := qb.CreateQuery(&entities.HistoricalLocation{}, &entities.Location{}, 1, nil)
	fmt.Println(sql7)

	fmt.Println("------------GET OBSERVATIONS------------")
	sql8, _, _ := qb.CreateQuery(&entities.Observation{}, nil, nil, nil)
	fmt.Println(sql8)
}
