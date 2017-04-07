package odata

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Operator are the supported ODATA operators
type Operator string

// List of all ODATA Operators.
const (
	//logical
	And Operator = "and"
	Or  Operator = "or"
	Not Operator = "not"

	//comparison
	Equals              Operator = "eq"
	NotEquals           Operator = "ne"
	GreaterThan         Operator = "gt"
	GreaterThanOrEquals Operator = "ge"
	LessThan            Operator = "lt"
	LessThanOrEquals    Operator = "le"
	Like                Operator = "like"
	IsNull              Operator = "is null"

	// arithmetic
	Addition       Operator = "add"
	Subtraction    Operator = "sub"
	Multiplication Operator = "mul"
	Division       Operator = "div"
	Modulo         Operator = "mod"
)

// IsLogical returns whether a defined operation is a logical operator or not.
func (o Operator) IsLogical() bool {
	if o == And || o == Or {
		return true
	}

	return false
}

// IsUnary returns whether a defined operation is unary or binary.  Will return true
// if the operation only supports a subject with no value.
func (o Operator) IsUnary() bool {
	if o == IsNull {
		return true
	}

	return false
}

// ToString representation of the ComparisonOperator
func (o Operator) ToString() string {
	return fmt.Sprintf("%s", o)
}

// ODATAOperators map, using: ODATAOperators["eq"]
var ODATAOperators = map[string]Operator{
	// logic
	And.ToString(): And,
	Or.ToString():  Or,
	Not.ToString(): Not,

	// comparison
	Equals.ToString():              Equals,
	NotEquals.ToString():           NotEquals,
	GreaterThan.ToString():         GreaterThan,
	GreaterThanOrEquals.ToString(): GreaterThanOrEquals,
	LessThan.ToString():            LessThan,
	LessThanOrEquals.ToString():    LessThanOrEquals,
	Like.ToString():                Like,
	IsNull.ToString():              IsNull,

	// arithmetic
	Addition.ToString():       LessThanOrEquals,
	Subtraction.ToString():    Subtraction,
	Multiplication.ToString(): Multiplication,
	Division.ToString():       Division,
	Modulo.ToString():         Modulo,
}

// Predicate is the basic model construct of the odata expression
type Predicate struct {
	Right    interface{}
	Left     interface{}
	Operator Operator
}

// Split separates the left and right predicates when containing a *Predicate
// and returns a slice of predicates without a predicate inside left or right,
// the predicates can be coupled with the logical operators inside the OdataOperator slice
//
//  predicate := ParseODATAFilter("id ge 10 and id lt 100 or user eq 'tim' and status 'active'")
//  ps, ops := predicate.Split()
//  for i, p := range ps {
//  	log.Printf("%v <- %v -> %v", p.Left, p.Operator, p.Right)
//  	if len(ops)-1 >= i {
//  		log.Printf("--%v--", ops[i])
//  	}
//  }
//
//  output = id <- ge -> 10 --and-- id <- lt -> 100 --or-- user <- eq -> 'tim' --and-- status <- eq -> 'active'
func (p *Predicate) Split() ([]*Predicate, []Operator) {
	var pr []*Predicate
	predicates := &pr

	var ops []Operator
	operators := &ops

	if p.Operator.IsLogical() {
		internalSplit(p, predicates, operators)
	} else {
		*predicates = append(*predicates, p)
	}

	return pr, ops
}

func internalSplit(p *Predicate, result *[]*Predicate, ops *[]Operator) {
	if p.Operator.IsLogical() {
		*ops = append(*ops, p.Operator)
		internalSplit(p.Left.(*Predicate), result, ops)
		internalSplit(p.Right.(*Predicate), result, ops)
	} else {
		*result = append(*result, p)
	}
}

// ParseODATAFilter parses a filter string into a predicate, simple implementation for now
// ToDo: handle separators ( ) -> also in split function
// ToDo: get the 'not' operator working
// ToDo: check for valid operators -> return error
// ToDo: handle arithmetic operators
// ToDo: handle functions (string, date, math, geospatial, spatialRelationship)
func ParseODATAFilter(filterStr string) (*Predicate, error) {
	if len(filterStr) == 0 {
		return nil, errorInvalidFilter
	}

	filter := strings.TrimSpace(filterStr)
	if len(filter) == 0 {
		return nil, errorInvalidFilter
	}

	fragment, err := parseFragment(filter)
	return fragment, err
}

