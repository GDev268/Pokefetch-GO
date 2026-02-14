package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const COLOR_STEP = 8

func ansiBG(color uint8) string {
	return fmt.Sprintf("\x1b[48;5;%dm", color)
}

// ANSI foreground
func ansiFG(color uint8) string {
	return fmt.Sprintf("\x1b[38;5;%dm", color)
}

// ANSI reset
func ansiReset() string {
	return "\x1b[0m"
}

func RandIntRange32(min, max int32) int32 {
	return rand.Int32N(max-min+1) + min
}

func ExtractTypes(pokemon *Pokemon) []string {
	var types []string

	typesVal, ok := pokemon.Types.([]interface{})

	if !ok {
		return types
	}

	for _, t := range typesVal {
		entryMap, ok := t.(map[string]interface{})

		if !ok {
			continue
		}

		typeMap, ok := entryMap["type"].(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := typeMap["name"].(string)

		if ok {
			types = append(types, name)
		}

	}

	return types
}

func QuantizeColor(r, g, b uint8, step uint8) (uint8, uint8, uint8) {
	return (r / step) * step, (g / step) * step, (b / step) * step
}

func ExtractColors(path string) (string, uint, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return "", 0, err
	}

	pokemon := string(data)

	var pokemonLines = len(strings.Split(strings.ReplaceAll(pokemon, "\r\n", "\n"), "\n"))

	var ffLines = math.Max(float64(pokemonLines+OFFSET), 0)

	re, err := regexp.Compile(`(?:38|48);2;(\d{1,3});(\d{1,3});(\d{1,3})`)
	if err != nil {
		return "", 0, err // handle the error
	}

	matches := re.FindAllStringSubmatch(pokemon, -1)

	type RGB struct {
		R, G, B uint8
	}

	var counts = make(map[RGB]uint)

	for _, match := range matches {
		var r, _ = strconv.Atoi(match[1])
		var g, _ = strconv.Atoi(match[2])
		var b, _ = strconv.Atoi(match[3])

		var dark = r < 90 && g < 90 && b < 90
		var light = r > 180 && g > 180 && b > 180

		if dark || light {
			continue
		}

		resultR, resultG, resultB := QuantizeColor(uint8(r), uint8(g), uint8(b), COLOR_STEP)

		var rgb = RGB{resultR, resultG, resultB}

		counts[rgb]++
	}

	var maxCount uint = 0
	var index *RGB

	for key, count := range counts {
		if maxCount < count {
			maxCount = count
			index = &key
		}
	}

	if index == nil {
		return "", 0, errors.New("no colors found")
	}

	return fmt.Sprintf("38;2;%d;%d;%d", index.R, index.G, index.B), uint(ffLines), nil
}

func Module(t, key, color string) interface{} {
	return map[string]interface{}{
		"type":       t,
		"key":        key,
		"keyColor":   color,
		"valueColor": color,
	}
}

func ForegroundForBG(bg uint8) uint8 {
	switch bg {
	case 1 | 5 | 8 | 21 | 99 | 236 | 55:
		return 255
	default:
		return 232
	}
}

func PokemonTypeColor(pokemonType string) uint8 {
	switch pokemonType {
	case "normal":
		return 101
	case "fire":
		return 202
	case "water":
		return 31
	case "electric":
		return 226
	case "grass":
		return 76
	case "ice":
		return 81
	case "fighting":
		return 124
	case "poison":
		return 127
	case "ground":
		return 178
	case "flying":
		return 98
	case "psychic":
		return 170
	case "bug":
		return 142
	case "rock":
		return 101
	case "ghost":
		return 55
	case "dragon":
		return 21
	case "dark":
		return 236
	case "steel":
		return 247
	case "fairy":
		return 219
	default:
		return 0
	}
}

func CreateTextBadge(text string, bgColor uint8, bold bool) string {
	var fgColor = ForegroundForBG(bgColor)

	var fg = ansiFG(fgColor)
	var bg = ansiBG(bgColor)
	var boldCode string

	if bold {
		boldCode = "\x1b[1m"
	} else {
		boldCode = ""
	}

	return fmt.Sprintf("%s%s%s %s %s%s", boldCode, fg, bg, text, ansiReset(), ansiReset())
}

func GetTypeBadges(types []string) string {
	var result string = ""
	for _, curType := range types {
		var color = PokemonTypeColor(curType)

		result += CreateTextBadge(strings.ToUpper(curType), color, false) + " "
	}

	return result
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

const POKEMON_NAME_COLOR = 15
const POKEMON_SHINY_COLOR = 220

func FormatPokemonDisplay(name string, types []string, isShiny bool) string {
	var capitalizedName = Capitalize(name)

	var nameBadge = CreateTextBadge(capitalizedName, POKEMON_NAME_COLOR, true)

	var shinyBadge string

	if isShiny {
		shinyBadge = CreateTextBadge(SHINY_TEXT, POKEMON_SHINY_COLOR, true)
	} else {
		shinyBadge = ""
	}

	var typeBadges = GetTypeBadges(types)

	var pokemonInfo = fmt.Sprintf("%s %s %s", nameBadge, shinyBadge, typeBadges)

	return pokemonInfo
}

func StripPokemonForm(name string) string {
	if i := strings.IndexByte(name, '-'); i >= 0 {
		return name[:i]
	}
	return name
}

func ChangeInvalidNames(pokemonID int, pokemon Pokemon) string {
	switch pokemonID {
	case 29:
		return "nidoran-f"
	case 32:
		return "nidoran-m"
	case 122:
		return "mr-mime"
	case 386:
		return "deoxys"
	case 413:
		return "wormadam"
	case 487:
		return "giratina"
	case 492:
		return "shaymin"
	case 550:
		return "basculin"
	case 555:
		return "darmanitan"
	case 641:
		return "tornadus"
	case 642:
		return "thundurus"
	case 645:
		return "landorus"
	case 647:
		return "keldeo"
	case 648:
		return "meloetta"
	case 678:
		return "meowstic"
	case 681:
		return "aegislash"
	case 710:
		return "pumpkaboo"
	case 711:
		return "gourgeist"
	case 718:
		return "zygarde"
	case 741:
		return "oricorio"
	case 745:
		return "lycanroc"
	case 746:
		return "wishiwashi"
	case 774:
		return "minior"
	case 778:
		return "mimikyu"
	case 849:
		return "toxtricity"
	case 875:
		return "eiscue"
	case 876:
		return "indeedee"
	case 877:
		return "morpeko"
	case 892:
		return "urshifu"
	case 902:
		return "basculegion"
	default:
		return pokemon.Name
	}
}
