package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandExit(c *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMapb(c *config) error {
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

func commandMap(c *config) error {
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
	}
}
