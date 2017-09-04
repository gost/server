package postgis

import (
	"fmt"
	"strings"
	"time"

	entities "github.com/gost/core"
	"github.com/gost/godata"
	"github.com/gost/server/sensorthings/odata"
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
	if qo != nil && qo.OrderBy != nil && len(qo.OrderBy.OrderByItems) > 0 {
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

			if qo.Select != nil && qo.Select.SelectItems != nil && len(qo.Select.SelectItems) > 0 {
				// skip if a qo has a select with nil (non user requested expand)
				if qo.Select.SelectItems[0].Segments[0].Value == "nil" {
					continue
				}
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

func (qb *QueryBuilder) createJoin(e1 entities.Entity, e2 entities.Entity, id interface{}, isExpand, generatedExpand bool, qo *odata.QueryOptions, qpi *QueryParseInfo, joinString string) string {
	if e2 != nil {
		nqo := qo
		et2 := e2.GetEntityType()

		asPrefix := ""
		if qpi != nil {
			if qpi.Parent != nil {
				asPrefix = qpi.Parent.AsPrefix
			} else {
				asPrefix = qpi.AsPrefix
			}
		}

		if !isExpand {
			nqo = &odata.QueryOptions{}
			nqo.Select = &godata.GoDataSelectQuery{SelectItems: []*godata.SelectItem{{Segments: []*godata.Token{{Value: "id"}}}}}
			if id != nil {
				var err error
				nqo.Filter, err = godata.ParseFilterString(fmt.Sprintf("id eq %v", id))
				if err != nil {
					fmt.Printf("\n\n ERROR %v \n\n", err)
				}
			}

			join := getJoin(qb.tables, et2, e1.GetEntityType(), asPrefix)
			lowerJoin := strings.ToLower(join)
			filterPrefix := "WHERE"
			if strings.Contains(lowerJoin, "where") {
				filterPrefix = "AND"
			}

			filter := qb.getFilterQueryString(et2, nqo, filterPrefix)

			joinString = fmt.Sprintf("%s"+
				"INNER JOIN LATERAL ("+
				"SELECT %s FROM %s %s "+
				"%s) "+
				"AS %s on true ", joinString,
				qb.getSelect(e2, nqo, nil, true, true, false, false, ""),
				qb.tables[et2],
				join,
				filter,
				tableMappings[et2])
		} else {
			join := getJoin(qb.tables, et2, e1.GetEntityType(), asPrefix)
			lowerJoin := strings.ToLower(join)
			filterPrefix := "WHERE"
			if strings.Contains(lowerJoin, "where") {
				filterPrefix = "AND"
			}

			joinType := "LEFT JOIN LATERAL"
			if generatedExpand {
				joinType = "INNER JOIN LATERAL"
			}

			joinString = fmt.Sprintf("%s"+
				"%s ("+
				"SELECT %s FROM %s %s "+
				"%s "+
				"ORDER BY %s "+
				"LIMIT %s OFFSET %s) AS %s on true ",
				joinString,
				joinType,
				qb.getSelect(e2, nqo, qpi, true, true, false, true, ""),
				qb.tables[et2],
				join,
				qb.getFilterQueryString(et2, nqo, filterPrefix),
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

			// if first select value is nil means the expand is not requested by the user so
			// supply qo to createJoin as a non Expand
			generatedExpand := isExpandGenerated(qo.Select)
			if !generatedExpand {
				qo.Select = &godata.GoDataSelectQuery{}
			}
			joinString = qb.createJoin(subQPI.Parent.Entity, subQPI.Entity, nil, true, generatedExpand, qo, subQPI, joinString)
		}
	}

	return joinString
}

func isExpandGenerated(sq *godata.GoDataSelectQuery) bool {
	generatedExpand := false
	if sq != nil && sq.SelectItems != nil && len(sq.SelectItems) > 0 {
		generatedExpand = sq.SelectItems[0].Segments[0].Value == "nil"
	}

	return generatedExpand
}

func (qb *QueryBuilder) constructQueryParseInfo(operations []*godata.ExpandItem, main *QueryParseInfo) {
	for _, o := range operations {
		for i, t := range o.Path {
			nQPI := &QueryParseInfo{}
			et, _ := entities.EntityFromString(strings.ToLower(t.Value))

			path := make([]entities.EntityType, 0)
			for p := 0; p < i+1; p++ {
				etfs, _ := entities.EntityFromString(o.Path[p].Value)
				path = append(path, etfs.GetEntityType())
			}

			parent := main.GetParent(path)
			exist, qi := main.QueryInfoExists(path)
			if exist {
				if i == len(o.Path)-1 {
					if qi.ExpandItem == nil {
						qi.ExpandItem = &godata.ExpandItem{
							Filter: o.Filter,
						}
					} else if qi.ExpandItem.Filter != nil && qi.ExpandItem.Filter.Tree != nil {

						temp := &godata.ParseNode{}
						*temp = *qi.ExpandItem.Filter.Tree
						fq := &godata.GoDataFilterQuery{
							Tree: &godata.ParseNode{
								Token: &godata.Token{
									Type:  6,
									Value: "and",
								},
								Children: []*godata.ParseNode{
									temp,
									o.Filter.Tree,
								},
							},
						}

						qi.ExpandItem.Filter = fq

					} else {
						qi.ExpandItem.Filter = o.Filter

					}
				}

				continue
			}

			if i == len(o.Path)-1 {
				nQPI.Init(et.GetEntityType(), main.GetNextQueryIndex(), parent, o)
				main.SubEntities = append(main.SubEntities, nQPI)

				if o.Expand != nil && len(o.Expand.ExpandItems) > 0 {
					qb.constructQueryParseInfo(o.Expand.ExpandItems, main)
				}
			} else {
				var ei *godata.ExpandItem
				ei = nil

				// if expand was generated and there is a select with nil value forward it into the inner expand
				if isExpandGenerated(o.Select) {
					ei = &godata.ExpandItem{
						Select: o.Select,
					}
				}

				nQPI.Init(et.GetEntityType(), main.GetNextQueryIndex(), parent, ei)
				main.SubEntities = append(main.SubEntities, nQPI)
			}

			continue
		}
	}
}

// createFilterQueryString converts an OData query string found in odata.QueryOptions.QueryFilter to a PostgreSQL query string
// ParamFactory is used for converting SensorThings parameter names to postgres field names
// Convert receives a name such as phenomenonTime and returns "data ->> 'id'" true, returns
// false if parameter cannot be converted
func (qb *QueryBuilder) getFilterQueryString(et entities.EntityType, qo *odata.QueryOptions, prefix string) string {
	q := ""
	if qo != nil && qo.Filter != nil {
		filterString := qb.createFilter(et, qo.Filter.Tree, false)
		if filterString == "" {
			return q
		}

		q += fmt.Sprintf("%s ", prefix)
		q += filterString
	}

	return q
}

func (qb *QueryBuilder) createFilter(et entities.EntityType, pn *godata.ParseNode, ignoreSelectAs bool) string {
	output := ""
	if pn == nil || pn.Token == nil {
		return output
	}

	switch pn.Token.Type {
	case godata.FilterTokenNav:
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
	case godata.FilterTokenLogical:
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
	case godata.FilterTokenFunc:
		if pn.Token.Value == "contains" {
			left := qb.createFilter(et, pn.Children[0], true)
			right := qb.createFilter(et, pn.Children[1], true)
			return fmt.Sprintf("%s LIKE %s", qb.createLike(left, LikeContains), qb.createLike(right, LikeContains))
		} else if pn.Token.Value == "substringof" {
			left := qb.createFilter(et, pn.Children[0], true)
			right := qb.createFilter(et, pn.Children[1], true)
			return fmt.Sprintf("%s LIKE %s", qb.createLike(right, LikeContains), qb.createLike(left, LikeContains))
		} else if pn.Token.Value == "endswith" {
			left := qb.createFilter(et, pn.Children[0], true)
			right := qb.createFilter(et, pn.Children[1], true)
			return fmt.Sprintf("%s LIKE %s", qb.createLike(left, LikeEndsWith), qb.createLike(right, LikeEndsWith))
		} else if pn.Token.Value == "startswith" {
			left := qb.createFilter(et, pn.Children[0], true)
			right := qb.createFilter(et, pn.Children[1], true)
			return fmt.Sprintf("%s LIKE %s", qb.createLike(left, LikeStartsWith), qb.createLike(right, LikeStartsWith))
		}
		if pn.Token.Value == "length" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("LENGTH(%s)", left)
		} else if pn.Token.Value == "indexof" {
			left := qb.createFilter(et, pn.Children[0], true)
			right := qb.createFilter(et, pn.Children[1], true)
			return fmt.Sprintf("STRPOS(%s, %s) -1", left, right)
		} else if pn.Token.Value == "substring" {
			left, right, right2 := "", "", ""
			left = qb.createFilter(et, pn.Children[0], true)
			right = qb.createFilter(et, pn.Children[1], true)
			if len(pn.Children) > 2 {
				right2 = qb.createFilter(et, pn.Children[2], true)
				return fmt.Sprintf("SUBSTRING(%s from (%s + 1) for %s)", left, right, right2)
			}

			return fmt.Sprintf("SUBSTRING(%s from (%s + 1) for LENGTH(%s))", left, right, left)
		} else if pn.Token.Value == "tolower" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("LOWER(%s)", left)
		} else if pn.Token.Value == "toupper" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("UPPER(%s)", left)
		} else if pn.Token.Value == "trim" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("TRIM(both ' ' from %s)", left)
		} else if pn.Token.Value == "concat" {
			left := qb.createFilter(et, pn.Children[0], true)
			right := qb.createFilter(et, pn.Children[1], true)
			return fmt.Sprintf("CONCAT(%s, %s)", left, right)
		} else if pn.Token.Value == "round" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("ROUND(CAST(%s as double precision))", strings.Replace(left, "->", "->>", -1))
		} else if pn.Token.Value == "floor" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("FLOOR(CAST(%s as double precision))", strings.Replace(left, "->", "->>", -1))
		} else if pn.Token.Value == "ceiling" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("CEILING(CAST(%s as double precision))", strings.Replace(left, "->", "->>", -1))
		} else if pn.Token.Value == "year" {
			return qb.createExtractDateQuery(pn, et, "YEAR")
		} else if pn.Token.Value == "month" {
			return qb.createExtractDateQuery(pn, et, "MONTH")
		} else if pn.Token.Value == "day" {
			return qb.createExtractDateQuery(pn, et, "DAY")
		} else if pn.Token.Value == "hour" {
			return qb.createExtractDateQuery(pn, et, "HOUR")
		} else if pn.Token.Value == "minute" {
			return qb.createExtractDateQuery(pn, et, "MINUTE")
		} else if pn.Token.Value == "second" {
			return qb.createExtractDateQuery(pn, et, "SECOND")
		} else if pn.Token.Value == "fractionalseconds" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("EXTRACT(MICROSECONDS FROM to_timestamp(%s,'YYYY-MM-DD\"T\"HH24:MI:SS.MS\"Z\"')) / 1000000", left)
		} else if pn.Token.Value == "date" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("(%s)::date", left)
		} else if pn.Token.Value == "time" {
			left := qb.createFilter(et, pn.Children[0], true)
			if strings.Contains(strings.ToLower(left), "time") {
				return fmt.Sprintf("((%s)::timestamp)::time", left)
			}

			return fmt.Sprintf("(%s)::time", left)
		} else if pn.Token.Value == "totaloffsetminutes" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("EXTRACT(TIMEZONE_MINUTE FROM to_timestamp(%s,'YYYY-MM-DD\"T\"HH24:MI:SS.MS\"Z\"'))", left)
		} else if pn.Token.Value == "now" {
			return fmt.Sprint("to_char(now()::timestamp at time zone 'UTC', 'YYYY-MM-DD\"T\"HH24:MI:SS.MS\"Z\"')")
		} else if pn.Token.Value == "maxdatetime" {
			return fmt.Sprint("'9999-12-31T23:59:59.999Z'")
		} else if pn.Token.Value == "mindatetime" {
			return fmt.Sprint("'0001-01-01T00:00:00.000Z'")
		} else if pn.Token.Value == "totalseconds" {
			left := qb.createFilter(et, pn.Children[0], true)
			return fmt.Sprintf("SELECT extract(epoch from (%s)::timestamp)", left)
		} else if pn.Token.Value == "geo.distance" {
			return qb.createSpatialQuery(pn, et, "ST_DISTANCE(%s, %s)", 2)
		} else if pn.Token.Value == "geo.length" {
			return qb.createSpatialQuery(pn, et, "ST_LENGTH(%s)", 1)
		} else if pn.Token.Value == "st_equals" {
			return qb.createSpatialQuery(pn, et, "ST_EQUALS(%s, %s)", 2)
		} else if pn.Token.Value == "st_touches" {
			return qb.createSpatialQuery(pn, et, "ST_TOUCHES(%s, %s)", 2)
		} else if pn.Token.Value == "st_overlaps" {
			return qb.createSpatialQuery(pn, et, "ST_OVERLAPS(%s, %s)", 2)
		} else if pn.Token.Value == "st_crosses" {
			return qb.createSpatialQuery(pn, et, "ST_CROSSES(%s, %s)", 2)
		} else if pn.Token.Value == "st_contains" {
			return qb.createSpatialQuery(pn, et, "ST_CONTAINS(%s, %s)", 2)
		} else if pn.Token.Value == "st_disjoint" {
			return qb.createSpatialQuery(pn, et, "ST_DISJOINT(%s, %s)", 2)
		} else if pn.Token.Value == "st_relate" {
			return qb.createSpatialQuery(pn, et, "ST_RELATE(%s, %s, %s)", 3)
		} else if pn.Token.Value == "st_within" {
			return qb.createSpatialQuery(pn, et, "ST_WITHIN(%s, %s)", 2)
		} else if pn.Token.Value == "st_intersects" || pn.Token.Value == "geo.intersects" {
			return qb.createSpatialQuery(pn, et, "ST_INTERSECTS(%s, %s)", 2)
		}
	case godata.FilterTokenOp:
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
	case godata.FilterTokenGeography:
		return fmt.Sprintf("ST_GeomFromText(%v)", pn.Children[0].Token.Value)
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
		p := selectMappings[et][strings.ToLower(pn.Token.Value)]
		if p != "" {
			return p
		}

		if ignoreSelectAs && selectMappingsIgnore[et][pn.Token.Value] {
			return fmt.Sprintf("%s.%v", et.ToString(), pn.Token.Value)
		}

		return pn.Token.Value
	}

	return output
}

