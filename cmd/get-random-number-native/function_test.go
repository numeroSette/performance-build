package getrandomnumbernative

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func GetRandomNumberNativeTest(b *testing.B) {

	// References
	// https://blog.questionable.services/article/testing-http-handlers-go/

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/get-random-number-native", nil)
	if err != nil {
		b.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetRandomNumberNative)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		b.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `^\{\"random_number"\:\"\d+\"\}$`
	matched, _ := regexp.MatchString(`^\{\"random_number"\:\"\d+\"\}$`, rr.Body.String())
	if matched {
		b.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func BenchmarkGetRandomNumberNativeTest(b *testing.B) { GetRandomNumberNativeTest(b) }
