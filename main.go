package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gravitonsmith/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiClient pokeapi.Client
	pokedex       map[string]pokeapi.PokemonInfo
	Next          *string
	Previous      *string
}

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	config := &config{
		pokedex:       map[string]pokeapi.PokemonInfo{},
		pokeapiClient: pokeClient,
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		commandName := input[0]
		args := []string{}
		if len(input) > 1 {
			args = input[1:]
		}
		command, ok := getCommands()[commandName]
		if ok {
			err := command.callback(config, args...)
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
