package postgis

import (
	"fmt"
	"strings"

	entities "github.com/gost/core"
	"github.com/gost/godata"
)

var funcToStringMap map[string]func(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string

// init is used to work around the initialization loop error (circular reference)
func init() {
	funcToStringMap = map[string]func(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string{
		"contains":           containsToString,
		"substringof":        substringofToString,
		"endswith":           endswithToString,
		"startswith":         startswithToString,
		"length":             lengthToString,
		"indexof":            indexofToString,
		"substring":          substringToString,
		"tolower":            tolowerToString,
		"toupper":            toupperToString,
		"trim":               trimToString,
		"concat":             concatToString,
		"round":              roundToString,
		"floor":              floorToString,
		"ceiling":            ceilingToString,
		"year":               yearToString,
		"month":              monthToString,
		"day":                dayToString,
		"hour":               hourToString,
		"minute":             minuteToString,
		"second":             secondToString,
		"fractionalseconds":  fractionalsecondsToString,
		"date":               dateToString,
		"time":               timeToString,
		"totaloffsetminutes": totaloffsetminutesToString,
		"now":                nowToString,
		"maxdatetime":        maxdatetimeToString,
		"mindatetime":        mindatetimeToString,
		"totalseconds":       totalsecondsToString,
		"geo.length":         geolengthToString,
		"geo.distance":       geodistanceToString,
		"geo.intersects":     stintersectsToString,
		"st_equals":          stequalsToString,
		"st_touches":         sttouchesToString,
		"st_overlaps":        stoverlapsToString,
		"st_crosses":         stcrossesToString,
		"st_contains":        stcontainsToString,
		"st_disjoint":        stdisjointToString,
		"st_relate":          strelateToString,
		"st_within":          stwithinToString,
		"st_intersects":      stintersectsToString,
	}
}

func getLeftAndRightFilter(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType, ignoreLeftSelectAs, ignoreRightSelectAs bool) (string, string) {
	left := qb.createFilter(et, pn.Children[0], ignoreLeftSelectAs)
	right := qb.createFilter(et, pn.Children[1], ignoreRightSelectAs)
	return left, right
}

func containsToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left, right := getLeftAndRightFilter(qb, pn, et, true, true)
	return fmt.Sprintf("%s LIKE %s", qb.createLike(left, LikeContains), qb.createLike(right, LikeContains))
}

func substringofToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left, right := getLeftAndRightFilter(qb, pn, et, true, true)
	return fmt.Sprintf("%s LIKE %s", qb.createLike(right, LikeContains), qb.createLike(left, LikeContains))
}

func endswithToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left, right := getLeftAndRightFilter(qb, pn, et, true, true)
	return fmt.Sprintf("%s LIKE %s", qb.createLike(left, LikeEndsWith), qb.createLike(right, LikeEndsWith))
}

func startswithToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left, right := getLeftAndRightFilter(qb, pn, et, true, true)
	return fmt.Sprintf("%s LIKE %s", qb.createLike(left, LikeStartsWith), qb.createLike(right, LikeStartsWith))
}

func lengthToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("LENGTH(%s)", left)
}

func indexofToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left, right := getLeftAndRightFilter(qb, pn, et, true, true)
	return fmt.Sprintf("STRPOS(%s, %s) -1", left, right)
}

func substringToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left, right := getLeftAndRightFilter(qb, pn, et, true, true)
	right2 := ""
	if len(pn.Children) > 2 {
		right2 = qb.createFilter(et, pn.Children[2], true)
		return fmt.Sprintf("SUBSTRING(%s from (%s + 1) for %s)", left, right, right2)
	}

	return fmt.Sprintf("SUBSTRING(%s from (%s + 1) for LENGTH(%s))", left, right, left)
}

func tolowerToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("LOWER(%s)", left)
}

func toupperToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("UPPER(%s)", left)
}

func trimToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("TRIM(both ' ' from %s)", left)
}

func concatToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left, right := getLeftAndRightFilter(qb, pn, et, true, true)
	return fmt.Sprintf("CONCAT(%s, %s)", left, right)
}

func roundToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("ROUND(CAST(%s as double precision))", strings.Replace(left, "->", "->>", -1))
}

func floorToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("FLOOR(CAST(%s as double precision))", strings.Replace(left, "->", "->>", -1))
}

func ceilingToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("CEILING(CAST(%s as double precision))", strings.Replace(left, "->", "->>", -1))
}

func yearToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createExtractDateQuery(pn, et, "YEAR")
}

func monthToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createExtractDateQuery(pn, et, "MONTH")
}

func dayToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createExtractDateQuery(pn, et, "DAY")
}

func hourToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createExtractDateQuery(pn, et, "HOUR")
}

func minuteToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createExtractDateQuery(pn, et, "MINUTE")
}

func secondToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createExtractDateQuery(pn, et, "SECOND")
}

func fractionalsecondsToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("EXTRACT(MICROSECONDS FROM to_timestamp(%s,'YYYY-MM-DD\"T\"HH24:MI:SS.MS\"Z\"')) / 1000000", left)
}

func dateToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("(%s)::date", left)
}

func timeToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	if strings.Contains(strings.ToLower(left), "time") {
		return fmt.Sprintf("((%s)::timestamp)::time", left)
	}

	return fmt.Sprintf("(%s)::time", left)
}

func totaloffsetminutesToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("EXTRACT(TIMEZONE_MINUTE FROM to_timestamp(%s,'YYYY-MM-DD\"T\"HH24:MI:SS.MS\"Z\"'))", left)
}

func nowToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return fmt.Sprint("to_char(now()::timestamp at time zone 'UTC', 'YYYY-MM-DD\"T\"HH24:MI:SS.MS\"Z\"')")
}

func maxdatetimeToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return fmt.Sprint("'9999-12-31T23:59:59.999Z'")
}

func mindatetimeToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return fmt.Sprint("'0001-01-01T00:00:00.000Z'")
}

func totalsecondsToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	left := qb.createFilter(et, pn.Children[0], true)
	return fmt.Sprintf("SELECT extract(epoch from (%s)::timestamp)", left)
}

func geodistanceToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_DISTANCE(%s, %s)", 2)
}

func geolengthToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_LENGTH(%s)", 1)
}

func stequalsToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_EQUALS(%s, %s)", 2)
}

func sttouchesToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_TOUCHES(%s, %s)", 2)
}

func stoverlapsToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_OVERLAPS(%s, %s)", 2)
}

func stcrossesToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_CROSSES(%s, %s)", 2)
}

func stcontainsToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_CONTAINS(%s, %s)", 2)
}

func stdisjointToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_DISJOINT(%s, %s)", 2)
}

func strelateToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_RELATE(%s, %s, %s)", 3)
}

func stwithinToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_WITHIN(%s, %s)", 2)
}

func stintersectsToString(qb QueryBuilder, pn *godata.ParseNode, et entities.EntityType) string {
	return qb.createSpatialQuery(pn, et, "ST_INTERSECTS(%s, %s)", 2)
}