var odataRegex = map[string]string{
	"regexParenthesis": "^([(](.*)[)])$",
	"regexAndor":       "^(.*?) (or|and)+ (.*)$",
	"regexOp":          "(.*\\w*) (eq|gt|lt|ge|le|ne) (datetimeoffset'(.*)'|(.*))",
	"regexStartsWith":  "^startswith[(](.*),'(.*)'[)]",
	"regexEndsWith":    "^endswith[(](.*),'(.*)'[)]",
	"regexContains":    "^contains[(](.*),'(.*)'[)]",
}

var errorInvalidFilter = errors.New("Invalid filter")

func parseFragment(filter string) (*Predicate, error) {
	var err error
	found := false
	predicate := &Predicate{}

	for k, regex := range odataRegex {
		if found {
			break
		}

		r, _ := regexp.Compile(regex)
		match := r.FindStringSubmatch(filter)
		if len(match) > 0 {
			if k == "regexParenthesis" {
				if len(match) > 2 {
					if strings.Index(match[2], ")") < strings.Index(match[2], "(") {
						continue
					}

					if predicate, err = parseFragment(match[2]); err != nil {
						return nil, errorInvalidFilter
					}
				}
			} else if k == "regexAndor" {
				if len(match) < 4 {
					return nil, errorInvalidFilter
				}

				predicate = &Predicate{
					Operator: ODATAOperators[match[2]],
				}

				if predicate.Left, err = parseFragment(match[1]); err != nil {
					return nil, errorInvalidFilter
				}

				if predicate.Right, err = parseFragment(match[3]); err != nil {
					return nil, errorInvalidFilter
				}
			} else if k == "regexOp" {
				if len(match) < 4 {
					return nil, errorInvalidFilter
				}

				var val interface{}

				// if not string value
				if strings.Index(match[3], "'") == -1 {
					if val, err = strconv.ParseFloat(match[3], 64); err != nil {
						val = fmt.Sprintf("'%s'", match[3])
					}
				} else {
					val = match[3]
				}

				predicate = &Predicate{
					Left:     match[1],
					Operator: ODATAOperators[match[2]],
					Right:    val,
				}

				/*if(predicate.Value.indexOf && predicate.Value.indexOf("datetimeoffset") == 0)
				{
					var m = predicate.Value.match(/^datetimeoffset'(.*)'$/);
					if( m && m.length > 1) {
						obj.value = new Date(m[1]);
					}
				}*/
			} else if k == "regexStartsWith" || k == "regexEndsWith" || k == "regexContains" {
				predicate, err = buildLike(match, k)
				if err != nil {
					return nil, err
				}
			}

			found = true
		}
	}

	return predicate, nil
}

func buildLike(match []string, key string) (*Predicate, error) {
	if len(match) < 3 {
		return nil, errorInvalidFilter
	}

	var right string
	if key == "startsWith" {
		right = match[2] + "*"
	} else {
		if key == "endsWith" {
			right = "*" + match[2]
		} else {
			right = "*" + match[2] + "*"
		}
	}

	p := &Predicate{
		Left:     match[1],
		Operator: Like,
		Right:    right,
	}

	return p, nil
}

// StringFunction are the supporter ODATA string functions
type StringFunction string

// List of all ODATA String Functions
const (
	OSSubstringOf StringFunction = "substringof"
	OSEndsWith    StringFunction = "endswith"
	OSStartsWith  StringFunction = "startswith"
	OSLength      StringFunction = "length"
	OSIndexOf     StringFunction = "indexof"
	OSSubstring   StringFunction = "substring"
	OSToLower     StringFunction = "tolower"
	OSToUpper     StringFunction = "toupper"
	OSTrim        StringFunction = "trim"
	OSConcat      StringFunction = "concat"
)

// ToString representation of a StringFunction
func (s StringFunction) ToString() string {
	return fmt.Sprintf("%s", s)
}

// StringFunctions map, using: StringFunctions["substring"]
var StringFunctions = map[string]StringFunction{
	OSSubstringOf.ToString(): OSSubstringOf,
	OSEndsWith.ToString():    OSEndsWith,
	OSStartsWith.ToString():  OSStartsWith,
	OSLength.ToString():      OSLength,
	OSIndexOf.ToString():     OSIndexOf,
	OSSubstring.ToString():   OSSubstring,
	OSToLower.ToString():     OSToLower,
	OSToUpper.ToString():     OSToUpper,
	OSTrim.ToString():        OSTrim,
	OSConcat.ToString():      OSConcat,
}

// DateFunction are the supporter ODATA date functions
type DateFunction string

