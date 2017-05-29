package odata

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Operator are the supported ODATA operators
type OdataOperator string

// ToString representation of the ComparisonOperator
func (o OdataOperator) ToString() string {
	return fmt.Sprintf("%s", o)
}

// IsLogical returns whether a defined operation is a logical operator or not.
func (o OdataOperator) IsLogical() bool {
	if o == And || o == Or {
		return true
	}

	return false
}

// IsUnary returns whether a defined operation is unary or binary.  Will return true
// if the operation only supports a subject with no value.
func (o OdataOperator) IsUnary() bool {
	return (o == IsNull)
}

// Function are the supported ODATA functions
type Function string

// ToString representation of a function
func (f Function) ToString() string {
	return fmt.Sprintf("%s", f)
}

type OdataFunction struct {
	Function Function
	Args     []interface{}
}

// Predicate is the basic model construct of the odata expression
type Predicate struct {
	Left     interface{}
	Operator OdataOperator
	Right    interface{}
	Function OdataFunction
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
func (p *Predicate) Split() ([]*Predicate, []OdataOperator) {
	var pr []*Predicate
	predicates := &pr

	var ops []OdataOperator
	operators := &ops

	if p.Operator.IsLogical() {
		internalSplit(p, predicates, operators)
	} else {
		*predicates = append(*predicates, p)
	}

	return pr, ops
}

func internalSplit(p *Predicate, result *[]*Predicate, ops *[]OdataOperator) {
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
	"regexParenthesis":   "^([(](.*)[)])$",
	"regexAndor":         "^(.*?) (or|and)+ (.*)$",
	"regexOp":            "(.*\\w*) (eq|gt|lt|ge|le|ne) (datetimeoffset'(.*)'|(.*))",
	"regexStartsWith":    "^startswith[(](.*),'(.*)'[)]",
	"regexEndsWith":      "^endswith[(](.*),'(.*)'[)]",
	"regexContains":      "^contains[(](.*),'(.*)'[)]",
	OSRWithin.ToString(): "^st_within[(](.*),(.*)[)]",
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
				if !strings.Contains(match[3], "'") {
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
			} else if k == "regexStartsWith" || k == "regexEndsWith" || k == "regexContains" || k == OSRWithin.ToString() {
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
	} else if key == "endsWith" {
		right = "*" + match[2]
	} else if key == OSRWithin.ToString() {
		return &Predicate{
			Function: OdataFunction{
				Args:     []interface{}{match[1], match[2]},
				Function: OSRWithin,
			},
		}, nil
	}

	p := &Predicate{
		Left:     match[1],
		Operator: Like,
		Right:    right,
	}

	return p, nil
}

// List of all ODATA Operators.
const (
	//logical
	And OdataOperator = "and"
	Or  OdataOperator = "or"
	Not OdataOperator = "not"

	//comparison
	Equals              OdataOperator = "eq"
	NotEquals           OdataOperator = "ne"
	GreaterThan         OdataOperator = "gt"
	GreaterThanOrEquals OdataOperator = "ge"
	LessThan            OdataOperator = "lt"
	LessThanOrEquals    OdataOperator = "le"
	Like                OdataOperator = "like"
	IsNull              OdataOperator = "is null"

	// arithmetic
	Addition       OdataOperator = "add"
	Subtraction    OdataOperator = "sub"
	Multiplication OdataOperator = "mul"
	Division       OdataOperator = "div"
	Modulo         OdataOperator = "mod"
)

// ODATAOperators map, using: ODATAOperators["eq"]
var ODATAOperators = map[string]OdataOperator{
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

// List of all ODATA Functions
const (
	OSSubstringOf        Function = "substringof"
	OSEndsWith           Function = "endswith"
	OSStartsWith         Function = "startswith"
	OSLength             Function = "length"
	OSIndexOf            Function = "indexof"
	OSSubstring          Function = "substring"
	OSToLower            Function = "tolower"
	OSToUpper            Function = "toupper"
	OSTrim               Function = "trim"
	OSConcat             Function = "concat"
	ODYear               Function = "year"
	ODMonth              Function = "month"
	ODDay                Function = "day"
	ODHour               Function = "hour"
	ODMinute             Function = "minute"
	ODSecond             Function = "second"
	ODFractionalSeconds  Function = "fractionalseconds"
	ODDate               Function = "date"
	ODTime               Function = "time"
	ODTotalOffsetMinutes Function = "totaloffsetminutes"
	ODNow                Function = "now"
	ODMinDateTime        Function = "mindatetime"
	ODMaxDateTime        Function = "maxdatetime"
	OMRound              Function = "round"
	OMFloor              Function = "floor"
	OMCeiling            Function = "ceiling"
	OGSDistance          Function = "geo.distance"
	OGSLength            Function = "geo.length"
	OGSIntersects        Function = "geo.intersects"
	OSREquals            Function = "st_equals"
	OSRDisjoint          Function = "st_disjoint"
	OSRTouches           Function = "st_touches"
	OSRWithin            Function = "st_within"
	OSROverlaps          Function = "st_overlaps"
	OSRCrosses           Function = "st_crosses"
	OSRIntersects        Function = "st_intersects"
	OSRContains          Function = "st_contains"
	OSRRelate            Function = "st_relate"
)

// Functions map, using: Functions["substring"]
var Functions = map[string]Function{
	OSSubstringOf.ToString():        OSSubstringOf,
	OSEndsWith.ToString():           OSEndsWith,
	OSStartsWith.ToString():         OSStartsWith,
	OSLength.ToString():             OSLength,
	OSIndexOf.ToString():            OSIndexOf,
	OSSubstring.ToString():          OSSubstring,
	OSToLower.ToString():            OSToLower,
	OSToUpper.ToString():            OSToUpper,
	OSTrim.ToString():               OSTrim,
	OSConcat.ToString():             OSConcat,
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
	OMRound.ToString():              OMRound,
	OMFloor.ToString():              OMFloor,
	OMCeiling.ToString():            OMCeiling,
	OGSDistance.ToString():          OGSDistance,
	OGSLength.ToString():            OGSLength,
	OGSIntersects.ToString():        OGSIntersects,
	OSREquals.ToString():            OSREquals,
	OSRDisjoint.ToString():          OSRDisjoint,
	OSRTouches.ToString():           OSRTouches,
	OSRWithin.ToString():            OSRWithin,
	OSROverlaps.ToString():          OSROverlaps,
	OSRCrosses.ToString():           OSRCrosses,
	OSRIntersects.ToString():        OSRIntersects,
	OSRContains.ToString():          OSRContains,
	OSRRelate.ToString():            OSRRelate,
}
