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
	Next          *string
	Previous      *string
}

func main() {
	pokeClient := pokeapi.NewClient(5 * time.Second)
	config := &config{
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
		command, ok := getCommands()[commandName]
		if ok {
			err := command.callback(config)
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
