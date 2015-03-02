package inmemory

import (
	"log"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/nordicdyno/resm-sketch/store"
)

type Storage struct {
	resources []*store.Resource
	// cached value
	left int
	sync.Mutex
}

func NewStorage(limit int) (*Storage, error) {
	resources := make([]*store.Resource, limit)
	for i, id := range store.GenResourcesIds(limit) {
		resources[i] = &store.Resource{id, ""}
	}
	return &Storage{
		resources: resources,
		left:      limit,
	}, nil
}

func (s *Storage) ListByUser(user string) (store.ResourcesList, error) {
	info, err := s.List()
	if err != nil {
		return nil, err
	}

	idsList := make([]string, 0)
	for _, pair := range info.Allocated {
		if user != pair.User {
			continue
		}
		idsList = append(idsList, pair.Id)
	}
	return store.ResourcesList(idsList), nil
}

func (s *Storage) List() (*store.ResourcesInfo, error) {
	s.Lock()
	defer s.Unlock()

	deallocated := make(store.ResourcesList, 0, s.left)
	allocated := make([]store.Resource, 0)
	for _, res := range s.resources {
		if res.User == "" {
			deallocated = append(deallocated, res.Id)
			continue
		}
		allocated = append(allocated, *res)
	}

	r := &store.ResourcesInfo{
		Allocated:   allocated,
		Deallocated: deallocated,
	}
	return r, nil
}

func (s *Storage) Allocate(user string) (string, error) {
	s.Lock()
	defer s.Unlock()
	if s.left == 0 {
		return "", store.ErrResourcesIsOver
	}

	var item string
	for _, res := range s.resources {
		if res.User != "" {
			continue
		}
		res.User = user
		item = res.Id
		break
	}
	s.left -= 1
	return item, nil
}

func (s *Storage) AddResource(id string) error {
	s.Lock()
	s.resources = append(s.resources, &store.Resource{id, ""})
	s.left += 1
	s.Unlock()
	return nil
}

func (s *Storage) Deallocate(id string) error {
	s.Lock()
	defer s.Unlock()

	var found bool
	var res *store.Resource
	for _, res = range s.resources {
		if res.Id != id {
			continue
		}
		found = true
		break
	}
	if !found {
		return store.ErrResourcesNotFound
	}
	if res.User == "" {
		return store.ErrResourcesIsFree
	}

	res.User = ""
	s.left += 1
	return nil
}

var _ = spew.Config
var _ = log.Prefix

func (s *Storage) Reset() error {
	s.Lock()
	defer s.Unlock()

	for _, res := range s.resources {
		res.User = ""
	}
	s.left = len(s.resources)
	return nil
}

func (s *Storage) Dump() {
	s.Lock()
	spew.Dump(s)
	s.Unlock()
}
