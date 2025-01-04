package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type config struct {
	Next     string
	Previous string
}

func commandExit(c *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range getCommands(c) {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMapb(c *config) error {
	res, err := http.Get(c.Previous)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	jsonData := batchLocation{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return err
	}

	if jsonData.Next != nil {
		c.Next = *jsonData.Next
	} else {
		c.Next = ""
	}
	if jsonData.Previous != nil {
		c.Previous = *jsonData.Previous
	} else {
		c.Previous = ""
	}

	for _, place := range jsonData.Results {
		fmt.Println(place.Name)
	}

	return nil
}

func commandMap(c *config) error {
	res, err := http.Get(c.Next)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	jsonData := batchLocation{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return err
	}

	if jsonData.Next != nil {
		c.Next = *jsonData.Next
	} else {
		c.Next = ""
	}
	if jsonData.Previous != nil {
		c.Previous = *jsonData.Previous
	} else {
		c.Previous = ""
	}

	for _, place := range jsonData.Results {
		fmt.Println(place.Name)
	}

	return nil
}

func getCommands(c *config) map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    func() error { return commandExit(c) },
		},
		"help": {
			name:        "help",
			description: "Shows help text",
			callback:    func() error { return commandHelp(c) },
		},
		"map": {
			name:        "map",
			description: "List next locations from pokemom map",
			callback:    func() error { return commandMap(c) },
		},
		"mapb": {
			name:        "mapb",
			description: "List previous locations from pokemom map",
			callback:    func() error { return commandMapb(c) },
		},
	}
}

func main() {
	config := config{Next: "https://pokeapi.co/api/v2/location-area/", Previous: ""}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		configPtr := &config
		commandName := input[0]
		command, ok := getCommands(configPtr)[commandName]
		if ok {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	cleaned := strings.Fields(strings.ToLower(text))
	return cleaned
}

type batchLocation struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
