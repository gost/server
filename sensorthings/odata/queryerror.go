package odata

import (
	"fmt"
	gostErrors "github.com/geodan/gost/errors"
)

// QueryErrorMessage describes what went wrong with the incoming request
type QueryErrorMessage string

// List of QueryErrorMessages, %s in the string will be replaced by a given value by running it trough FormatQueryError
const (
	QueryTopInvalid          QueryErrorMessage = "The value %s for $top is invalid, please provide a non-negative integer"
	QuerySkipInvalid         QueryErrorMessage = "The value %s for $skip is invalid, please provide a non-negative integer"
	QueryOrderByInvalid      QueryErrorMessage = "The value %s for $orderby is invalid, please use the following format $orderby=\"propertyname\" \"asc/desc\""
	QueryCountInvalid        QueryErrorMessage = "The value %s for $count is invalid, available options: \"true\" or \"false\" "
	QueryResultFormatInvalid QueryErrorMessage = "The value %s for $resultFormat is invalid, available options: dataArray"
	QueryUnknown             QueryErrorMessage = "The query parameter %s is not supported"
	QueryNotAvailable        QueryErrorMessage = "Query %s is not available on endpoint %s"
	QueryExpandAvailable     QueryErrorMessage = "Expand %s is not available on endpoint %s"
)

// CreateQueryError formats a query error, adding a value into the defined message
// for example the QueryErrorMessage "The value %v for $top is invalid, please provide a non-negative integer"
// will be formatted into "The value -1 for $top is invalid, please provide a non-negative integer"
func CreateQueryError(msg QueryErrorMessage, status int, value ...string) error {
	m := fmt.Sprintf("%s", msg)
	return gostErrors.NewErrorWithStatusCode(fmt.Errorf(m, value), status)
}
