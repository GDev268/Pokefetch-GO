package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// The Struct should be the Noun
type PokeFetch struct {
	cachedPath string
	configPath string
	colorFmt   string
	ffLines    uint
}

// The Constructor should be the Verb
func NewPokeFetch(cachedPath, configPath string) (*PokeFetch, error) {
	// ExtractColors is assumed to be defined in your other file (colors.go)
	color, lines, err := ExtractColors(cachedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract colors: %w", err)
	}

	return &PokeFetch{
		cachedPath: cachedPath,
		configPath: configPath,
		colorFmt:   color,
		ffLines:    lines,
	}, nil
}

func PaddedKey(icon string, label string) string {
	return fmt.Sprintf("│ %s %-8s", icon, label)
}

func (pf *PokeFetch) Run(pokemonDisplay string) error {
	config := make(map[string]interface{})

	// 1. DYNAMIC VERTICAL CENTERING
	const moduleHeight = 13
	var paddingTop uint = 1
	if pf.ffLines > moduleHeight {
		paddingTop = (pf.ffLines - moduleHeight) / 2
	}

	config["display"] = map[string]interface{}{
		"separator": " │ ",
		"color": map[string]interface{}{
			"keys": "default",
		},
	}

	config["logo"] = map[string]interface{}{
		"type":   "command-raw",
		"source": fmt.Sprintf("cat %s", pf.cachedPath),
		"padding": map[string]interface{}{
			"top":  paddingTop,
			"left": 2,
		},
	}

	var modules []interface{}

	// 2. THE BOX TOP
	modules = append(modules, map[string]interface{}{
		"type":   "custom",
		"format": "╭────────────╮",
	})

	// 3. SYSTEM INFO
	// PaddedKey provides the "│" border on the far left
	modules = append(modules, map[string]interface{}{"type": "title", "key": PaddedKey("", "user")})
	modules = append(modules, map[string]interface{}{"type": "os", "key": PaddedKey("󰋄", "distro")})
	modules = append(modules, map[string]interface{}{"type": "kernel", "key": PaddedKey("⚙", "kernel")})
	modules = append(modules, map[string]interface{}{"type": "uptime", "key": PaddedKey("󱎫", "uptime")})
	modules = append(modules, map[string]interface{}{"type": "shell", "key": PaddedKey("󱆃", "shell")})
	modules = append(modules, map[string]interface{}{"type": "terminal", "key": PaddedKey("", "term")})

	modules = append(modules, map[string]interface{}{
		"type":   "cpu",
		"key":    PaddedKey("", "cpu"),
		"format": "{1}",
	})
	modules = append(modules, map[string]interface{}{
		"type":   "gpu",
		"key":    PaddedKey("󰢮", "gpu"),
		"format": "{2}",
	})

	modules = append(modules, map[string]interface{}{"type": "memory", "key": PaddedKey("󰍛", "memory")})
	modules = append(modules, map[string]interface{}{"type": "disk", "key": PaddedKey("󰋊", "disk")})

	// 4. POKEMON INFO
	modules = append(modules, map[string]interface{}{
		"type":   "custom",
		"key":    PaddedKey("󰐿", "poke"),
		"format": pokemonDisplay,
	})

	// 5. COLORS
	modules = append(modules, map[string]interface{}{
		"type":   "colors",
		"key":    PaddedKey("󰈊", "colors"),
		"symbol": "circle",
	})

	// 6. THE BOX BOTTOM
	modules = append(modules, map[string]interface{}{
		"type":   "custom",
		"format": "╰────────────╯",
	})

	config["modules"] = modules

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pf.configPath, data, 0644)
}

func loadConfig(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return make(map[string]interface{}), nil
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result, nil
}
