package kount

import (
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/go-resty/resty"
)

type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) write(buf []byte) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.Write(buf)
}

// Car is the data we want from the response
type Car struct {
	RemoteID   int
	BatchID    int
	Title      string
	Price      int
	DatePosted string
	Seen       string
}

// CheckThenPanic scan of the inventory for 'Subaru' on CL
func CheckThenPanic() (results []byte, numResults int, err error) {
	BatchID := 1

	//initial request
	resp, err := resty.R().Get("https://boise.craigslist.org/search/cta?query=subaru")

	// get number of pages to scrape, and results from the first page
	cars, numResults, err := ParsePage(resp.String(), BatchID)
	check(err)

	//paginate
	for i := 120; i < numResults; i += 120 {
		resp, err := resty.R().Get("https://boise.craigslist.org/search/cta?query=subaru&s=" + strconv.Itoa(i))
		check(err)
		additionalCars, _, err := ParsePage(resp.String(), BatchID)
		check(err)

		cars = append(cars, additionalCars...) //Ok, that (...) is awesome
	}

	json, err := json.Marshal(cars)
	check(err)

	return json, numResults, err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ParsePage from the html
func ParsePage(html string, BatchID int) ([]Car, int, error) {
	var results []Car

	doc := soup.HTMLParse(html)

	NumResults, err := strconv.Atoi(doc.Find("span", "class", "totalcount").Text())
	check(err)

	cars := doc.Find("ul", "class", "rows").FindAll("li")

	for _, car := range cars {
		RemoteID, _ := strconv.Atoi(car.Attrs()["data-pid"])
		Price := FindPrice(car)

		Title := car.Find("a", "class", "result-title").Text()
		DatePosted := car.Find("time", "class", "result-date").Attrs()["datetime"]
		Seen := time.Now().String()
		result := Car{RemoteID, BatchID, Title, Price, DatePosted, Seen}

		results = append(results, result)
	}

	return results, NumResults, nil
}

// FindPrice without errors
func FindPrice(car soup.Root) int {
	Price := 0
	PriceText := car.Find("span", "class", "result-price")

	if PriceText.Error == nil {
		Price, _ = strconv.Atoi(strings.Trim(PriceText.Text(), "$")) // TODO: don't ignore this error
	}

	return Price
}
