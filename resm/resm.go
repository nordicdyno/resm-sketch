package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"
	// 3rd party libs
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	// project's libs
	"github.com/nordicdyno/resm-sketch/store"
	"github.com/nordicdyno/resm-sketch/store/inbolt"
	"github.com/nordicdyno/resm-sketch/store/inmemory"
)

const (
	DefaultResourcesLimit = 10
	HandlerContextKey     = "Handler"
	BadRequestMessage     = "Bad Request.\n"
)

// hack
var _ = spew.Config

type positiveIntVal int

func (nv *positiveIntVal) String() string { return "PositiveInteger" }

func (nv *positiveIntVal) Set(s string) error {
	value, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}
	*nv = positiveIntVal(value)
	if *nv <= 0 {
		return fmt.Errorf("value should be positive")
	}
	return nil
}

var (
	file     = flag.String("file", "", "bolt DB file")
	limitVal = positiveIntVal(DefaultResourcesLimit)
	bind     = flag.String("bind", ":9090", "[host]:port where to serve on")
	verbose  = flag.Bool("verbose", false, "chatty mode")
)

func init() {
	flag.Var(&limitVal, "limit", "resources limit")

	flag.Parse()
}

func main() {
	if *verbose {
		fmt.Println("Run on", *bind)
	}

	rh := NewResourceHandler(int(limitVal), *file)
	rh.LogStackTrace = *verbose
	//os.Args = nil // monkeypatch: hide commandline from expvar
	http.Handle("/", rh)

	err := http.ListenAndServe(*bind, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getHandler(r *http.Request) *ResourceHandler {
	hVal := context.Get(r, HandlerContextKey)
	h, ok := hVal.(*ResourceHandler)
	if !ok || h == nil {
		panic("Global Handler not found in request context")
	}
	return h
}

type ResourceHandler struct {
	Router        *mux.Router
	Logger        *log.Logger
	LogStackTrace bool

	Storage store.ResourceAllocater
}

func NewResourceHandler(limit int, filePath string) *ResourceHandler {
	var (
		storage store.ResourceAllocater
		err     error
	)
	switch filePath {
	case "":
		storage, err = inmemory.NewStorage(limit)
		if err != nil {
			log.Fatalln("Create inmemory storage failed", err)
		}
	default:
		storage, err = inbolt.NewStorage(filePath, limit)
		if err != nil {
			log.Panicln("Bolt db creation failed ", err)
		}
	}

	r := mux.NewRouter()

	getHandlers := []struct {
		path string
		fn   func(http.ResponseWriter, *http.Request)
	}{
		{"/allocate/{user}", allocateResourceUser},
		{"/deallocate/{resource_id}", deallocateResource},
		{"/reset", resetResources},
		{"/list", listResources},
		{"/list/{user}", listUserResources},
	}
	for _, gH := range getHandlers {
		r.HandleFunc(gH.path, gH.fn).Methods("GET")
	}
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(BadRequestMessage))
	})

	return &ResourceHandler{
		Router:  r,
		Storage: storage,
	}
}

func (h *ResourceHandler) LogLine(code int, start *time.Time, req *http.Request) {
	message := ""
	if h.LogStackTrace {
		trace := debug.Stack()
		message = fmt.Sprintf("\n%s", trace)
	}
	h.Logger.Println(
		code,
		&start,
		time.Now().Sub(*start),
		req.Method,
		req.URL.RequestURI(),
		message,
	)
}

func (h *ResourceHandler) ServeHTTP(wOrig http.ResponseWriter, r *http.Request) {
	w := &ResponseWrapper{wOrig, 200, false}
	start := time.Now()

	// set a default Logger
	h.Logger = log.New(os.Stderr, "", log.LstdFlags)

	// catch user code's panic, and convert to http response
	// (this does not use the JSON error response on purpose)
	defer func() {
		if rec := recover(); rec != nil {
			message := "Internal Server Error"
			http.Error(w, message, http.StatusInternalServerError)

			// log response
			h.LogLine(http.StatusInternalServerError, &start, r)
		}
	}()

	context.Set(r, HandlerContextKey, h)
	// defer context.Clear(r) <- no need with gorilla's mux/pat
	h.Router.ServeHTTP(w, r)
	h.LogLine(w.Code, &start, r)
}

type ResponseWrapper struct {
	http.ResponseWriter
	Code        int // the HTTP response code from WriteHeader
	wroteHeader bool
}

func (rw *ResponseWrapper) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.Code = code
	}
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}
