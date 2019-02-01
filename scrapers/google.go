package scrapers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

// Location contains the latitude & longitude
type Location struct {
	Lat float64
	Lng float64
}

// GetImages searches on google for images & returns an array of image urls
func GetImages(query string, safe bool) ([]string, error) {
	// Replace spaces with '+'
	query = strings.Replace(query, " ", "+", -1)

	url := "http://images.google.com/search?tbm=isch&q=" + query

	if safe {
		url += "&safe=on"
	}

	res, err := fetch(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return nil, err
	}

	var images []string

	doc.Find(".rg_di .rg_meta").Each(func(i int, token *goquery.Selection) {
		imageJSON := token.Text()
		imageURL := gjson.Get(imageJSON, "ou").String()

		if len(imageURL) > 0 {
			images = append(images, imageURL)
		}
	})

	return images, nil
}

// GetLocation returns a location based on the google maps API
func GetLocation(query string, apiKey string) (Location, bool) {
	res, err := http.Get("https://maps.google.com/maps/api/geocode/json?address=" + query + "&key=" + apiKey)

	if err != nil {
		return Location{0, 0}, false
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return Location{0, 0}, false
	}

	defer res.Body.Close()

	json := string(body)

	if !gjson.Get(json, "results.0.geometry.location").Exists() {
		return Location{0, 0}, false
	}

	location := Location{
		Lat: gjson.Get(json, "results.0.geometry.location.lat").Num,
		Lng: gjson.Get(json, "results.0.geometry.location.lng").Num,
	}

	return location, true
}
