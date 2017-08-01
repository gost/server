package endpoint

import (
	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"sort"
	"strings"
)

// Endpoint contains all information for creating and handling a main SensorThings endpoint.
// A SensorThings endpoint contains multiple EndpointOperations
// Endpoint can be marshalled to JSON for returning endpoint information requested
// by the user: http://www.sensorup.com/docs/#resource-path
type Endpoint struct {
	Name                  string                     `json:"name"` // Name of the endpoint
	URL                   string                     `json:"url"`  // External URL to the endpoint
	EntityType            entities.EntityType        `json:"-"`
	OutputInfo            bool                       `json:"-"` //Output when BasePathInfo is requested by the user
	Operations            []models.EndpointOperation `json:"-"`
	SupportedExpandParams []string                   `json:"-"`
	SupportedSelectParams []string                   `json:"-"`
}

// GetName returns the endpoint name
func (e *Endpoint) GetName() string {
	return e.Name
}

// ShowOutputInfo returns true if the endpoint should output his info when BasePathInfo is requested
func (e *Endpoint) ShowOutputInfo() bool {
	return e.OutputInfo
}

// GetURL returns the external url
func (e *Endpoint) GetURL() string {
	return e.URL
}

// GetOperations returns all operations for this endpoint such as GET, POST
func (e *Endpoint) GetOperations() []models.EndpointOperation {
	return e.Operations
}

// GetSupportedExpandParams returns which entities can be expanded
func (e *Endpoint) GetSupportedExpandParams() []string {
	return e.SupportedExpandParams
}

// GetSupportedSelectParams returns the supported select parameters for this endpoint
func (e *Endpoint) GetSupportedSelectParams() []string {
	return e.SupportedSelectParams
}

// SortedEndpoints is a slice of *EndpointWrapper's, it implements some functions to be able to
// sort the slice
type SortedEndpoints []*EndpointWrapper

// EndpointWrapper combines a SensorThings endpoint and operation in preparation to add
// it to the router
type EndpointWrapper struct {
	Endpoint  models.Endpoint
	Operation models.EndpointOperation
}

// Len returns the number of elements in the collection
func (a SortedEndpoints) Len() int { return len(a) }

// Swap swaps the elements
func (a SortedEndpoints) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less holds the custom sorting logic
func (a SortedEndpoints) Less(i, j int) bool {
	firstDynamic := isDynamic(a[i].Operation.Path)
	secondDynamic := isDynamic(a[j].Operation.Path)

	if strings.Contains(a[i].Operation.Path, "{c:.*}") && !strings.Contains(a[j].Operation.Path, "{c:.*}") {
		return false
	}

	if !strings.Contains(a[i].Operation.Path, "{c:.*}") && strings.Contains(a[j].Operation.Path, "{c:.*}") {
		return true
	}

	if firstDynamic && !secondDynamic {
		return false
	}

	if !firstDynamic && secondDynamic {
		return true
	}

	if firstDynamic && secondDynamic {
		dynamicI := strings.Count(a[i].Operation.Path, "{")
		dynamicJ := strings.Count(a[j].Operation.Path, "{")
		if dynamicI == dynamicJ {
			return len(a[i].Operation.Path) > len(a[j].Operation.Path)
		}
		return strings.Count(a[i].Operation.Path, "{") < strings.Count(a[j].Operation.Path, "{")
	}

	if len(a[i].Operation.Path) != len(a[j].Operation.Path) {
		return len(a[i].Operation.Path) > len(a[j].Operation.Path)
	}

	if a[i].Operation.OperationType != a[j].Operation.OperationType {
		return a[i].Operation.OperationType != models.HTTPOperationGet
	}

	if a[i].Operation.Path == a[j].Operation.Path {
		panic("Two endpoints can't be same")
	}

	return true
}

// EndpointsToSortedList sorts all the endpoints so they can be added
// to the routes in the right order else requests will be picked up by the wrong handlers
func EndpointsToSortedList(endpoints *map[entities.EntityType]models.Endpoint) SortedEndpoints {
	eps := SortedEndpoints{}
	for _, endpoint := range *endpoints {
		for _, op := range endpoint.GetOperations() {
			e := &EndpointWrapper{Endpoint: endpoint, Operation: op}
			eps = append(eps, e)
		}
	}

	sort.Sort(eps)
	return eps
}

func isDynamic(url string) bool {
	return strings.Contains(url, "{") && strings.Contains(url, "}")
}