// List of all ODATA date functions
const (
	ODYear               DateFunction = "year"
	ODMonth              DateFunction = "month"
	ODDay                DateFunction = "day"
	ODHour               DateFunction = "hour"
	ODMinute             DateFunction = "minute"
	ODSecond             DateFunction = "second"
	ODFractionalSeconds  DateFunction = "fractionalseconds"
	ODDate               DateFunction = "date"
	ODTime               DateFunction = "time"
	ODTotalOffsetMinutes DateFunction = "totaloffsetminutes"
	ODNow                DateFunction = "now"
	ODMinDateTime        DateFunction = "mindatetime"
	ODMaxDateTime        DateFunction = "maxdatetime"
)

// ToString representation of a DateFunction
func (d DateFunction) ToString() string {
	return fmt.Sprintf("%s", d)
}

// DateFunctions map, using: DateFunctions["now"]
var DateFunctions = map[string]DateFunction{
	ODYear.ToString():               ODYear,
	ODMonth.ToString():              ODMonth,
	ODDay.ToString():                ODDay,
	ODHour.ToString():               ODHour,
	ODMinute.ToString():             ODMinute,
	ODSecond.ToString():             ODSecond,
	ODFractionalSeconds.ToString():  ODFractionalSeconds,
	ODDate.ToString():               ODDate,
	ODTime.ToString():               ODTime,
	ODTotalOffsetMinutes.ToString(): ODTotalOffsetMinutes,
	ODNow.ToString():                ODNow,
	ODMinDateTime.ToString():        ODMinDateTime,
	ODMaxDateTime.ToString():        ODMaxDateTime,
}

// MathFunction are the supporter ODATA math functions
type MathFunction string

// List of all ODATA math functions
const (
	OMRound   MathFunction = "round"
	OMFloor   MathFunction = "floor"
	OMCeiling MathFunction = "ceiling"
)

// ToString representation of a MathFunction
func (m MathFunction) ToString() string {
	return fmt.Sprintf("%s", m)
}

// MathFunctions map, using: MathFunctions["round"]
var MathFunctions = map[string]MathFunction{
	OMRound.ToString():   OMRound,
	OMFloor.ToString():   OMFloor,
	OMCeiling.ToString(): OMCeiling,
}

// GeospatialFunction are the supporter ODATA geospatial functions
type GeospatialFunction string

// List of all ODATA geospatial functions
const (
	OGSDistance   GeospatialFunction = "geo.distance"
	OGSLength     GeospatialFunction = "geo.length"
	OGSIntersects GeospatialFunction = "geo.intersects"
)

// ToString representation of a GeospatialFunction
func (g GeospatialFunction) ToString() string {
	return fmt.Sprintf("%s", g)
}

// GeospatialFunctions map, using: GeospatialFunctions["distance"]
var GeospatialFunctions = map[string]GeospatialFunction{
	OGSDistance.ToString():   OGSDistance,
	OGSLength.ToString():     OGSLength,
	OGSIntersects.ToString(): OGSIntersects,
}

// SpatialRelationshipFunction are the supporter ODATA spatial relationship functions
type SpatialRelationshipFunction string

// List of all ODATA spatial relationship functions
const (
	OSREquals     SpatialRelationshipFunction = "st_equals"
	OSRDisjoint   SpatialRelationshipFunction = "st_disjoint"
	OSRTouches    SpatialRelationshipFunction = "st_touches"
	OSRWithin     SpatialRelationshipFunction = "st_within"
	OSROverlaps   SpatialRelationshipFunction = "st_overlaps"
	OSRCrosses    SpatialRelationshipFunction = "st_crosses"
	OSRIntersects SpatialRelationshipFunction = "st_intersects"
	OSRContains   SpatialRelationshipFunction = "st_contains"
	OSRRelate     SpatialRelationshipFunction = "st_relate"
)

// ToString representation of a GeospatialFunction
func (g SpatialRelationshipFunction) ToString() string {
	return fmt.Sprintf("%s", g)
}

// SpatialRelationshipFunctions map, using: SpatialRelationshipFunctions["st_equals"]
var SpatialRelationshipFunctions = map[string]SpatialRelationshipFunction{
	OSREquals.ToString():     OSREquals,
	OSRDisjoint.ToString():   OSRDisjoint,
	OSRTouches.ToString():    OSRTouches,
	OSRWithin.ToString():     OSRWithin,
	OSROverlaps.ToString():   OSROverlaps,
	OSRCrosses.ToString():    OSRCrosses,
	OSRIntersects.ToString(): OSRIntersects,
	OSRContains.ToString():   OSRContains,
	OSRRelate.ToString():     OSRRelate,
}
