package odata

var supportedExpandParamsMap = map[string][]string{}
var supportedSelectParamsMap = map[string][]string{}

// Init sets the supportedExpandParamsMap and supportedSelectParamsMap, the IsValid functions under some queries need these
// maps to check if the query is supported map[string][]string{} = map[endpoint name][]supported strings array
func Init(expandParams, selectParams map[string][]string) {
	supportedExpandParamsMap = expandParams
	supportedSelectParamsMap = selectParams
}
