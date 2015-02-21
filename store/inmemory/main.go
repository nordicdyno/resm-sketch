package inmemory

import (
	"log"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/nordicdyno/resm-sketch/store"
)

type Resource struct {
	id      string
	free    bool
	ownedBy string
}

type Storage struct {
	resources []*Resource
	// cached value
	left int
	sync.Mutex
}

func NewStorage(limit int) (*Storage, error) {
	resources := make([]*Resource, limit)
	for i, id := range store.GenResourcesIds(limit) {
		resources[i] = &Resource{id, true, ""}
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

	list, ok := info.Allocated[user]
	if !ok {
		list = make(store.ResourcesList, 0)
	}
	return list, nil
}

func (s *Storage) List() (*store.ResourcesInfo, error) {
	s.Lock()
	defer s.Unlock()
	//spew.Dump(s)

	//log.Println("allocate list with size of", s.left)
	deallocated := make(store.ResourcesList, 0, s.left)
	allocatedByUser := make(map[string]store.ResourcesList)
	for _, res := range s.resources {
		if res.free {
			deallocated = append(deallocated, res.id)
		} else {
			list, ok := allocatedByUser[res.ownedBy]
			if !ok {
				list = make(store.ResourcesList, 0, 1)
			}
			list = append(list, res.id)
			allocatedByUser[res.ownedBy] = list
		}
	}
	r := &store.ResourcesInfo{
		Allocated:   allocatedByUser,
		Deallocated: deallocated,
	}
	// spew.Dump(r)
	return r, nil
}

func (s *Storage) Allocate(user string) (string, error) {
	s.Lock()
	defer s.Unlock()
	if s.left == 0 {
		return "", store.ErrResourcesIsOver
	}

	//spew.Dump(s)

	var item string
	for _, res := range s.resources {
		if !res.free {
			continue
		}
		//log.Println("ALLOCATE", res)
		res.ownedBy = user
		res.free = false
		item = res.id
		break
	}
	s.left -= 1
	//log.Println("RETURN", item)
	return item, nil
}

func (s *Storage) AddResource(id string) error {
	s.Lock()
	s.resources = append(s.resources, &Resource{id, true, ""})
	s.left += 1
	s.Unlock()
	return nil
}

func (s *Storage) Deallocate(id string) error {
	s.Lock()
	defer s.Unlock()

	var found bool
	var res *Resource
	for _, res = range s.resources {
		if res.id != id {
			continue
		}
		found = true
		break
	}
	if !found {
		return store.ErrResourcesNotFound
	}
	if res.free {
		return store.ErrResourcesNotFound
	}

	res.free = true
	s.left += 1
	return nil
}

var _ = spew.Config
var _ = log.Prefix

func (s *Storage) Reset() error {
	s.Lock()
	defer s.Unlock()

	for _, res := range s.resources {
		res.free = true
	}
	s.left = len(s.resources)
	return nil
}

func (s *Storage) Dump() {
	s.Lock()
	spew.Dump(s)
	s.Unlock()
}