func (qb *QueryBuilder) prepareFilter(et entities.EntityType, originalLeft, left, originalRight, right string) (string, string) {
	for i := 0; i < 2; i++ {
		var oStr []*string
		var str []*string
		if i == 0 {
			oStr = []*string{&originalLeft, &originalRight}
			str = []*string{&left, &right}
		} else {
			oStr = []*string{&originalRight, &originalLeft}
			str = []*string{&right, &left}
		}

		e := strings.Replace(fmt.Sprintf("%v", *oStr[1]), "'", "", -1)
		property := strings.ToLower(fmt.Sprintf("%v", *oStr[0]))

		if property == "encodingtype" {
			et, err := entities.CreateEncodingType(e)
			if err == nil {
				*str[1] = fmt.Sprintf("%v", et.Code)
			}
			return left, right
		}

		if property == "observationtype" {
			et, err := entities.GetObservationTypeByValue(e)
			if err == nil {
				*str[1] = fmt.Sprintf("%v", et.Code)
			}
			return left, right
		}

		if property == "phenomenontime" || property == "resulttime" || property == "time" {
			if t, err := time.Parse(time.RFC3339Nano, e); err == nil {
				*str[1] = fmt.Sprintf("'%s'", t.UTC().Format("2006-01-02T15:04:05.000Z"))
			} else {
			}

			return left, right
		}
	}

	return left, right
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

func (qb *QueryBuilder) createArithmetic(et entities.EntityType, pn *godata.ParseNode, operator, castTo string) string {
	left := qb.CastObservationResult(qb.createFilter(et, pn.Children[0], true), castTo)
	right := qb.CastObservationResult(qb.createFilter(et, pn.Children[1], true), castTo)

	return fmt.Sprintf("%s %s %s", left, operator, right)
}

// CastObservationResult converts an observation result query to a specified type (castTo)
func (qb *QueryBuilder) CastObservationResult(input string, castTo string) string {
	if input == "observation.data -> 'result'" {
		return fmt.Sprintf("(observation.data ->> 'result')::%s", castTo)
	}

	return input
}

func (qb *QueryBuilder) createExtractDateQuery(pn *godata.ParseNode, et entities.EntityType, function string) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("EXTRACT(%s FROM to_timestamp(%s,'YYYY-MM-DD\"T\"HH24:MI:SS.MS\"Z\"'))", function, left)
}

