package postgis

import (
	"fmt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/odata"
	"github.com/gost/godata"
	"strings"
	"time"
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
	if qo != nil && qo.Top != nil {
		return fmt.Sprintf("%v", *qo.Top)
	}
	return fmt.Sprintf("%v", qb.maxTop)
}

// getOffset returns the offset, this number is set by ODATA's
// $skip, if not provided do not skip anything = return "0"
func (qb *QueryBuilder) getOffset(qo *odata.QueryOptions) string {
	if qo != nil && qo.Skip != nil {
		return fmt.Sprintf("%v", *qo.Skip)
	}
	return "0"
}

// getOrderBy returns the string that needs to be placed after ORDER BY, this is set using
// ODATA's $orderby if not given use the default ORDER BY "table".id DESC
func (qb *QueryBuilder) getOrderBy(et entities.EntityType, qo *odata.QueryOptions) string {
	if qo != nil && qo.OrderBy != nil {
		obString := ""
		for _, obi := range qo.OrderBy.OrderByItems {
			propertyName := selectMappings[et][strings.ToLower(obi.Field.Value)]
			if len(obString) == 0 {
				obString = fmt.Sprintf("%s %s", propertyName, obi.Order)
			} else {
				obString = fmt.Sprintf("%s, %s %s", obString, propertyName, obi.Order)
			}
		}

		return obString
	}

	return fmt.Sprintf("%s DESC", selectMappings[et][idField])
}

