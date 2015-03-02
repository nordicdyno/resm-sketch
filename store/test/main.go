package test

import (
	"log"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/nordicdyno/resm-sketch/store"
	"strconv"
)

//var sTen store.ResourceAllocater

func ConfigureMain() {
	if len(os.Getenv("VERBOSE")) < 1 {
		logfile, _ := os.OpenFile(os.DevNull, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		log.SetOutput(logfile)
	}
}

var _ = spew.Config

func TenTestList(sTen store.ResourceAllocater, t *testing.T) {
	l, err := sTen.List()
	if err != nil {
		t.Error("call failed", err)
	}
	//spew.Dump(l)

	if len(l.Allocated) != 0 {
		t.Error("allocation not expected here", err)
	}

	if len(l.Deallocated) != 10 {
		t.Error(10, "items should be allocated")
	}
	_ = l
}

func TenTestAllocateByUser(sTen store.ResourceAllocater, t *testing.T) {
	for i := 0; i < 10; i++ {
		id, err := sTen.Allocate("t")
		if err != nil {
			t.Error("Allocation failed on step", i, err)
		}
		log.Println(id, "resource created for user 't'")
	}

	l, err := sTen.List()
	if err != nil {
		t.Error("call failed", err)
	}
	if len(l.Deallocated) != 0 {
		t.Error(len(l.Deallocated), "items found, but they should be out of stock")
	}

}

func TenTestDeallocateById(sTen store.ResourceAllocater, t *testing.T) {
	l, err := sTen.List()
	if err != nil {
		t.Error("call failed", err)
	}

	for user, val := range l.Allocated {

		for _, id := range val {
			err := sTen.Deallocate(id)
			if err != nil {
				t.Error("Deallocate failed for resource", user, err)
			}
			log.Println(id, "resource deallocated for user", user)
		}
	}

	l, err = sTen.List()
	if err != nil {
		t.Error("call failed", err)
	}
	if len(l.Deallocated) != 10 {
		t.Error(10, "items expected, but found ", len(l.Deallocated))
	}
}

func TenTestReset(sTen store.ResourceAllocater, t *testing.T) {
	for i := 0; i < 10; i++ {
		id, err := sTen.Allocate("t" + strconv.Itoa(i))
		if err != nil {
			t.Error("Allocation failed on step", i, err)
		}
		log.Println(id, "resource created for user 't'")
	}

	err := sTen.Reset()
	if err != nil {
		t.Error("Reset() call failed", err)
	}
	l, err := sTen.List()
	if err != nil {
		t.Error("List() call failed", err)
	}

	if len(l.Allocated) != 0 {
		t.Error(len(l.Allocated), "allocated items found, but all items should be free")
	}
}
