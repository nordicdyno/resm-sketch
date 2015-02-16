package store

import (
    "errors"
    "strconv"
)


type ResourcesList []string

// ResourcesInfo struct for report
type ResourcesInfo struct {
	Allocated   map[string]ResourcesList
	Deallocated ResourcesList
}

type ResourceAllocater interface {
	Allocate(user string) (id string, err error)
	AddResource(id string) error
	Deallocate(id string) error
	List() (*ResourcesInfo, error)
	ListByUser(user string) (list ResourcesList, err error)

	Reset() error
}

var (
    ErrResourcesIsOver = errors.New("Resources are over")
    ErrResourcesNotFound = errors.New("Resource not found")
)

func GenResourcesIds(limit int) []string {
    ids := make([]string, limit)
    for i := 0; i< limit; i++ {
        ids[i] = "r" + strconv.Itoa(i)
    }
    return ids
}