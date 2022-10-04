package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/nevadooo/improved-waddle/pkg/spicyest"
)

var (
	address     *string
	spicyAPIKey *string
)

func init() {
	address = flag.String("address", "", "collection address")
	spicyAPIKey = flag.String("spicy-api-key", "", "api key for spicy")
}

func main() {
	flag.Parse()

	if *address == "" {
		fmt.Println("address is required")
		return
	}

	c := &http.Client{}
	client := spicyest.NewHttpClient(c, *spicyAPIKey)

	prices, err := client.GetPrices(*address, []string{"9181"}, 1, "")
	if err != nil {
		fmt.Printf("Found error: %v\n", err)
		return
	}
	fmt.Printf("This is the response: %v\n", prices)
}
