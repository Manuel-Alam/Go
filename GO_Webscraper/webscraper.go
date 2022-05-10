package main

import(
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
	"strings"
	"context"

	"github.com/gocolly/colly/v2"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/readpref"
)

//struct used to store json data for an item
type Item struct{
	Name string 	 `json:"name"`
	Colourway string `json:"style"` 
	Price string	 `json:"price"`
	ImgUrl string	 `json:"url"`
	ID string	     `json:"ID"`
	PairNumber int   `json:"PairNumber"`
}

func main(){
	
	//calling webscraper method
	webscraper()

	// Login to MongoDB account 
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://MannySchool:letrollMongoDB1!@cluster0.9buev.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil{
		log.Fatal(err)
	}

	// Declare Context type object for managing multiple API requests
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to the MongoDB and return Client instance
	err = client.Connect(ctx)
	if err !=nil{
		log.Fatal(err)
	}

	// Access a MongoDB collection through a database
	col := client.Database("Footlocker-DB").Collection("Nike-Shoes")

	// Load values from JSON file to model
	byteValues, err := ioutil.ReadFile("shoes.json")

	// Declare an empty slice for the items
	var items []Item

	// Unmarshal the encoded JSON byte string into the slice
	err = json.Unmarshal(byteValues, &items)

	// Iterate the slice of MongoDB struct docs
	for i := range items {

		// Put the document element in a new variable
		item := items[i]

		// Call the InsertOne() method and pass the context and doc objects
		col.InsertOne(ctx, item)
	}

}

//webscraper method
func webscraper(){

	//declare which domains are allowed to be webscraped
	c:= colly.NewCollector(
		colly.AllowedDomains("www.footlocker.ca"),
	)

	//empty slice of items.
	shoes := []Item{}

	//keep track of number of products.
	num := 1

	//each div with this element will create an item object and add it to the slice.
	c.OnHTML("div.ProductCard", func(element *colly.HTMLElement){

		product := element.DOM
		
		// This will find the link for each product and later store the url.
		str := element.ChildAttr("a","href")
		split := strings.Split(str,"/")
		shoeID := strings.ReplaceAll(split[4],".html","")

		// This will find the style for each item.
		style := product.Find("span.ProductName-alt").Text()
		split2 := strings.ReplaceAll(style,"Men's•","")

		//creating object here to store all the product data.
		shoe := Item{
			Name: product.Find("span.ProductName-primary").Text(),
			Colourway: split2,
			Price: product.Find("span.ProductPrice").Text(),
			ImgUrl: "footlocker.ca"+element.ChildAttr("a","href"),
			ID: shoeID,
			PairNumber: num,

		}

		//adding the new item to the slice.
		shoes = append(shoes, shoe)

		//incrementing number of items found.
		num++
		
	})
	
	//this is the specific webpage being scraped.
	c.Visit("https://www.footlocker.ca/en/category/brands/adidas.html?query=adidas%3AnewArrivals%3Agender%3AMen%27s%3Abrand%3Aadidas&sort=newArrivals")

	//writing the array of objects to json file.
	writeJSON(shoes)

	//print success if all done correctly.
	fmt.Println("SUCCESS")
}

//method to write the slice data to a json file.
func writeJSON(data []Item){
	f, err:= json.MarshalIndent(data, "", " ")

	if err != nil{
		log.Fatal(err)
		return 
	}

	_ = ioutil.WriteFile("shoes.json",f,0644)
}