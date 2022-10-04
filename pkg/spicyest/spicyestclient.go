package spicyest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type GetPricesResponse struct {
	NextPage string     `json:"next_page"`
	Prices   []NftPrice `json:"prices"`
}

type NftPrice struct {
	Address      string    `json:"address"`
	TokenID      string    `json:"tokenID"`
	CalculatedAt time.Time `json:"calculatedAt"`
	Currency     string    `json:"currency"`
	Price        float64   `json:"price"`
	PriceLower   float64   `json:"price_lower"`
	PriceUpper   float64   `json:"price_upper"`
}
type SpicyestAPI interface {
	GetPrices(address string, tokenIds []string, limit int32, nextPage string) (*GetPricesResponse, error)
}

type httpClient struct {
	client   *http.Client
	endpoint string
	apiKey   string
}

func NewHttpClient(client *http.Client, apiKey string) SpicyestAPI {
	return &httpClient{
		endpoint: "http://api.spicyest.com/",
		client:   client,
		apiKey:   apiKey,
	}
}

func (c *httpClient) GetPrices(address string, tokenIds []string, limit int32, nextPage string) (*GetPricesResponse, error) {
	url := fmt.Sprintf("%s/prices?address=%s&limit=%d", c.endpoint, address, limit)
	if nextPage != "" {
		url += "&next_page=" + nextPage
	}

	if len(tokenIds) != 0 {
		tokenIdsStr := strings.Join(tokenIds, "%2C")
		url += "&tokenIDs=" + tokenIdsStr
	}
	fmt.Printf("This is the url:% v\n", url)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "*/*")
	req.Header.Add("X-API-Key", c.apiKey)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error calling spicyest API")
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected API to not error. Got: %s", res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, errors.Wrap(err, "error reading response")
	}

	priceRes := &GetPricesResponse{}
	err = json.Unmarshal(body, &priceRes)

	if err != nil {
		return nil, errors.Wrap(err, "error unmmarhsalling response from API")
	}
	return priceRes, nil
}