func (qb *QueryBuilder) createSpatialQuery(pn *godata.ParseNode, et entities.EntityType, function string, params int) string {
	if params == 1 {
		return fmt.Sprintf(function, qb.addGeomFromGeoJSON(pn.Children[0].Token, qb.createFilter(et, pn.Children[0], true)))
	} else if params == 2 {
		return fmt.Sprintf(function,
			qb.addGeomFromGeoJSON(pn.Children[0].Token, qb.createFilter(et, pn.Children[0], true)),
			qb.addGeomFromGeoJSON(pn.Children[1].Token, qb.createFilter(et, pn.Children[1], true)))
	} else if params == 3 {
		return fmt.Sprintf(function,
			qb.addGeomFromGeoJSON(pn.Children[0].Token, qb.createFilter(et, pn.Children[0], true)),
			qb.addGeomFromGeoJSON(pn.Children[1].Token, qb.createFilter(et, pn.Children[1], true)),
			qb.createFilter(et, pn.Children[2], true))
	}

	return function
}

func (qb *QueryBuilder) addGeomFromGeoJSON(token *godata.Token, input string) string {
	if token.Type != godata.FilterTokenGeography {
		input = fmt.Sprintf("ST_GeomFromGeoJSON(%s)", input)
	}

	return input
}