// getSelect return the select string that needs to be placed after SELECT in the query
// select is set by ODATA's $select, if not set get all properties for the given entity (return all)
// addID to true if it needs to be added and isn't in QuerySelect.Params, addAs to true if a field needs to be
// outputted with AS [name]
func (qb *QueryBuilder) getSelect(et entities.Entity, qo *odata.QueryOptions, qpi *QueryParseInfo, addID bool, addAs bool, fromAs bool, isExpand bool, selectString string) string {
	var properties []string
	if qo == nil || qo.Select == nil || len(qo.Select.SelectItems) == 0 {
		properties = et.GetPropertyNames()
	} else {
		idAdded := false
		for _, p := range qo.Select.SelectItems {
			for _, s := range p.Segments {
				if s.Value == idField {
					idAdded = true
				}
			}

			for _, pn := range et.GetPropertyNames() {
				if strings.ToLower(p.Segments[0].Value) == strings.ToLower(pn) {
					if p.Segments[0].Value == idField {
						properties = append([]string{idField}, properties...)
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
	if isExpand {
		if et.GetEntityType() == entities.EntityTypeObservation {
			properties = append([]string{observationFeatureOfInterestID}, properties...)
		}

		if et.GetEntityType() == entities.EntityTypeDatastream {
			properties = append([]string{datastreamThingID, datastreamObservedPropertyID, datastreamSensorID}, properties...)
		}

		if et.GetEntityType() == entities.EntityTypeHistoricalLocation {
			properties = append([]string{historicalLocationThingID}, properties...)
		}

		if et.GetEntityType() == entities.EntityTypeObservation {
			properties = append([]string{observationStreamID}, properties...)
		}
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
			qo := &odata.QueryOptions{}
			if subQPI.ExpandItem != nil {
				qo = odata.ExpandItemToQueryOptions(subQPI.ExpandItem)
			}

			selectString = qb.getSelect(subQPI.Entity, qo, subQPI, true, true, true, false, selectString)
		}
	}

	return selectString
}

// addAsPrefix adds a prefix in front of the current as for example A_, B_ to be able to
// distinguish the different results if multiple tables are requested of the same type
func (qb *QueryBuilder) addAsPrefix(qpi *QueryParseInfo, as string) string {
	if qpi == nil {
		return as
	}

	return fmt.Sprintf("%v_%s", qpi.AsPrefix, as)
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
			nqo = &odata.QueryOptions{}
			nqo.Select = &godata.GoDataSelectQuery{SelectItems: []*godata.SelectItem{{Segments: []*godata.Token{{Value: "id"}}}}}

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
			qo := &odata.QueryOptions{}
			if subQPI.ExpandItem != nil {
				qo = odata.ExpandItemToQueryOptions(subQPI.ExpandItem)
			}
			joinString = qb.createJoin(subQPI.Parent.Entity, subQPI.Entity, true, qo, subQPI, joinString)
		}
	}

	return joinString
}

func (qb *QueryBuilder) constructQueryParseInfo(operations []*godata.ExpandItem, main *QueryParseInfo, from *QueryParseInfo) {
	for _, o := range operations {
		currentFrom := from
		for i, t := range o.Path {
			nQPI := &QueryParseInfo{}
			et, _ := entities.EntityFromString(strings.ToLower(t.Value))

			if i == len(o.Path)-1 {
				nQPI.Init(et.GetEntityType(), main.GetNextQueryIndex(), currentFrom, o)
				main.SubEntities = append(main.SubEntities, nQPI)

				if o.Expand != nil && len(o.Expand.ExpandItems) > 0 {
					qb.constructQueryParseInfo(o.Expand.ExpandItems, main, nQPI)
				}
			} else {
				nQPI.Init(et.GetEntityType(), main.GetNextQueryIndex(), currentFrom, nil)
				main.SubEntities = append(main.SubEntities, nQPI)
			}

			currentFrom = nQPI
		}
	}
}

// createFilterQueryString converts an OData query string found in odata.QueryOptions.QueryFilter to a PostgreSQL query string
// ParamFactory is used for converting SensorThings parameter names to postgres field names
// Convert receives a name such as phenomenonTime and returns "data ->> 'id'" true, returns
// false if parameter cannot be converted
func (qb *QueryBuilder) getFilterQueryString(et entities.EntityType, qo *odata.QueryOptions, addWhere bool) string {
	q := ""
	if qo != nil && qo.Filter != nil {
		if addWhere {
			q += " WHERE "
		}

		q += qb.createFilter(qo.Filter.Tree)
	}

	return q
}

func (qb *QueryBuilder) createFilter(pn *godata.ParseNode) string {
	output := ""
	switch pn.Token.Type {
	case godata.FilterTokenOpenParen:
		return ""
	case godata.FilterTokenCloseParen:
		return ""
	case godata.FilterTokenWhitespace:
		return ""
	case godata.FilterTokenNav:
		q := ""
		for i, part := range pn.Children {
			if i == 0 {
				q += fmt.Sprintf("%v ", strings.ToLower(qb.createFilter(part)))
				continue
			}

			arrow := "->"
			if i+1 == len(pn.Children) {
				arrow = "->>"
			}
			q += fmt.Sprintf("%v '%v'", arrow, qb.createFilter(part))
		}
		return q
	case godata.FilterTokenColon:
		return ""
	case godata.FilterTokenComma: //5
		return ""
	case godata.FilterTokenLogical:
		left := qb.createFilter(pn.Children[0])
		right := qb.createFilter(pn.Children[1])
		right = qb.prepareFilterRight(left, right)
		return fmt.Sprintf("%v %v %v", strings.ToLower(left), qb.odataLogicalOperatorToPostgreSQL(pn.Token.Value), right)
	case godata.FilterTokenOp:
		return ""
	case godata.FilterTokenFunc:
		return ""
	case godata.FilterTokenLambda:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenNull: // 10
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenIt:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenRoot:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenFloat:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenInteger:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenString: //15
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenDate:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenTime:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenDateTime:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenBoolean:
		return fmt.Sprintf("%v", pn.Token.Value)
	case godata.FilterTokenLiteral: // 20
		return fmt.Sprintf("%v", pn.Token.Value)
	}

	return output
}

func (qb *QueryBuilder) prepareFilterRight(left, right string) string {
	e := strings.Replace(fmt.Sprintf("%v", right), "'", "", -1)
	property := strings.ToLower(fmt.Sprintf("%v", left))

	if property == "encodingtype" {
		et, err := entities.CreateEncodingType(e)
		if err == nil {
			right = fmt.Sprintf("%v", et.Code)
		}
		return right
	}

	if property == "observationtype" {
		et, err := entities.GetObservationTypeByValue(e)
		if err == nil {
			right = fmt.Sprintf("%v", et.Code)
		}
		return right
	}

	if property == "phenomenontime" || property == "resulttime" || property == "time" {
		if t, err := time.Parse(time.RFC3339Nano, e); err == nil {
			right = fmt.Sprintf("'%s'", t.UTC().Format("2006-01-02T15:04:05.000Z"))
		}
		return right
	}

	return right
}

// odataLogicalOperatorToPostgreSQL converts a logical operator to a PostgreSQL string representation
func (qb *QueryBuilder) odataLogicalOperatorToPostgreSQL(o string) string {
	switch o {
	case "and":
		return "AND"
	case "or":
		return "OR"
	case "not":
		return "NOT"
	case "has":
		return "HAS"
	case "eq":
		return "="
	case "ne":
		return "!="
	case "gt":
		return ">"
	case "ge":
		return ">="
	case "lt":
		return "<"
	case "le":
		return "<="
	}

	return ""
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

	eo := &godata.ExpandItem{}
	if qo != nil {
		eo.Filter = qo.Filter
		eo.Expand = qo.Expand
		eo.OrderBy = qo.OrderBy
		eo.Search = qo.Search
		eo.Select = qo.Select
		eo.Skip = qo.Skip
		eo.Top = qo.Top
	}

	qpi := &QueryParseInfo{}
	qpi.Init(et1, 0, nil, eo)

	if qo != nil && qo.Expand != nil {
		qpi.SubEntities = make([]*QueryParseInfo, 0)
		if len(qo.Expand.ExpandItems) > 0 {
			qb.constructQueryParseInfo(qo.Expand.ExpandItems, qpi, qpi)
		}
	}

	etIdField := fmt.Sprintf("%s.id", tableMappings[et1])
	orderBy := qb.getOrderBy(et1, qo)
	queryString := fmt.Sprintf("SELECT %s FROM %s %s", qb.getSelect(e1, qo, qpi, true, true, false, false, ""), qb.tables[et1], qb.createJoin(e1, e2, false, qo, qpi, ""))
	queryString = fmt.Sprintf("%s WHERE %s IN (SELECT %s FROM %s %s", queryString, etIdField, etIdField, qb.tables[et1], qb.createJoin(e1, e2, false, qo, qpi, ""))

	if id != nil {
		if e2 == nil {
			queryString = fmt.Sprintf("%s WHERE %s = %v", queryString, selectMappings[et2][idField], id)
		} else {
			queryString = fmt.Sprintf("%s WHERE %s.%s = %v", queryString, tableMappings[et2], asMappings[et2][idField], id)
		}
	}

	if qo != nil && qo.Filter != nil {
		if id != nil {
			queryString = fmt.Sprintf("%s AND %s", queryString, qb.getFilterQueryString(et1, qo, false))
		} else {
			queryString = fmt.Sprintf("%s %s", queryString, qb.getFilterQueryString(et1, qo, true))
		}
	}

	limit := ""
	if qo != nil && qo.Top != nil && int(*qo.Top) != -1 {
		limit = fmt.Sprintf("LIMIT %s", qb.getLimit(qo))
	}
	queryString = fmt.Sprintf("%s ORDER BY %s %s OFFSET %s)", queryString, orderBy, limit, qb.getOffset(qo))
	queryString = fmt.Sprintf("%s ORDER BY %s", queryString, orderBy)

	return queryString, qpi
}

// CreateCountQuery creates the correct count query based on the given info
//   e1: entity to get
//   e2: from entity
//   id: e2 == nil: where e1.id = ... | e2 != nil: where e2.id = ...
// Returns an empty string if ODATA Query Count is set to false.
// example: Datastreams(1)/Thing = CreateCountQuery(&entities.Thing, &entities.Datastream, 1, nil)
func (qb *QueryBuilder) CreateCountQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, qo *odata.QueryOptions) string {
	if qo != nil && qo.Count != nil && bool(*qo.Count) == false {
		return ""
	}

	et1 := e1.GetEntityType()
	et2 := e1.GetEntityType()
	if e2 != nil { // 2nd entity is given, this means get e1 by e2
		et2 = e2.GetEntityType()
	}

	queryString := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", qb.tables[et1], qb.createJoin(e1, e2, false, nil, nil, ""))
	if id != nil {
		queryString = fmt.Sprintf("%s WHERE %s.%s = %v", queryString, tableMappings[et2], asMappings[et2][idField], id)
	}

	if qo != nil && qo.Filter != nil {
		if id != nil {
			queryString = fmt.Sprintf("%s AND %s", queryString, qb.getFilterQueryString(et1, qo, false))
		} else {
			queryString = fmt.Sprintf("%s %s", queryString, qb.getFilterQueryString(et1, qo, true))
		}
	}

	return queryString
}
