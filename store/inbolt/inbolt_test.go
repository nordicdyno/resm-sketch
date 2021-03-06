package inbolt

import (
	"testing"

	"github.com/nordicdyno/resm-sketch/store"
	storetest "github.com/nordicdyno/resm-sketch/store/test"
)

var sTen store.ResourceAllocater

func TestMain(t *testing.T) {
	var err error
	sTen, err = NewStorage("test.db", 10)
	if err != nil {
		t.Error("Create inmemory storage failed", err)
	}
	storetest.ConfigureMain()
}

func TestList(t *testing.T) {
	storetest.TenTestList(sTen, t)
}

func TestAllocateByUser(t *testing.T) {
	storetest.TenTestAllocateByUser(sTen, t)

}

func TestDeallocateById(t *testing.T) {
	storetest.TenTestDeallocateById(sTen, t)
}

func TestReset(t *testing.T) {
	storetest.TenTestReset(sTen, t)
}
