package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"

	"github.com/nordicdyno/resm-sketch/store"
)

var (
	outOfResourcesStr = "Out of resources.\n"
	notAllocatedStr   = "Not allocated.\n"
)
var (
	_ = spew.Config
	_ = log.Prefix()
)

func allocateResourceUser(w http.ResponseWriter, r *http.Request) {
	h := getHandler(r)
	user, ok := mux.Vars(r)["user"]
	if !ok || len(user) < 1 {
		// write 404 or 503 bad request?
		//log.Println("user not found")
		return
	}
	id, err := h.Storage.Allocate(user)
	if err != nil {
		if err == store.ErrResourcesIsOver {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(outOfResourcesStr))
			return
		}
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func deallocateResource(w http.ResponseWriter, r *http.Request) {
	h := getHandler(r)
	id, _ := mux.Vars(r)["resource_id"]
	err := h.Storage.Deallocate(id)
	if err != nil {
		switch err {
		case store.ErrResourcesIsFree:
			fallthrough
		case store.ErrResourcesNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(notAllocatedStr))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func resetResources(w http.ResponseWriter, r *http.Request) {
	h := getHandler(r)
	if err := h.Storage.Reset(); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func listResources(w http.ResponseWriter, r *http.Request) {
	var (
		info *store.ResourcesInfo
		err  error
	)
	h := getHandler(r)
	if info, err = h.Storage.List(); err != nil {
		panic(err)
	}

	var bDeallocated []byte
	bDeallocated, err = json.Marshal(info.Deallocated)
	if err != nil {
		panic(err)
	}

	bAllocated := new(bytes.Buffer)
	for n, resource := range info.Allocated {
		if n != 0 {
			fmt.Fprintf(bAllocated, ",")
		}
		fmt.Fprintf(bAllocated, `"%s":"%s"`, resource.Id, resource.User)
	}

	listJSONstr := fmt.Sprintf(`{"allocated":{%s}, "deallocated": %s}`,
		bAllocated.String(), string(bDeallocated),
	)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(listJSONstr))
}

func listUserResources(w http.ResponseWriter, r *http.Request) {
	var (
		list store.ResourcesList
		err  error
	)
	h := getHandler(r)
	user, _ := mux.Vars(r)["user"]
	if list, err = h.Storage.ListByUser(user); err != nil {
		panic(err)
	}

	var b []byte
	b, err = json.Marshal(list)
	if err != nil {
		panic(err)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
