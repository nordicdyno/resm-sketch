package inbolt

import (
	"log"
	"time"

	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"

	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nordicdyno/resm-sketch/store"
)

const BucketName = "Pool"

type Resource struct {
	Id      string
	Free    bool
	OwnedBy string
}

type Storage struct {
	db *bolt.DB
	//	resources []*Resource
	// cached value
	//	left int
	//	sync.Mutex
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
		//log.Println("bucket deleted!")

		b, err := tx.CreateBucket(name)
		if err != nil {
			return fmt.Errorf("create bucket %s: %s", string(name), err)
		}

		for _, id := range store.GenResourcesIds(limit) {
			//log.Println("allocate", id)
			buf := &bytes.Buffer{}
			r := Resource{
				Id:      id,
				Free:    true,
				OwnedBy: "",
			}
			err := gob.NewEncoder(buf).Encode(r)
			if err != nil {
				return fmt.Errorf("encode data %v: %s", r, err)
			}
			err = b.Put([]byte(id), buf.Bytes())
			if err != nil {
				return fmt.Errorf("put data: %s", err)
			}
			//log.Println("put key", id)
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
	//log.Println("allocate list with size of", s.left)
	deallocated := make(store.ResourcesList, 0, 0)
	allocated := make([]store.ResourceUserPair, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			res := Resource{}
			dec := gob.NewDecoder(bytes.NewBuffer(v))

			if err := dec.Decode(&res); err != nil {
				return fmt.Errorf("decode data %v: %s", v, err)
			}

			//log.Printf("key=%s, value=%v\n", k, res)
			if res.Free {
				deallocated = append(deallocated, res.Id)
				continue
			}
			allocated = append(allocated, store.ResourceUserPair{
				Id:   res.Id,
				User: res.OwnedBy,
			})
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
			res := Resource{}
			dec := gob.NewDecoder(bytes.NewBuffer(v))

			if err := dec.Decode(&res); err != nil {
				return fmt.Errorf("decode data %v: %s", v, err)
			}

			if !res.Free {
				continue
			}

			id = res.Id
			buf := &bytes.Buffer{}
			res.OwnedBy = user
			res.Free = false

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
			//log.Println("not found", id)
			return store.ErrResourcesNotFound
		}

		res := Resource{}

		dec := gob.NewDecoder(bytes.NewBuffer(v))
		if err := dec.Decode(&res); err != nil {
			return fmt.Errorf("fail decode data %v: %s", v, err)
		}
		//log.Printf("'%v' Get result => %v", id, spew.Sdump(res))
		if res.OwnedBy == "" {
			return store.ErrResourcesNotFound
		}

		buf := &bytes.Buffer{}
		res.OwnedBy = ""
		res.Free = true
		err := gob.NewEncoder(buf).Encode(res)
		if err != nil {
			return fmt.Errorf("encode data %v: %s", res, err)
		}

		err = b.Put([]byte(res.Id), buf.Bytes())
		if err != nil {
			return fmt.Errorf("put data: %s", err)
		}
		//log.Println("put key", res.Id)
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
			res := Resource{}
			dec := gob.NewDecoder(bytes.NewBuffer(v))

			if err := dec.Decode(&res); err != nil {
				return fmt.Errorf("decode data %v: %s", v, err)
			}

			buf := &bytes.Buffer{}
			res.OwnedBy = ""
			res.Free = true

			err := gob.NewEncoder(buf).Encode(res)
			if err != nil {
				return fmt.Errorf("encode data %v: %s", res, err)
			}
			err = b.Put([]byte(res.Id), buf.Bytes())
			if err != nil {
				return fmt.Errorf("put data: %s", err)
			}
			//log.Println("put key", res.Id)
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
