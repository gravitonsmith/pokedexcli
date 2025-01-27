package main

import (
	"errors"
	"fmt"
	"math/rand"
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
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon and add it to the Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon that you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Check your pokedex!",
			callback:    commandPokedex,
		},
	}
}

func commandPokedex(c *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range c.pokedex {
		fmt.Printf("-%v\n", pokemon.Name)
	}
	return nil
}

func commandInspect(c *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("Please provide a pokemon name")
	}

	name := args[0]
	pokemon, ok := c.pokedex[name]
	if !ok {
		return errors.New("you have not caught that pokemon")
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("-%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Type:")
	for _, typ := range pokemon.Types {
		fmt.Printf("-%v\n", typ.Type.Name)
	}
	return nil
}

func commandCatch(c *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("You must specify a pokemon to try and catch")
	}

	pokemon := args[0]
	pokemonData, err := c.pokeapiClient.CatchPoke(pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonData.Name)
	if rand.Intn(pokemonData.BaseExperience) > 70 {
		fmt.Printf("%v escaped...\n", pokemonData.Name)
		return nil
	}

	fmt.Printf("%v was caught!\nAdding to pokedex!\n", pokemonData.Name)
	c.pokedex[pokemonData.Name] = pokemonData

	return nil
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
