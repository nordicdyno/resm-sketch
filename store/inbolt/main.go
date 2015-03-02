package inbolt

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"

	"github.com/nordicdyno/resm-sketch/store"
)

const BucketName = "Pool"

type Storage struct {
	db *bolt.DB
	// resources []*Resource
	// cached value
	// left int
	// sync.Mutex
}

func NewStorage(file string, limit int) (*Storage, error) {
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		name := []byte(BucketName)
		err := tx.DeleteBucket(name)
		if err != nil && err != bolt.ErrBucketNotFound {
			return fmt.Errorf("delete bucket %s: %s", string(name), err)
		}

		b, err := tx.CreateBucket(name)
		if err != nil {
			return fmt.Errorf("create bucket %s: %s", string(name), err)
		}

		for _, id := range store.GenResourcesIds(limit) {
			buf := &bytes.Buffer{}
			r := store.Resource{
				Id:   id,
				User: "",
			}
			err := gob.NewEncoder(buf).Encode(r)
			if err != nil {
				return fmt.Errorf("encode data %v: %s", r, err)
			}
			err = b.Put([]byte(id), buf.Bytes())
			if err != nil {
				return fmt.Errorf("put data: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
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
	deallocated := make(store.ResourcesList, 0, 0)
	allocated := make([]store.Resource, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			res := store.Resource{}
			dec := gob.NewDecoder(bytes.NewBuffer(v))

			if err := dec.Decode(&res); err != nil {
				return fmt.Errorf("decode data %v: %s", v, err)
			}

			//log.Printf("key=%s, value=%v\n", k, res)
			if res.User == "" {
				deallocated = append(deallocated, res.Id)
				continue
			}
			allocated = append(allocated, res)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	r := &store.ResourcesInfo{
		Allocated:   allocated,
		Deallocated: deallocated,
	}
	// spew.Dump(r)
	return r, nil
}

func (s *Storage) Allocate(user string) (string, error) {
	id := ""
	err := s.db.Update(func(tx *bolt.Tx) error {
		name := []byte(BucketName)
		b := tx.Bucket(name)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			res := store.Resource{}
			dec := gob.NewDecoder(bytes.NewBuffer(v))

			if err := dec.Decode(&res); err != nil {
				return fmt.Errorf("decode data %v: %s", v, err)
			}

			if res.User != "" {
				continue
			}

			id = res.Id
			buf := &bytes.Buffer{}
			res.User = user

			err := gob.NewEncoder(buf).Encode(res)
			if err != nil {
				return fmt.Errorf("encode data %v: %s", res, err)
			}
			err = b.Put([]byte(res.Id), buf.Bytes())
			if err != nil {
				return fmt.Errorf("put data: %s", err)
			}
			//log.Println("put key", res.Id)
			break
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if id == "" {
		return "", store.ErrResourcesIsOver
	}
	return id, nil
}

func (s *Storage) AddResource(id string) error {
	return nil
}

func (s *Storage) Deallocate(id string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		name := []byte(BucketName)
		b := tx.Bucket(name)

		v := b.Get([]byte(id))
		if v == nil {
			return store.ErrResourcesNotFound
		}

		res := store.Resource{}

		dec := gob.NewDecoder(bytes.NewBuffer(v))
		if err := dec.Decode(&res); err != nil {
			return fmt.Errorf("fail decode data %v: %s", v, err)
		}
		if res.User == "" {
			return store.ErrResourcesIsFree
		}

		buf := &bytes.Buffer{}
		res.User = ""
		err := gob.NewEncoder(buf).Encode(res)
		if err != nil {
			return fmt.Errorf("encode data %v: %s", res, err)
		}

		err = b.Put([]byte(res.Id), buf.Bytes())
		if err != nil {
			return fmt.Errorf("put data: %s", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

var _ = spew.Config
var _ = log.Prefix

func (s *Storage) Reset() error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		name := []byte(BucketName)
		b := tx.Bucket(name)

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			res := store.Resource{}
			dec := gob.NewDecoder(bytes.NewBuffer(v))

			if err := dec.Decode(&res); err != nil {
				return fmt.Errorf("decode data %v: %s", v, err)
			}

			buf := &bytes.Buffer{}
			res.User = ""

			err := gob.NewEncoder(buf).Encode(res)
			if err != nil {
				return fmt.Errorf("encode data %v: %s", res, err)
			}
			err = b.Put([]byte(res.Id), buf.Bytes())
			if err != nil {
				return fmt.Errorf("put data: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Dump() {
	spew.Dump(s)
}
