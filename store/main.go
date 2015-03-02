package store

import (
	"errors"
	"strconv"
)

type ResourcesList []string

type Resource struct {
	Id   string
	User string
}

// ResourcesInfo struct for report
type ResourcesInfo struct {
	Allocated   []Resource
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
	ErrResourcesIsOver   = errors.New("Resources are over")
	ErrResourcesNotFound = errors.New("Resource not found")
	ErrResourcesIsFree   = errors.New("Resource already free")
)

func GenResourcesIds(limit int) []string {
	ids := make([]string, limit)
	for i := 0; i < limit; i++ {
		ids[i] = "r" + strconv.Itoa(i)
	}
	return ids
}
