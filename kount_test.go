package kount

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestCheckThenPanic(t *testing.T) {
	// Test that we're getting a response from CL
	results, numResults, err := CheckThenPanic()

	if err != nil {
		t.Errorf("There was an error with the CheckThenPanic method")
	}

	if len(results) == 0 {
		t.Errorf("Response came back empty, no results")
	}

	var c []Car
	err = json.Unmarshal(results, &c)

	//test to see if string is JSON
	if err != nil {
		t.Errorf("Response was not JSON")
	}

	//check to see that we're getting close to the amount that we paginated for
	if len(c) < numResults-10 {
		t.Errorf("Not enough results in response Expected: %d Got: %d", numResults, len(c))
	}
}

func BenchmarkCheckThenPanic(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _, _ = CheckThenPanic()
	}
}

func TestParsePage(t *testing.T) {
	// test that we're getting cars back
	html, err := ioutil.ReadFile("test.html")
	check(err)

	cars, NumResults, err := ParsePage(string(html), 1)
	check(err)

	if NumResults < 120 {
		t.Errorf("Did not fetch a full page of results")
	}

	if len(cars) != 120 {
		t.Errorf("cars is expected to be 120, got: %s", string(len(cars)))
	}

	//check that all cars have numeric prices, titles, etc (all fields)
}

func BenchmarkParsePage(b *testing.B) {
	html, err := ioutil.ReadFile("test.html")
	check(err)

	for n := 0; n < b.N; n++ {
		ParsePage(string(html), 1)
	}
}
