package main

// ------------------------------------------------------
// IMPORTS
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// ------------------------------------------------------
// STRUCTS

type Listing struct {
	name, bairro, energy_rating string
	features                    []string
	lat, lon                    float64
	price, area, rooms          uint
	// to implement: ID as a algbebraic enum
	// area, Energy rating, number of rooms
	//
}

type Bairro struct {
	name, url string
	// in the future this should also have data about the average price and other economic metrics
}

// ------------------------------------------------------
// FUNCTIONS

// returns the bairro info stored in config_files/bairros.csv
// this is just name,url pairs so far
func bairros2scrape() []Bairro {
	var bairro_list []Bairro

	// read csv file with bairro info and populate the array
	// open file logic
	file_name := "/home/fsargento/go/projects/realista/config_files/bairros.csv"
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal("Error while opening the file!!", err)
	}
	defer file.Close()

	// read the csv
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading csv file!")
	}
	// convert each line to question, answer and push them to the vec
	for _, line := range lines {
		bairro_list = append(bairro_list, Bairro{name: line[0], url: line[1]})
	}
	fmt.Print("Got bairro list: ")
	fmt.Println(bairro_list)
	return bairro_list
}

func scrape_supercasa() []Listing {
	fmt.Println("Scraping supercasa.pt")
	var listings []Listing

	// iterate thorugh each bairro
	// each bairro has several pages, each with several listings
	max_pages_per_bairro := 65
	bairro_list := bairros2scrape()
	for _, bairro := range bairro_list {
		fmt.Println("On bairro:" + bairro.name)
		c := colly.NewCollector()
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

			// parse features from string to num of rooms, area, CE rating
			// The features array always has therse values in order
			// number of rooms
			var rooms int
			re = regexp.MustCompile(`(\d+) quarto`)
			match = re.FindStringSubmatch(strings.Join(features, " "))
			if len(match) == 2 {
				rooms, err = strconv.Atoi(match[1])
				if err != nil {
					log.Fatal("Unable to parse number of rooms from string to int")
				}
			} else {
				fmt.Println("Unable to find number of rooms in features")
				rooms = 0
			}

			// area
			var area int
			re = regexp.MustCompile(`(\d+) m²`)
			match = re.FindStringSubmatch(strings.Join(features, " "))
			if len(match) == 2 {
				area, err = strconv.Atoi(match[1])
				if err != nil {
					log.Fatal("Unable to parse area from string to int")
				}
			} else {
				fmt.Println("Unable to find area in features")
				area = 0
			}

			// energy rating
			var energy_rating string
			re = regexp.MustCompile(`C\.E\.: ?(\w)`)
			match = re.FindStringSubmatch(strings.Join(features, " "))
			if len(match) == 2 {
				energy_rating = match[1]
			} else {
				fmt.Println("Unable to find energy rating in features")
				energy_rating = "?"
			}

			// latitude conversion
			latitudeF64, err := strconv.ParseFloat(latitude, 64)
			if err != nil {
				log.Fatal("Latitude string to float conversion failed!")
			}

			// longitude conversion
			longitudeF64, err := strconv.ParseFloat(longitude, 64)
			if err != nil {
				log.Fatal("Longitude string to float conversion failed!")
			}

			// Finally, append this listing to the listing list
			listings = append(listings, Listing{name: propertyID, bairro: bairro.name, features: features, energy_rating: energy_rating, price: uint(propertyPriceInt), rooms: uint(rooms), area: uint(area), lat: latitudeF64, lon: longitudeF64})
			fmt.Println("----------------------------------------")
		})
		for i := 0; i < max_pages_per_bairro; i++ {

			url := bairro.url + "/pagina-" + strconv.Itoa(i)
			// Visit the URL and start scraping
			// Visit() starts the scraping process, whcih will call the on HTML callback we set up
			err := c.Visit(url)
			if err != nil {
				fmt.Println("scrape_supercasa scraped " + strconv.Itoa(i) + " pages")
				fmt.Println(err)
				break
			}
		}
	}

	return listings
}

// ------------------------------------------------------
// MAIN

func main() {
	fmt.Println("Welcome to realista!")
	listings := scrape_supercasa()
	fmt.Println("Out of scrape supercasa!\n-------------------------------------------")
	fmt.Println(listings[len(listings)-1])
}
