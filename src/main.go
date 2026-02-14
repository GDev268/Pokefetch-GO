package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var pokemon Pokemon

	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		panic(err)
	}

	var pokemonName = StripPokemonForm(ChangeInvalidNames(pokemon.ID, pokemon))

	var display = FormatPokemonDisplay(pokemonName, ExtractTypes(&pokemon), isShiny)

	home, _ := os.UserHomeDir()

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

	cmdErr := cmd.Run()
	if cmdErr != nil {
		return
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

	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	cmdErr2 := cmd2.Run()
	if cmdErr2 != nil {
		_, fprintErr := fmt.Fprintln(os.Stderr, "Failed to run fastfetch!")

		if fprintErr != nil {
			panic(fprintErr)
		}
		panic(cmdErr2)
	}
}
