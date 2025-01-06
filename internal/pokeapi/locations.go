package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gravitonsmith/pokedexcli/internal/pokecache"
)

const baseUrl = "https://pokeapi.co/api/v2"

type Client struct {
	cache      pokecache.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cacheInterval),
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

	if val, ok := c.cache.Get(url); ok {
		locations := BatchLocation{}
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return BatchLocation{}, err
		}
		return locations, nil
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

	c.cache.Add(url, jsonData)

	locations := BatchLocation{}
	err = json.Unmarshal(jsonData, &locations)
	if err != nil {
		return BatchLocation{}, nil
	}

	return locations, nil
}

func (c *Client) ExploreLocation(location string) (LocationInfo, error) {
	url := baseUrl + "/location-area/" + location

	if val, ok := c.cache.Get(url); ok {
		locationInfo := LocationInfo{}
		err := json.Unmarshal(val, &locationInfo)
		if err != nil {
			return LocationInfo{}, err
		}
		return locationInfo, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationInfo{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationInfo{}, err
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationInfo{}, err
	}

	c.cache.Add(url, jsonData)

	locationInfo := LocationInfo{}
	err = json.Unmarshal(jsonData, &locationInfo)
	if err != nil {
		return LocationInfo{}, nil
	}
	return locationInfo, nil
}
