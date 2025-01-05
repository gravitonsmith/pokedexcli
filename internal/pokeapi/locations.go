package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const baseUrl = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) ListLocations(pageUrl *string) (BatchLocation, error) {
	url := baseUrl + "/location-area"
	if pageUrl != nil {
		url = *pageUrl
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BatchLocation{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return BatchLocation{}, err
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return BatchLocation{}, err
	}

	locations := BatchLocation{}
	err = json.Unmarshal(jsonData, &locations)
	if err != nil {
		return BatchLocation{}, nil
	}

	return locations, nil
}
