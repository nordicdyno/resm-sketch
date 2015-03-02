package main

import (
	"encoding/json"
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

	//log.Println("allocateResourceUser with id=", id)
	///spew.Dump(mux.Vars(r))
}

func deallocateResource(w http.ResponseWriter, r *http.Request) {
	h := getHandler(r)
	id, _ := mux.Vars(r)["resource_id"]
	err := h.Storage.Deallocate(id)
	//log.Println("deallocateResource err => ", err)
	if err != nil {
		//log.Println("Resource", id, "not allocated", err)
		if err == store.ErrResourcesNotFound {
			//log.Println("write header", http.StatusNotFound)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(notAllocatedStr))
		} else {
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

	//log.Println("/reset OK, header:", http.StatusNoContent)
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

	var b []byte
	b, err = json.Marshal(info)
	if err != nil {
		panic(err)
	}

	//log.Println(string(b))
	//log.Println("/reset OK, header:", http.StatusNoContent)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
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

	//log.Println(string(b))
	//log.Println("/reset OK, header:", http.StatusNoContent)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
