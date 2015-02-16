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
    r.HandleFunc("/allocate/{user}", allocateResourceUser)
    r.HandleFunc("/deallocate/{resource_id}", deallocateResource)
    r.HandleFunc("/reset", resetResources)
    r.HandleFunc("/list", listResources)
    r.HandleFunc("/list/{user}", listUserResources)

	//r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return &ResourceHandler{
		Router:  r,
		Storage: storage,
	}
}

func (h *ResourceHandler) LogLine(code int, start *time.Time, req *http.Request) {
	h.Logger.Println(
		code,
		&start,
		time.Now().Sub(*start),
		req.Method,
		req.URL.RequestURI(),
	)

}

func (h *ResourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// set a default Logger
	h.Logger = log.New(os.Stderr, "", log.LstdFlags)

	// catch user code's panic, and convert to http response
	// (this does not use the JSON error response on purpose)
	defer func() {
		if rec := recover(); rec != nil {
			trace := debug.Stack()
			h.Logger.Printf("%s\n%s", r, trace)

			message := "Internal Server Error"
			if h.LogStackTrace {
				message = fmt.Sprintf("%s\n\n%s", r, trace)
			}
			http.Error(w, message, http.StatusInternalServerError)

			// log response
			h.LogLine(http.StatusInternalServerError, &start, r)
		}
	}()

	context.Set(r, HandlerContextKey, h)
	h.Router.ServeHTTP(w, r)
	// FIXME: not always 200!
	h.LogLine(200, &start, r)
}
