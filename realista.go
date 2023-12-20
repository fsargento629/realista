package main

// ------------------------------------------------------
// IMPORTS
import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

// ------------------------------------------------------
// STRUCTS

type Listing struct {
	Bairro, Energy_rating, Source, Url string
	Features                           []string
	Lat, Lon                           float64
	Price, Area, Rooms, Id             uint
	// pos processing values
	Price_per_m2, Price_offset float32 // price offet is 1 if it is the same, 2 if double, etc
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
	file_name := "config_files/bairros.csv"
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

			// Link for listing
			url := e.ChildAttr(".property-list-title a", "href")
			fmt.Printf("Url: %s\n", url)

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
				Id: uint(id), Bairro: bairro.name, Features: features, Source: "super_casa", Url: "https://supercasa.pt" + url,
				Energy_rating: energy_rating, Price: uint(propertyPriceInt),
				Rooms: uint(rooms), Area: uint(area),
				Lat: latitudeF64, Lon: longitudeF64})
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
		if listing.Area > 0 {
			listing.Price_per_m2 = float32(listing.Price / listing.Area)
			newListings = append(newListings, listing)
		} else {
			// remove listing if area is not positive.
			fmt.Println("Listing " + strconv.FormatUint(uint64(listing.Id), 10) + " had non-positive area")
		}
	}

	// hasmaps to support mean calculation
	bairro_average_ppm2s := make(map[string]float32)
	bairro_counters := make(map[string]uint)
	bairro_price_sum := make(map[string]uint)

	// iterate all listing toget the total price sum for each bairro
	for _, listing := range newListings {
		bairro_price_sum[listing.Bairro] += uint(listing.Price_per_m2)
		bairro_counters[listing.Bairro]++
	}
	// divide price sums by counters to get averages
	for bairro, price_sum := range bairro_price_sum {
		bairro_average_ppm2s[bairro] = float32(float64(price_sum) / float64(bairro_counters[bairro]))
		fmt.Printf("Bairro:%s | Average PPMsqr:%f\n", bairro, bairro_average_ppm2s[bairro])
	}

	// iterate new listings to get price offset to neighbourhood
	for i, listing := range newListings {
		newListings[i].Price_offset = listing.Price_per_m2 / bairro_average_ppm2s[listing.Bairro]
	}

	return newListings

}

// scrape all implemented websites and convert them to a csv
func scrape() []Listing {
	listings := scrape_supercasa()
	listings = post_processing(listings)

	// convert data to csv
	listings2csv(listings) // this is working, but we dont need to always call this

	return listings

}

// helper function to quickly print data about listing.
// this should be moved to a Listing specfic module
func print_listing(listing Listing) {
	fmt.Println("----------------------------")
	fmt.Printf("ID: %d\nURL:%s\nBairro:%s\nCE:%s\nLat:%f ; Lon:%f\nPrice(k eur):%d\nArea:%d\nRooms:%d\nPPMsqr:%f\nBairro price offset:%f\n",
		listing.Id, listing.Url, listing.Bairro, listing.Energy_rating, listing.Lat, listing.Lon, listing.Price/1000, listing.Area,
		listing.Rooms, listing.Price_per_m2, listing.Price_offset)
	fmt.Println("----------------------------")

}

// create a csv file with all the listings
func listings2csv(listings []Listing) {
	currentTime := time.Now()
	csvFileName := fmt.Sprintf("data/listings_%s.csv", currentTime.Format("2006_01_02_15_04_05"))

	// Create or open the CSV file
	file, err := os.Create(csvFileName)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header to the CSV file
	header := []string{"id", "price", "area", "rooms", "energyRating", "pricePerM2", "bairro", "priceOffset", "lat", "lon", "url"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing CSV header:", err)
		return
	}

	// Write each listing to the CSV file
	for _, listing := range listings {
		row := []string{
			fmt.Sprintf("%d", listing.Id),
			fmt.Sprintf("%d", listing.Price),
			fmt.Sprintf("%d", listing.Area),
			fmt.Sprintf("%d", listing.Rooms),
			listing.Energy_rating,
			fmt.Sprintf("%.2f", listing.Price_per_m2),
			listing.Bairro,
			fmt.Sprintf("%.2f", listing.Price_offset),
			fmt.Sprintf("%.8f", listing.Lat),
			fmt.Sprintf("%.8f", listing.Lon),
			listing.Url,
		}

		err := writer.Write(row)
		if err != nil {
			fmt.Println("Error writing CSV row:", err)
			return
		}
	}

	fmt.Printf("CSV file '%s' created successfully.\n", csvFileName)
}

// get the array of listings from .csv
func csv2listings(csv_file string) []Listing {

	file, err := os.Open(csv_file)
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
	var listings []Listing
	for _, line := range lines[1:] {
		// Each line should be like:
		//id,price,area,rooms,energyRating,pricePerM2,bairro,priceOffset,lat,lon

		id, _ := strconv.ParseUint(line[0], 10, 16)
		price, _ := strconv.ParseUint(line[1], 10, 32)
		area, _ := strconv.ParseUint(line[2], 10, 32)
		rooms, _ := strconv.ParseUint(line[3], 10, 32)
		price_per_m2, _ := strconv.ParseFloat(line[5], 64)
		price_offset, _ := strconv.ParseFloat(line[7], 64)
		lat, _ := strconv.ParseFloat(line[8], 64)
		lon, _ := strconv.ParseFloat(line[9], 64)

		// append to listings. Should we have a step here to filter out bad parsings?
		listing := Listing{Id: uint(id), Price: uint(price), Area: uint(area), Rooms: uint(rooms), Energy_rating: line[4],
			Price_per_m2: float32(price_per_m2), Bairro: line[6], Price_offset: float32(price_offset), Lat: lat, Lon: lon, Url: line[10]}
		listings = append(listings, listing)
	}

	fmt.Printf("Imported %d listings from %s\n", len(listings), csv_file)
	return listings
}

// API function to return a random listing from the database of listings we have
func get_random_listing(c *gin.Context) {

	rand.New(rand.NewSource(time.Now().Local().Unix()))
	randomidx := rand.Intn(len(all_listings))

	fmt.Printf("Random idx is %d\n", randomidx)
	c.JSON(http.StatusOK, all_listings[randomidx]) // this is not working
}

// ------------------------------------------------------
// MAIN

// GLOBAL VARS
var all_listings []Listing

func main() {
	fmt.Println("Welcome to realista!")

	// all_listings = scrape()
	// print_listing(all_listings[100])

	all_listings = csv2listings("data/current.csv")

	// Initialize APIs
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/rand_house", get_random_listing)

	// run API
	fmt.Println("Starting API at. Access it with  curl http://localhost:8080/rand_house")
	router.Run("localhost:8080")
}
