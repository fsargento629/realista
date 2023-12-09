package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

type Listing struct {
	name, url, features, description string
	lat, lon                         float32
	price                            uint
}

func scrape_supercasa() {
	fmt.Println("Scraping supercasa.pt")

	// iterate through super casa website
	// get all the listings for all the pages
	max_pages2scrape := 5
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
			fmt.Println("----------------------------------------")
		})

		// Visit the URL and start scraping
		err := c.Visit(url)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	fmt.Println("Welcome to realista!")
	scrape_supercasa()
}