// LikeType describes the type of like
type LikeType int

// LikeType is a "enumeration" of the Like types, LikeStartsWith = startsWith input%, LikeEndsWith = endsWith %input, LikeContains = contains %input%
const (
	LikeStartsWith LikeType = 0
	LikeEndsWith   LikeType = 1
	LikeContains   LikeType = 2
)

func (qb *QueryBuilder) createLike(input string, like LikeType) string {
	if !strings.HasPrefix(input, "'") || !strings.HasSuffix(input, "'") {
		return input
	}

	input = input[1 : len(input)-1]

	switch like {
	case LikeStartsWith:
		{
			return fmt.Sprintf("%s%s%s", "'", input, "%'")
		}
	case LikeEndsWith:
		{
			return fmt.Sprintf("%s%s%s", "'%", input, "'")
		}
	case LikeContains:
		{
			return fmt.Sprintf("%s%s%s", "'%", input, "%'")
		}
	}

	return input
}

func (qb *QueryBuilder) sortQueryOptions(qo *odata.QueryOptions) {
	// Check filters for query on expanded item (FilterTokenNAV) and move it to the desired expand
	// if expand is not requested create inner query
	if qo != nil && qo.Filter != nil {
		ce := make([]string, 0)
		qb.sortFilter(qo, qo.Filter.Tree, qo.Filter.Tree, nil, &ce)
		pn := &godata.ParseNode{}
		fq := &godata.GoDataFilterQuery{}
		fq.Tree = pn
		pn.Children = append(pn.Children, qo.Filter.Tree)

		qb.cleanupFilter(fq.Tree)

		if len(fq.Tree.Children) > 0 {
			*qo.Filter.Tree = *fq.Tree.Children[0]
		} else {
			*qo.Filter.Tree = godata.ParseNode{}
		}
	}
}

