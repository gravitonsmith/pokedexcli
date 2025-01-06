package main

import (
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Shows help text",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "List next locations from pokemom map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List previous locations from pokemom map",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location from the map",
			callback:    commandExplore,
		},
	}
}

func commandExplore(c *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("You must provide a location")
	}

	location := args[0]
	locationData, err := c.pokeapiClient.ExploreLocation(location)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s ...\n", locationData.Name)
	fmt.Println("Pokemon in this location include:")
	for _, pokemon := range locationData.PokemonEncounters {
		fmt.Printf("-- %v\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandExit(c *config, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMapb(c *config, args ...string) error {
	jsonData, err := c.pokeapiClient.ListLocations(c.Previous)
	if err != nil {
		return err
	}

	c.Next = jsonData.Next
	c.Previous = jsonData.Previous

	for _, place := range jsonData.Results {
		fmt.Println(place.Name)
	}

	return nil
}

func commandMap(c *config, args ...string) error {
	jsonData, err := c.pokeapiClient.ListLocations(c.Next)
	if err != nil {
		return err
	}

	c.Next = jsonData.Next
	c.Previous = jsonData.Previous

	for _, place := range jsonData.Results {
		fmt.Println(place.Name)
	}

	return nil
}
