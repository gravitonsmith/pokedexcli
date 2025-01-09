package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) CatchPoke(name string) (PokemonInfo, error) {
	url := baseUrl + "/pokemon/" + name

	if val, ok := c.cache.Get(url); ok {
		pokemonData := PokemonInfo{}
		err := json.Unmarshal(val, &pokemonData)
		if err != nil {
			return PokemonInfo{}, err
		}
		return pokemonData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokemonInfo{}, nil
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return PokemonInfo{}, err
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonInfo{}, err
	}

	c.cache.Add(url, jsonData)

	pokemonData := PokemonInfo{}
	err = json.Unmarshal(jsonData, &pokemonData)
	if err != nil {
		return PokemonInfo{}, nil
	}

	return pokemonData, nil
}
