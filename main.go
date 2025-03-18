package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type commands struct {
	name        string
	description string
	callback    func() error
}

var supportedCommands map[string]commands

func commandExit() (err error) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return
}

func commandHelp() (err error) {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for k, v := range supportedCommands {
		fmt.Printf("%s: %s\n", k, v.description)
	}
	return
}

func cleanInput(text string) []string {
	splitted := strings.Split(strings.Trim(text, " "), " ")
	filtered := make([]string, 0)
	for _, word := range splitted {
		if word != "" {
			filtered = append(filtered, word)
		}
	}
	return filtered

}

func init() {
	supportedCommands = map[string]commands{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)
		command, exists := supportedCommands[strings.ToLower(cleaned[0])]
		if !exists {
			fmt.Println("Unknown command")
		} else {
			err := command.callback()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}
}
