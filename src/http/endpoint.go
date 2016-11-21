package http

import (
	"github.com/geodan/gost/src/sensorthings/models"
	"strings"
)

// HttpEndpoints is a slice of *HttpEndpoint's, it implements some functions to be able to
// sort the slice
type HttpEndpoints []*HttpEndpoint

// HttpEndpoint combines a SensorThings endpoint and operation in preparation to add
// it to the router
type HttpEndpoint struct {
	Endpoint  models.Endpoint
	Operation models.EndpointOperation
}

// Len returns the number of elements in the collection
func (a HttpEndpoints) Len() int { return len(a) }

// Swap swaps the elements
func (a HttpEndpoints) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less holds the custom sorting logic
func (a HttpEndpoints) Less(i, j int) bool {
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
		} else {
			return strings.Count(a[i].Operation.Path, "{") < strings.Count(a[j].Operation.Path, "{")
		}
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

func isDynamic(url string) bool {
	return strings.Contains(url, "{") && strings.Contains(url, "}")
}