func (qb *QueryBuilder) sortFilter(qo *odata.QueryOptions, pn *godata.ParseNode, parentNode *godata.ParseNode, startNavNode *godata.ParseNode, currentExpand *[]string) {
	// navigational filter found
	if pn.Token != nil && pn.Token.Type == godata.FilterTokenNav {
		if startNavNode == nil {
			startNavNode = pn
		}
		if pn.Children[0].Token.Type == godata.FilterTokenNav {
			*currentExpand = append([]string{pn.Children[1].Token.Value}, *currentExpand...)
			qb.sortFilter(qo, pn.Children[0], parentNode, startNavNode, currentExpand)
		} else {
			*currentExpand = append([]string{pn.Children[1].Token.Value}, *currentExpand...)
			*currentExpand = append([]string{pn.Children[0].Token.Value}, *currentExpand...)

			cp := *currentExpand
			field := cp[len(cp)-1]
			cp = cp[:len(cp)-1]
			afterCouplingNode := findParseNodeAfterCoupling(parentNode)

			// set navigation node to literal with last item from path (which should be the fieldname i.e Datastreams/Observations/id)
			*startNavNode = godata.ParseNode{
				Token: &godata.Token{
					Value: field,
					Type:  godata.FilterTokenLiteral,
				},
				Parent: parentNode,
			}

			// check if expand exist for current path, if so add a new filter to it
			addExpand := true
			if qo.Expand != nil {
				for i := 0; i < len(qo.Expand.ExpandItems); i++ {
					e := qo.Expand.ExpandItems[i]
					if len(cp) != len(e.Path) {
						continue
					}

					for j := 0; j < len(cp); j++ {
						if cp[j] == e.Path[j].Value {
							if j == len(cp)-1 {
								addFilterToExpandItem(afterCouplingNode, e)
								addExpand = false
							}
						} else {
							break
						}
					}

					if !addExpand {
						break
					}
				}
			}

			// expand is not defined, create an expand with filter and define that the data should not be return to the requester
			if addExpand {
				// add new expand and set to not export since this data was not requested by the user
				if qo.Expand == nil {
					qo.Expand = &godata.GoDataExpandQuery{}
					qo.Expand.ExpandItems = make([]*godata.ExpandItem, 0)
				}

				expandString := ""
				for _, s := range cp {
					if len(expandString) != 0 {
						expandString += "/"
					}
					expandString += s
				}

				expand, _ := godata.ParseExpandString(expandString)
				newExpandItem := &godata.ExpandItem{
					Path: expand.ExpandItems[0].Path,
				}

				// add a select of nil, if set to nil the QueryBuilder knows that the expand is generated
				newExpandItem.Select = &godata.GoDataSelectQuery{
					SelectItems: []*godata.SelectItem{
						{
							Segments: []*godata.Token{
								{
									Value: "nil",
								},
							},
						},
					},
				}

				// add filter to expand
				addFilterToExpandItem(afterCouplingNode, newExpandItem)
				qo.Expand.ExpandItems = append(qo.Expand.ExpandItems, newExpandItem)
			}

			// remove node from filter since it sits inside the expand now
			*afterCouplingNode = godata.ParseNode{}
		}
	} else {
		// look for more navigational filters
		for i := 0; i < len(pn.Children); i++ {
			c := pn.Children[i]
			ne := make([]string, 0)
			c.Parent = pn
			qb.sortFilter(qo, c, pn, nil, &ne)
		}
	}
}

