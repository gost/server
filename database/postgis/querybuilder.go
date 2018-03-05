package postgis

import (
	"fmt"
	"strings"
	"time"

	entities "github.com/gost/core"
	"github.com/gost/godata"
	gostLog "github.com/gost/server/log"
	"github.com/gost/server/sensorthings/odata"
	log "github.com/sirupsen/logrus"
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
// supply an int to extra to ad to the limit
func (qb *QueryBuilder) getLimit(qo *odata.QueryOptions, extra int) int {
	if qo != nil && qo.Top != nil {
		return int(*qo.Top) + extra
	}
	return qb.maxTop + extra
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
func (qb *QueryBuilder) getOrderBy(et entities.EntityType, qo *odata.QueryOptions, fromAs bool) string {
	if qo != nil && qo.OrderBy != nil && len(qo.OrderBy.OrderByItems) > 0 {
		obString := ""
		for _, obi := range qo.OrderBy.OrderByItems {
			propertyName := ""
			if fromAs {
				propertyName = asMappings[et][strings.ToLower(obi.Field.Value)]
			} else {
				propertyName = selectMappings[et][strings.ToLower(obi.Field.Value)]
			}

			if len(obString) == 0 {
				obString = fmt.Sprintf("%s %s", propertyName, obi.Order)
			} else {
				obString = fmt.Sprintf("%s, %s %s", obString, propertyName, obi.Order)
			}
		}

		return obString
	}

	if fromAs {
		return fmt.Sprintf("%s DESC", asMappings[et][idField])
	}

	return fmt.Sprintf("%s DESC", selectMappings[et][idField])
}

// getSelect return the select string that needs to be placed after SELECT in the query
// select is set by ODATA's $select, if not set get all properties for the given entity (return all)
// addID to true if it needs to be added and isn't in QuerySelect.Params, addAs to true if a field needs to be
// outputted with AS [name]
func (qb *QueryBuilder) getSelect(et entities.Entity, qo *odata.QueryOptions, qpi *QueryParseInfo, addID bool, addAs bool, fromAs bool, isExpand bool, selectString string) string {
	entityType := et.GetEntityType()
	properties := qb.getProperties(et, qo, addID)

	// ToDo: this is a fix for supporting $expand=Observations/FeatureOfInterest, try to add observationFeatureOfInterestID in a different way
	if isExpand {
		if entityType == entities.EntityTypeObservation {
			properties = append([]string{observationFeatureOfInterestID}, properties...)
		}

		if entityType == entities.EntityTypeDatastream {
			properties = append([]string{datastreamThingID, datastreamObservedPropertyID, datastreamSensorID}, properties...)
		}

		if entityType == entities.EntityTypeHistoricalLocation {
			properties = append([]string{historicalLocationThingID}, properties...)
		}

		if entityType == entities.EntityTypeObservation {
			properties = append([]string{observationStreamID}, properties...)
		}
	}

	selectString = qb.propertiesToSelectString(entityType, qpi, properties, selectString, addAs, fromAs, isExpand)
	selectString = qb.subEntitiesToSelectString(qpi, selectString)

	return selectString
}

func (qb *QueryBuilder) getProperties(et entities.Entity, qo *odata.QueryOptions, addID bool) []string {
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

	return properties
}

func (qb *QueryBuilder) subEntitiesToSelectString(qpi *QueryParseInfo, selectString string) string {
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

func (qb *QueryBuilder) propertiesToSelectString(entityType entities.EntityType, qpi *QueryParseInfo, properties []string, selectString string, addAs, fromAs, isExpand bool) string {
	for _, p := range properties {
		toAdd := ""
		if len(selectString) > 0 {
			toAdd += ", "
		}

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

		if isExpand {
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
				"LIMIT %v OFFSET %s) AS %s on true ",
				joinString,
				joinType,
				qb.getSelect(e2, nqo, qpi, true, true, false, true, ""),
				qb.tables[et2],
				join,
				qb.getFilterQueryString(et2, nqo, filterPrefix),
				qb.getOrderBy(et2, nqo, false),
				qb.getLimit(nqo, 0),
				qb.getOffset(nqo),
				qb.addAsPrefix(qpi, tableMappings[et2]))
		}
	}

	if qpi != nil && len(qpi.SubEntities) > 0 {
		for _, subQPI := range qpi.SubEntities {
			// set innerjoinExpand to false to left join the expand
			innerjoinExpand := false

			qo := &odata.QueryOptions{}
			if subQPI.ExpandItem != nil {
				qo = odata.ExpandItemToQueryOptions(subQPI.ExpandItem)

				// Inner join if there was an outer filter found and added to the expandItem
				if subQPI.ExpandItem.HasOuterFilter {
					innerjoinExpand = true
				}
			}

			// if first select value is nil means the expand is not requested by the user so
			// supply qo to createJoin as a non Expand
			generatedExpand := isExpandGenerated(qo.Select)
			if !generatedExpand {
				//qo.Select = &godata.GoDataSelectQuery{}
			}

			joinString = qb.createJoin(subQPI.Parent.Entity, subQPI.Entity, nil, true, innerjoinExpand, qo, subQPI, joinString)
		}
	}

	return joinString
}

func (qb *QueryBuilder) createSelectByRelationString(e1 entities.Entity, e2 entities.Entity, id interface{}, qpi *QueryParseInfo) string {
	et2 := e2.GetEntityType()

	nqo := &odata.QueryOptions{}
	nqo.Select = &godata.GoDataSelectQuery{SelectItems: []*godata.SelectItem{{Segments: []*godata.Token{{Value: "id"}}}}}
	if id != nil {
		var err error
		nqo.Filter, err = godata.ParseFilterString(fmt.Sprintf("id eq %v", id))
		if err != nil {
			fmt.Printf("\n\n ERROR %v \n\n", err)
		}
	}

	join := getJoin(qb.tables, et2, e1.GetEntityType(), "")
	lowerJoin := strings.ToLower(join)
	filterPrefix := "WHERE"
	if strings.Contains(lowerJoin, "where") {
		filterPrefix = "AND"
	}

	filter := qb.getFilterQueryString(et2, nqo, filterPrefix)

	joinString := fmt.Sprintf(
		"(SELECT %s "+
			"FROM %s %s "+
			"%s) "+
			"IS NOT NULL",
		qb.getSelect(e2, nqo, nil, true, true, false, false, ""),
		qb.tables[et2],
		join,
		filter)

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
	// Only select last Location since we only have 1 encoding
	if et == entities.EntityTypeLocation {
		top := godata.GoDataTopQuery(1)
		if qo == nil {
			nQo := &odata.QueryOptions{}
			qo = nQo
		}

		qo.Top = &top
	}

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
	if pn == nil || pn.Token == nil {
		return ""
	}

	return qb.filterToString(pn, et, ignoreSelectAs)
}

func (qb *QueryBuilder) filterToString(pn *godata.ParseNode, et entities.EntityType, ignoreSelectAs bool) string {
	if convertFunction, ok := filterToStringMap[pn.Token.Type]; ok {
		return convertFunction(*qb, pn, et, ignoreSelectAs)
	}

	return ""
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

		// stil navigation node in children move further
		if pn.Children[0].Token.Type == godata.FilterTokenNav {
			*currentExpand = append([]string{pn.Children[1].Token.Value}, *currentExpand...)
			qb.sortFilter(qo, pn.Children[0], parentNode, startNavNode, currentExpand)
			return
		}

		// prepend entity and field to keep track of path
		*currentExpand = append([]string{pn.Children[1].Token.Value}, *currentExpand...)
		*currentExpand = append([]string{pn.Children[0].Token.Value}, *currentExpand...)

		// if [0] is not an entity do nothing
		_, err := entities.EntityFromString(pn.Children[0].Token.Value)
		if err != nil {
			return
		}

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
		qb.addExpand(qo, afterCouplingNode, cp)

		// remove node from filter since it sits inside the expand now
		*afterCouplingNode = godata.ParseNode{}
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

func iterateParseNodeChilds(pn *godata.ParseNode, str string) string {
	if pn == nil || pn.Token == nil {
		return str
	}

	str = fmt.Sprintf("%v %v", str, pn.Token.Value)

	for _, p := range pn.Children {
		if p != nil && len(pn.Children) > 0 {
			str = iterateParseNodeChilds(p, str)
		}
	}

	return str
}

func (qb *QueryBuilder) addExpand(qo *odata.QueryOptions, afterCouplingNode *godata.ParseNode, cp []string) {
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
		newExpandItem.IsGenerated = true
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
			removeUnusedNodeFromParent(i, pn, current)
		}
	}
}

func removeUnusedNodeFromParent(i int, parent *godata.ParseNode, current *godata.ParseNode) {
	if len(current.Children) == 2 {
		// If 2 childs are found with empty token remove the node from it's parent
		if current.Children[0].Token == nil && current.Children[1].Token == nil {
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			return
		}

		// If child is empty, push back child 1
		if current.Children[0].Token == nil && current.Children[1].Token != nil {
			current = current.Children[1]
		} else if current.Children[0].Token != nil && current.Children[1].Token == nil {
			current = current.Children[0]
		}
	}

	if len(current.Children) == 1 && current.Children[0].Token == nil {
		parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
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
		copyCouplingNode := &godata.ParseNode{}
		cpn := findFirstCouplingParseNode(copyParseNode)
		*copyCouplingNode = *cpn // make copy else filter will be added to root
		copyCouplingNode.Children = make([]*godata.ParseNode, 2)
		copyCouplingNode.Children[0] = &temp
		copyCouplingNode.Children[1] = copyParseNode

		ei.Filter = &godata.GoDataFilterQuery{Tree: copyCouplingNode}
	}

	ei.HasOuterFilter = true
}

func findParseNodeAfterCoupling(pn *godata.ParseNode) *godata.ParseNode {
	if pn.Parent == nil {
		return pn
	}

	// if parent node is and/or return parseNode whihc is the node after coupling
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
	if queryOptions.Count == nil {
		return ""
	}

	if logger.Logger.Level == log.DebugLevel {
		defer gostLog.DebugWithElapsedTime(logger, time.Now(), "constructing count query")
	}

	queryString, _ := qb.getQueryString(e1, e2, id, queryOptions, true)

	return queryString
}

// CreateQuery creates a new count query based on given input
//   e1: entity to get
//   e2: from entity
//   id: e2 == nil: where e1.id = ... | e2 != nil: where e2.id = ...
// example: Datastreams(1)/Thing = CreateQuery(&entities.Thing, &entities.Datastream, 1, nil)
func (qb *QueryBuilder) CreateQuery(e1 entities.Entity, e2 entities.Entity, id interface{}, queryOptions *odata.QueryOptions) (string, *QueryParseInfo) {
	if logger.Logger.Level == log.DebugLevel {
		defer gostLog.DebugWithElapsedTime(logger, time.Now(), "constructing select query")
	}

	queryString, qpi := qb.getQueryString(e1, e2, id, queryOptions, false)

	return queryString, qpi
}

func (qb *QueryBuilder) getQueryString(e1 entities.Entity, e2 entities.Entity, id interface{}, queryOptions *odata.QueryOptions, isCount bool) (string, *QueryParseInfo) {
	var qo *odata.QueryOptions
	qo = nil

	if queryOptions != nil {
		qo = &odata.QueryOptions{}
		*qo = *queryOptions
		qb.sortQueryOptions(qo)
	}

	et1 := e1.GetEntityType()

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

	firstSelect := fmt.Sprintf("COUNT(DISTINCT %s)", qb.addAsPrefix(qpi, fmt.Sprintf("%s.%s", tableMappings[et1], asMappings[et1][idField])))
	if !isCount {
		firstSelect = qb.getSelect(e1, qo, qpi, true, true, true, false, "")
	}

	queryString := fmt.Sprintf("SELECT %s FROM (SELECT %s FROM %s",
		firstSelect,
		qb.getSelect(e1, qo, nil, true, true, false, true, ""),
		qb.tables[et1],
	)

	where := ""

	if id != nil && e2 == nil {
		where = fmt.Sprintf("%s WHERE %s = %v", where, selectMappings[et1][idField], id)
	}

	if qo != nil && qo.Filter != nil {
		if id != nil {
			if e2 == nil {
				where = fmt.Sprintf("%s AND %s", where, qb.getFilterQueryString(et1, qo, ""))
			} else {
				where = fmt.Sprintf("%s %s", where, qb.getFilterQueryString(et1, qo, "WHERE"))
			}
		} else {
			where = fmt.Sprintf("%s %s", where, qb.getFilterQueryString(et1, qo, "WHERE"))
		}
	}

	queryString = fmt.Sprintf("%s %s", queryString, strings.TrimSpace(where))

	// get entity by other entity
	if e2 != nil {
		selectBy := qb.createSelectByRelationString(e1, e2, id, qpi)
		prefix := "WHERE"
		if where != "" {
			prefix = "AND"
		}

		queryString = fmt.Sprintf("%s %s %s", queryString, prefix, selectBy)
	}

	if isCount {
		queryString = fmt.Sprintf("%s ORDER BY %s) AS %s",
			queryString,
			qb.getOrderBy(et1, qo, true),
			qb.addAsPrefix(qpi, tableMappings[et1]),
		)
	} else {
		limit := ""
		if qo != nil && qo.Top != nil && int(*qo.Top) != -1 {
			limit = fmt.Sprintf("LIMIT %v", qb.getLimit(qo, 1))
		}

		queryString = fmt.Sprintf("%s ORDER BY %s %s OFFSET %s) AS %s",
			queryString,
			qb.getOrderBy(et1, qo, true),
			limit,
			qb.getOffset(qo),
			qb.addAsPrefix(qpi, tableMappings[et1]),
		)
	}

	queryString = fmt.Sprintf("%s %s",
		queryString,
		qb.createJoin(e1, e2, id, false, false, qo, qpi, ""),
	)

	return queryString, qpi
}
