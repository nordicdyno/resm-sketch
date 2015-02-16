package inmemory

import (
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
    "github.com/nordicdyno/resm-sketch/store"
    "strconv")

var sTen *Storage

func TestMain(t *testing.T) {
	var err error
	sTen, err = NewStorage(10)
    if err != nil {
        t.Error("Create inmemory storage failed", err)
    }
}

var _ = spew.Config

func TestList(t *testing.T) {
	l, err := sTen.List()
	if err != nil {
		t.Error("call failed", err)
	}
	// spew.Dump(l)
	if len(l.Allocated) != 0 {
		t.Error("allocation not expected here", err)
	}

	if len(l.Deallocated) != 10 {
		t.Error(10, "items should be allocated")
	}
	_ = l
}

func TestAllocateByUser(t *testing.T) {
	for i := 0; i < 10; i++ {
        user := "t" + strconv.Itoa(i)
		id, err := sTen.Allocate(user)
		if err != nil {
			t.Error("Allocation failed on step", i, err)
		}
        if id == "" {
            t.Error("Allocation failed on step", i, "id is empty")
        }
		log.Println(id, "resource created for user", user)
	}

	l, err := sTen.List()
	if err != nil {
		t.Error("call failed", err)
	}
	if len(l.Deallocated) != 0 {
		t.Error(len(l.Deallocated), "items found, but they should be out of stock")
	}
}

func TestDeallocateById(t *testing.T) {
	l, err := sTen.List()
	if err != nil {
		t.Error("call failed", err)
	}
	//spew.Dump(l)
	for user, val := range l.Allocated {
		for _, id := range val {
			err := sTen.Deallocate(id)
			if err != nil {
				t.Error("Deallocate failed for resource", id, err)
			}
			log.Println(id, "resource deallocated for user", user)
		}
	}
    firstId := store.GenResourcesIds(1)[0]
    err = sTen.Deallocate(firstId)
    if err == nil {
        t.Error("Deallocate not failed for resource", firstId, err)
    }
	//sTen.Dump()
	//log.Fatalln("bye!")

	l, err = sTen.List()
	if err != nil {
		t.Error("call failed", err)
	}
	//sTen.Dump()
	//log.Println("List result:")
	//spew.Dump(l)
	if len(l.Deallocated) != 10 {
		t.Error(10, "items expected, but found ", len(l.Deallocated))
	}
}