func (qb *QueryBuilder) cleanupFilter(pn *godata.ParseNode) {
	for i := len(pn.Children) - 1; i >= 0; i-- {
		current := pn.Children[i]
		current.Parent = pn
		if current.Children != nil && len(current.Children) > 0 {
			qb.cleanupFilter(current)
		}
	}

	for i := len(pn.Children) - 1; i >= 0; i-- {
		current := pn.Children[i]
		current.Parent = pn
		if current.Children != nil && len(current.Children) > 0 {
			if len(current.Children) == 2 {
				// If 2 childs are found with empty token remove the node from it's parent
				if current.Children[0].Token == nil && current.Children[1].Token == nil {
					pn.Children = append(pn.Children[:i], pn.Children[i+1:]...)
					continue
				}

				// If child is empty, push back child 1
				if current.Children[0].Token == nil && current.Children[1].Token != nil {
					current = current.Children[1]
				} else if current.Children[0].Token != nil && current.Children[1].Token == nil {
					current = current.Children[0]
				}
			}

			if len(current.Children) == 1 {
				if current.Children[0].Token == nil {
					pn.Children = append(pn.Children[:i], pn.Children[i+1:]...)
					continue
				}
			}
		}
	}
}

func addFilterToExpandItem(pn *godata.ParseNode, ei *godata.ExpandItem) {
	copyParseNode := &godata.ParseNode{}
	*copyParseNode = *pn

	//add filter to tree, if there is already a filter get if there is an and or or in front
	if ei.Filter == nil {
		ei.Filter = &godata.GoDataFilterQuery{Tree: copyParseNode}
	} else {
		temp := *ei.Filter.Tree
		cpn := findFirstCouplingParseNode(copyParseNode)
		cpn.Children = make([]*godata.ParseNode, 2)
		cpn.Children[0] = &temp
		cpn.Children[1] = copyParseNode
		ei.Filter = &godata.GoDataFilterQuery{Tree: cpn}
	}
}

func findParseNodeAfterCoupling(pn *godata.ParseNode) *godata.ParseNode {
	if pn.Parent == nil {
		return pn
	}

	tokenValue := strings.ToLower(pn.Parent.Token.Value)
	if tokenValue == "and" || tokenValue == "or" {
		return pn
	}

	return findParseNodeAfterCoupling(pn.Parent)
}

