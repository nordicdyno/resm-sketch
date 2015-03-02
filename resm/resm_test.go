package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"testing"
	//"strings"
	//"fmt"
	//"log"
	"log"
	//    "github.com/nordicdyno/resm-sketch/store"
)

type AnyJSON interface{}

type handlerTest struct {
	Desc      string
	Method    string
	Path      string
	Params    url.Values
	Status    int
	Match     map[string]bool
	MatchJSON string
	StopOnMe  bool
}

func tablesRun(t *testing.T, rh *ResourceHandler, tests []handlerTest) {
	for _, test := range tests {

		record := httptest.NewRecorder()
		req := &http.Request{
			Method: test.Method,
			URL:    &url.URL{Path: test.Path},
		}
		if test.Params != nil {
			req.Form = test.Params
		}

		rh.ServeHTTP(record, req)

		got, want := record.Code, test.Status
		_ = log.Prefix()
		//log.Println("Got code:", got, record.Code)
		//log.Println("Wait code:", want, test.Status)

		if got != want {
			//log.Println(got, "!=", want, "?")
			//t.Errorf("%s: response code = %d, want %d", test.Desc, got, want)
			t.Fatalf("%s: response code = %d, want %d", test.Desc, got, want)
		}

		if test.Match != nil {
			for re, match := range test.Match {
				if got := regexp.MustCompile(re).Match(record.Body.Bytes()); got != match {
					t.Errorf("%s: %q ~ /%s/ = %v, want %v", test.Desc, record.Body, re, got, match)
				}
			}
		}

		if test.MatchJSON != "" {
			var exp, got AnyJSON
			gotBytes := record.Body.Bytes()
			//log.Println(test.MatchJSON, "cmp", string(gotBytes))

			json.Unmarshal(gotBytes, &got)
			json.Unmarshal([]byte(test.MatchJSON), &exp)
			//log.Println(exp, "VS", got, "->", reflect.DeepEqual(got, exp))

			if !reflect.DeepEqual(got, exp) {
				t.Errorf("%s: Got: %v, Expect: %s", test.Desc, record.Body, test.MatchJSON)
			}
		}

		if test.StopOnMe {
			return
		}
	}
}

//
//var handlers []store.ResourceAllocater
//
//func TestMain(t *testing.T) {
//    rh_bolt := NewResourceHandler(2, "test.db")
//
//}

func TestAllocateResetAndList(t *testing.T) {
	tests := []handlerTest{
		{
			Desc:   "allocate 1 OK",
			Method: "GET",
			Path:   "/allocate/my",
			Status: 201,
			Match: map[string]bool{
				`^r\d+$`: true,
			},
		},
		{
			Desc:   "allocate 2 OK",
			Method: "GET",
			Path:   "/allocate/him",
			Status: 201,
			Match: map[string]bool{
				`^r\d+$`: true,
			},
		},
		{
			Desc:      "list n1",
			Method:    "GET",
			Path:      "/list",
			Status:    200,
			MatchJSON: `{"Allocated":{"r0":"my","r1":"him"},"Deallocated":[]}`,
		},
		{
			Desc:      "list by unknown",
			Method:    "GET",
			Path:      "/list/unknown",
			Status:    200,
			MatchJSON: `[]`,
			//StopOnMe:  true,
		},
		{
			Desc:      "list by him",
			Method:    "GET",
			Path:      "/list/him",
			Status:    200,
			MatchJSON: `["r1"]`,
		},
		{
			Desc:   "allocate Fail",
			Method: "GET",
			Path:   "/allocate/me",
			Status: 503,
			Match: map[string]bool{
				`^r\d+`: false,
			},
			//StopOnMe:  true,
		},
		{
			Desc:   "Reset",
			Method: "GET",
			Path:   "/reset",
			Status: 204,
			Match: map[string]bool{
				`^.+$`: false,
			},
		},
		{
			Desc:   "allocate 1 again OK",
			Method: "GET",
			Path:   "/allocate/my",
			Status: 201,
			Match: map[string]bool{
				`^r\d+$`: true,
			},
		},
	}

	rh := NewResourceHandler(2, "")
	tablesRun(t, rh, tests)

	rh_bolt := NewResourceHandler(2, "test.db")
	tablesRun(t, rh_bolt, tests)
}

func TestDeallocate(t *testing.T) {
	tests := []handlerTest{
		{
			Desc:   "deallocate not allocated resource",
			Method: "GET",
			Path:   "/deallocate/r0",
			Status: 404,
			Match: map[string]bool{
				`^Not allocated`: true,
			},
		},
		{
			Desc:   "allocate resorce for deallocation",
			Method: "GET",
			Path:   "/allocate/r0",
			Status: 201,
		},
		{
			Desc:   "deallocate OK",
			Method: "GET",
			Path:   "/deallocate/r0",
			Status: 204,
			Match: map[string]bool{
				`^.+$`: false,
			},
		},
		{
			Desc:   "Try deallocate not existed resource",
			Method: "GET",
			Path:   "/deallocate/any",
			Params: nil,
			Status: 404,
			Match: map[string]bool{
				`^Not allocated`: true,
			},
			//StopOnMe: true,
		},
	}

	rh := NewResourceHandler(1, "")
	tablesRun(t, rh, tests)

	rh_bolt := NewResourceHandler(1, "test2.db")
	tablesRun(t, rh_bolt, tests)
}
