package killgrave

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestMatcherBySchema(t *testing.T) {
	bodyA := ioutil.NopCloser(bytes.NewReader([]byte("{\"type\": \"gopher\"}")))
	bodyB := ioutil.NopCloser(bytes.NewReader([]byte("{\"type\": \"cat\"}")))
	emptyBody := ioutil.NopCloser(bytes.NewReader([]byte("")))

	schemaGopherFile := "test/testdata/imposters/schemas/type_gopher.json"
	schemaCatFile := "test/testdata/imposters/schemas/type_cat.json"
	schemeFailFile := "test/testdata/imposters/schemas/type_gopher_fail.json"

	requestWithoutSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: nil,
	}

	requestWithSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: &schemaGopherFile,
	}

	requestWithNonExistingSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: &schemaCatFile,
	}

	requestWithWrongSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: &schemeFailFile,
	}

	okResponse := Response{Status: http.StatusOK}

	var matcherData = []struct {
		name string
		fn   mux.MatcherFunc
		req  *http.Request
		res  bool
	}{
		{"imposter without request schema", MatcherBySchema(Imposter{Request: requestWithoutSchema, Response: okResponse}), &http.Request{Body: bodyA}, true},
		{"correct request schema", MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), &http.Request{Body: bodyA}, true},
		{"incorrect request schema", MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), &http.Request{Body: bodyB}, false},
		{"non-existing schema file", MatcherBySchema(Imposter{Request: requestWithNonExistingSchema, Response: okResponse}), &http.Request{Body: bodyB}, false},
		{"malformatted schema file", MatcherBySchema(Imposter{Request: requestWithWrongSchema, Response: okResponse}), &http.Request{Body: bodyB}, false},
		{"empty body with required schema file", MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), &http.Request{Body: emptyBody}, false},
	}

	for _, tt := range matcherData {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.fn(tt.req, nil)
			if res != tt.res {
				t.Fatalf("error while matching by request schema - expected: %t, given: %t", tt.res, res)
			}
		})

	}
}