func findFirstCouplingParseNode(pn *godata.ParseNode) *godata.ParseNode {
	if pn.Parent == nil {
		return pn
	}

	tokenValue := strings.ToLower(pn.Parent.Token.Value)
	if tokenValue == "and" || tokenValue == "or" {
		return pn.Parent
	}

	return findFirstCouplingParseNode(pn.Parent)
}

// CreateCountQuery creates the correct count query based on the given info
//   e1: entity to get
//   e2: from entity
//   id: e2 == nil: where e1.id = ... | e2 != nil: where e2.id = ...
// Returns an empty string if ODATA Query Count is set to false.
// example: Datastreams(1)/Thing = CreateCountQuery(&entities.Thing, &entities.Datastream, 1, nil)
func (qb *QueryBuilder) CreateCountQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, queryOptions *odata.QueryOptions) string {
	var qo *odata.QueryOptions

	if queryOptions != nil {
		qo = &odata.QueryOptions{}
		*qo = *queryOptions
		qb.sortQueryOptions(qo)
	}

	et1 := e1.GetEntityType()
	et2 := e1.GetEntityType()
	if e2 != nil { // 2nd entity is given, this means get e1 by e2
		et2 = e2.GetEntityType()
	}

	queryString := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", qb.tables[et1], qb.createJoin(e1, e2, id, false, false, nil, nil, ""))
	if id != nil {
		queryString = fmt.Sprintf("%s WHERE %s.%s = %v", queryString, tableMappings[et2], asMappings[et2][idField], id)
	}

	if qo != nil && qo.Filter != nil {
		if id != nil {
			queryString = fmt.Sprintf("%s AND %s", queryString, qb.getFilterQueryString(et1, qo, ""))
		} else {
			queryString = fmt.Sprintf("%s %s", queryString, qb.getFilterQueryString(et1, qo, "WHERE"))
		}
	}

	return queryString
}

// CreateQuery creates a new count query based on given input
//   e1: entity to get
//   e2: from entity
//   id: e2 == nil: where e1.id = ... | e2 != nil: where e2.id = ...
// example: Datastreams(1)/Thing = CreateQuery(&entities.Thing, &entities.Datastream, 1, nil)
func (qb *QueryBuilder) CreateQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, queryOptions *odata.QueryOptions) (string, *QueryParseInfo) {
	var qo *odata.QueryOptions
	qo = nil

	if queryOptions != nil {
		qo = &odata.QueryOptions{}
		*qo = *queryOptions
		qb.sortQueryOptions(qo)
	}

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
			qb.constructQueryParseInfo(qo.Expand.ExpandItems, qpi)
		}
	}

	queryString := fmt.Sprintf("SELECT %s FROM (SELECT %s FROM %s",
		qb.getSelect(e1, qo, qpi, true, true, true, false, ""),
		qb.getSelect(e1, qo, nil, true, true, false, true, ""),
		qb.tables[et1],
	)

	if id != nil && e2 == nil {
		queryString = fmt.Sprintf("%s WHERE %s = %v", queryString, selectMappings[et2][idField], id)
	}

	if qo != nil && qo.Filter != nil {
		if id != nil {
			if e2 == nil {
				queryString = fmt.Sprintf("%s AND %s", queryString, qb.getFilterQueryString(et1, qo, ""))
			} else {
				queryString = fmt.Sprintf("%s %s", queryString, qb.getFilterQueryString(et1, qo, "WHERE"))
			}
		} else {
			queryString = fmt.Sprintf("%s %s", queryString, qb.getFilterQueryString(et1, qo, "WHERE"))
		}
	}

	limit := ""
	if qo != nil && qo.Top != nil && int(*qo.Top) != -1 {
		limit = fmt.Sprintf("LIMIT %s", qb.getLimit(qo))
	}
	queryString = fmt.Sprintf("%s ORDER BY %s )", queryString, qb.getOrderBy(et1, qo))
	queryString = fmt.Sprintf("%s AS %s %s %s OFFSET %s",
		queryString,
		qb.addAsPrefix(qpi, tableMappings[et1]),
		qb.createJoin(e1, e2, id, false, false, qo, qpi, ""),
		limit,
		qb.getOffset(qo),
	)

	//fmt.Printf("%s\n", queryString)
	return queryString, qpi
}
