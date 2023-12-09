package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Listing struct {
	name     string
	features []string
	lat, lon float64
	price    uint
}

func scrape_supercasa() []Listing {
	fmt.Println("Scraping supercasa.pt")

	// iterate through super casa website
	// get all the listings for all the pages
	max_pages2scrape := 5
	var listings []Listing
	for i := 0; i < max_pages2scrape; i++ {
		c := colly.NewCollector()
		url := "https://supercasa.pt/comprar-casas/lisboa/alvalade/pagina-" + strconv.Itoa(i)
		c.OnHTML(".property", func(e *colly.HTMLElement) {
			// Extract all text content within the current element
			textContent := e.Text

			// Extract property ID
			propertyID := e.Attr("id")

			// Extract property price
			propertyPrice := e.ChildText(".property-price span")

			// Extract property features
			var features []string
			e.ForEach(".property-features span", func(_ int, el *colly.HTMLElement) {
				features = append(features, el.Text)
			})

			// Extract latitude and longitude using a regular expression
			re := regexp.MustCompile(`latitude":([\d.-]+),"longitude":([\d.-]+)`)
			match := re.FindStringSubmatch(textContent)
			var latitude, longitude string
			if len(match) == 3 {
				latitude = match[1]
				longitude = match[2]
			}

			// Print the captured information
			fmt.Printf("Property ID: %s\n", propertyID)
			fmt.Printf("Property Price: %s\n", propertyPrice)
			fmt.Printf("Property Features: %s\n", features)
			fmt.Printf("Latitude: %s\n", latitude)
			fmt.Printf("Longitude: %s\n", longitude)

			// type conversions from string to int
			propertyPrice = strings.Split(propertyPrice, "\n")[0]
			propertyPrice = strings.ReplaceAll(propertyPrice, ".", "")
			propertyPrice = strings.TrimSuffix(propertyPrice, "€")
			propertyPrice = strings.TrimSpace(propertyPrice)
			propertyPriceInt, err := strconv.ParseUint(propertyPrice, 10, 64)
			if err != nil {
				log.Fatal("Price string to int conversion failed!")
			}

			latitudeF64, err := strconv.ParseFloat(latitude, 64)
			if err != nil {
				log.Fatal("Latitude string to float conversion failed!")
			}

			longitudeF64, err := strconv.ParseFloat(longitude, 64)
			if err != nil {
				log.Fatal("Longitude string to float conversion failed!")
			}

			listings = append(listings, Listing{name: propertyID, features: features, price: uint(propertyPriceInt), lat: latitudeF64, lon: longitudeF64})
			fmt.Println("----------------------------------------")
		})

		// Visit the URL and start scraping
		err := c.Visit(url)
		if err != nil {
			log.Fatal(err)
		}
	}

	return listings
}

func main() {
	fmt.Println("Welcome to realista!")
	listings := scrape_supercasa()
	fmt.Println("Out of scrape supercasa!\n-------------------------------------------")
	fmt.Println(listings[len(listings)-1])
}
