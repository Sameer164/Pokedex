package main

import (
	"Pokedex/internal/pokecache"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const Interval = 5

var cache = pokecache.NewCache(Interval)

type config struct {
	Next     string
	Previous string
}

type commands struct {
	name        string
	description string
	callback    func(*config, []string) error
}

var supportedCommands map[string]commands

func commandExit(c *config, args []string) (err error) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return
}

func commandHelp(c *config, args []string) (err error) {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
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

func commandMap(c *config, args []string) (err error) {
	var url string
	if c.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = c.Next
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("There was an error in fetching the Pokemon Location Areas.")
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("There was an error in reading the Pokemon Location Areas.")
	}
	cache.Set(url, data)
	var httpResponse map[string]interface{}

	err = json.Unmarshal(data, &httpResponse)
	if err != nil {
		return fmt.Errorf("There was an error in unmarshaling the data.")
	}

	if httpResponse["next"] == nil {
		c.Next = ""
	} else {
		c.Next = httpResponse["next"].(string)
	}
	if httpResponse["previous"] == nil {
		c.Previous = ""
	} else {
		c.Previous = httpResponse["previous"].(string)
	}

	results, ok := httpResponse["results"].([]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for results")
	}

	for _, m := range results {
		m2, ok := m.(map[string]interface{})
		if !ok {
			return fmt.Errorf("There was an error in reading the Pokemon Location Areas.")
		}
		fmt.Println(m2["name"])
	}
	return
}

func commandMapb(c *config, args []string) (err error) {
	var url string
	if c.Previous == "" {
		return fmt.Errorf("you're on the first page")
	} else {
		url = c.Previous
	}

	val, ok := cache.Get(url)
	var data []byte
	if ok {
		data = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("There was an error in fetching the Pokemon Location Areas.")
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("There was an error in reading the Pokemon Location Areas.")
		}
	}
	var httpResponse map[string]interface{}
	err = json.Unmarshal(data, &httpResponse)
	if err != nil {
		return fmt.Errorf("There was an error in unmarshaling the data.")
	}

	if httpResponse["next"] == nil {
		c.Next = ""
	} else {
		c.Next = httpResponse["next"].(string)
	}
	if httpResponse["previous"] == nil {
		c.Previous = ""
	} else {
		c.Previous = httpResponse["previous"].(string)
	}

	results, ok := httpResponse["results"].([]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for results")
	}

	for _, m := range results {
		m2, ok := m.(map[string]interface{})
		if !ok {
			return fmt.Errorf("There was an error in reading the Pokemon Location Areas.")
		}
		fmt.Println(m2["name"])
	}
	return
}

func commandExplore(c *config, args []string) error {
	fmt.Printf("Exploring %s...\n", args[0])
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]
	val, exists := cache.Get(url)
	var data []byte
	if exists {
		data = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("There was an error in fetching the Pokemons.")
		}
		defer res.Body.Close()
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("There was an error in reading the Pokemons")
		}
		cache.Set(url, data)
	}
	var httpResponse map[string]interface{}
	err := json.Unmarshal(data, &httpResponse)
	if err != nil {
		return fmt.Errorf("There was an error in unmarshaling the data.")
	}

	pokemons, ok := httpResponse["pokemon_encounters"].([]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for pokemon_encounters")
	}

	for _, m := range pokemons {
		m2, ok := m.(map[string]interface{})
		if !ok {
			return fmt.Errorf("There was an error in reading the Pokemons.")
		}
		n, ok := m2["pokemon"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("There was an error in reading the Pokemons.")
		}
		fmt.Println(n["name"])
	}
	return nil
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
		"map": {
			name:        "map",
			description: "Displays the names of next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays the names of all the pokemons in a certain location",
			callback:    commandExplore,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	mapConfig := config{}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)
		fmt.Printf("%v\n", cleaned)
		command, exists := supportedCommands[strings.ToLower(cleaned[0])]
		if !exists {
			fmt.Println("Unknown command")
		} else {
			err := command.callback(&mapConfig, cleaned[1:])
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}
}
