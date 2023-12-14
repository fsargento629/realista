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
	bairro, energy_rating, source string
	features                      []string
	lat, lon                      float64
	price, area, rooms, id        uint
	// pos processing values
	price_per_m2, price_offset float32 // price offet is 1 if it is the same, 2 if double, etc
}

type Bairro struct {
	name, url string
	// average_price_per_m2 float32
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

// Scrapes supercasa.pt and returns an array of Listings.
// The
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

			// property id from string to uint
			id, err := strconv.Atoi(strings.TrimPrefix(propertyID, "property_"))
			if err != nil {
				id = 0
			}

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
			listings = append(listings, Listing{
				id: uint(id), bairro: bairro.name, features: features, source: "super_casa",
				energy_rating: energy_rating, price: uint(propertyPriceInt),
				rooms: uint(rooms), area: uint(area),
				lat: latitudeF64, lon: longitudeF64})
			fmt.Println("----------------------------------------")
		})
		for i := 0; i < max_pages_per_bairro; i++ {

			url := bairro.url + "/pagina-" + strconv.Itoa(i)
			// Visit the URL and start scraping
			// Visit() starts the scraping process, whcih will call the on HTML callback we set up
			err := c.Visit(url)
			if err != nil {
				fmt.Println("scrape_supercasa scraped " + strconv.Itoa(i) + " pages from bairro " + bairro.name)
				fmt.Println(err)
				break
			}
		}
		fmt.Printf("Scraped %d house from super_casa.pt", len(listings))
	}

	return listings
}

// add statistics to each listing
// (price/m2, price offset to neighbourhood average (%), price offset to selection)
// todo: this should also return the bairro data. we should have a struct for that. Maybe just reuse bairro.
func post_processing(listings []Listing) []Listing {

	var newListings []Listing
	// price per area
	for _, listing := range listings {
		if listing.area > 0 {
			listing.price_per_m2 = float32(listing.price / listing.area)
			newListings = append(newListings, listing)
		} else {
			// remove listing if area is not positive.
			fmt.Println("Listing " + strconv.FormatUint(uint64(listing.id), 10) + " had non-positive area")
		}
	}

	// hasmaps to support mean calculation
	bairro_average_ppm2s := make(map[string]float32)
	bairro_counters := make(map[string]uint)
	bairro_price_sum := make(map[string]uint)

	// iterate all listing toget the total price sum for each bairro
	for _, listing := range newListings {
		bairro_price_sum[listing.bairro] += uint(listing.price_per_m2)
		bairro_counters[listing.bairro]++
	}
	// divide price sums by counters to get averages
	for bairro, price_sum := range bairro_price_sum {
		bairro_average_ppm2s[bairro] = float32(float64(price_sum) / float64(bairro_counters[bairro]))
		fmt.Printf("Bairro:%s | Average PPMsqr:%f\n", bairro, bairro_average_ppm2s[bairro])
	}

	// iterate new listings to get price offset to neighbourhood
	for i, listing := range newListings {
		newListings[i].price_offset = listing.price_per_m2 / bairro_average_ppm2s[listing.bairro]
	}

	return newListings

}

// helper function to quickly print data about listing.
// this should be moved to a Listing specfic module
func print_listing(listing Listing) {
	fmt.Println("----------------------------")
	fmt.Printf("ID: %d\nBairro:%s\nCE:%s\nLat:%f ; Lon:%f\nPrice(k eur):%d\nArea:%d\nRooms:%d\nPPMsqr:%f\nBairro price offset:%f\n",
		listing.id, listing.bairro, listing.energy_rating, listing.lat, listing.lon, listing.price/1000, listing.area,
		listing.rooms, listing.price_per_m2, listing.price_offset)
	fmt.Println("----------------------------")

}

// ------------------------------------------------------
// MAIN

func main() {
	fmt.Println("Welcome to realista!")
	listings := scrape_supercasa()
	listings = post_processing(listings)
	print_listing(listings[0])
	print_listing(listings[10])
	print_listing(listings[20])
	print_listing(listings[30])
	print_listing(listings[40])
	print_listing(listings[50])

	fmt.Println("Out of scrape supercasa!\n-------------------------------------------")
	fmt.Println(listings[len(listings)-1])
}
