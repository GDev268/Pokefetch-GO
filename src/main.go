package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const OFFSET int = -3
const SHINY_TEXT string = "\u2605 Shiny! \u2605"
const SHINY_PROBABILITY int32 = 1

func main() {
	var isShiny = RandIntRange32(1, 4) == SHINY_PROBABILITY

	var random_int = RandIntRange32(1, 904)
	var url = fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", random_int)

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var pokemon Pokemon

	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		panic(err)
	}

	var pokemonName = StripPokemonForm(pokemon.Name)

	var display = FormatPokemonDisplay(pokemonName, ExtractTypes(&pokemon), isShiny)

	home, homeErr := os.UserHomeDir()
	if err != nil {
		panic(homeErr)
	}

	cached := filepath.Join(home, ".cache", "pokemon.txt")
	ffCfg := filepath.Join(home, ".config", "fastfetch", "config.jsonc")

	var shinyArg string = ""

	if isShiny {
		shinyArg = "-s"
	}

	cmdStr := fmt.Sprintf(
		"pokemon-colorscripts -n %s %s --no-title > %s",
		pokemonName,
		shinyArg,
		cached,
	)

	cmd := exec.Command("sh", "-c", cmdStr)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmdErr1 := cmd.Run()
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, "Failed to generate Pokémon ASCII!")
		if err != nil {
			return
		}
		panic(cmdErr1)
	}

	var fetch, fetchErr = NewPokeFetch(cached, ffCfg)

	if fetchErr != nil {
		panic(fetchErr)
	}

	runErr := fetch.Run(display)

	if runErr != nil {
		panic(runErr)
	}

	cmd2 := exec.Command("fastfetch")

	// Attach output streams to the terminal
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	cmdErr2 := cmd2.Run()
	if cmdErr2 != nil {
		fmt.Fprintln(os.Stderr, "Failed to run fastfetch!")
		panic(cmdErr2)
	}
}
