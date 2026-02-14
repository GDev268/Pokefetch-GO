package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type PokeFetch struct {
	cachedPath string
	configPath string
	colorFmt   string
	ffLines    uint
}

func NewPokeFetch(cachedPath, configPath string) (*PokeFetch, error) {
	if color, lines, err := ExtractColors(cachedPath); err == nil {
		return &PokeFetch{
			cachedPath: cachedPath,
			configPath: configPath,
			colorFmt:   color,
			ffLines:    lines,
		}, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Failed to create PokeFetch!: %s", err.Error()))
	}
}

func (pf *PokeFetch) Run(pokemonDisplay string) error {
	var config, err = loadConfig(pf.configPath)

	if err != nil {
		return err
	}

	config["logo"] = map[string]interface{}{
		"type":   "command-raw",
		"source": fmt.Sprintf("cat %s", pf.cachedPath),
		"padding": map[string]interface{}{
			"top": 2,
		},
	}

	var display = config["display"].(map[string]interface{})

	if display == nil {
		display = make(map[string]interface{})
		config["display"] = display
	}

	var color = display["color"].(map[string]interface{})

	if color == nil {
		color = make(map[string]interface{})
	}

	color["title"] = pf.colorFmt
	color["keys"] = pf.colorFmt

	var modules = []interface{}{}

	modules = append(modules, "title")
	modules = append(modules, "separator")

	modules = append(modules, Module("os", "os    ", pf.colorFmt))
	modules = append(modules, Module("kernel", "kernel", pf.colorFmt))
	modules = append(modules, Module("uptime", "uptime", pf.colorFmt))
	modules = append(modules, Module("processes", "proc  ", pf.colorFmt))
	modules = append(modules, Module("packages", "pkgs  ", pf.colorFmt))
	modules = append(modules, Module("shell", "shell ", pf.colorFmt))
	modules = append(modules, Module("monitor", "mon   ", pf.colorFmt))
	modules = append(modules, Module("terminal", "term  ", pf.colorFmt))

	modules = append(modules, map[string]interface{}{
		"type":            "cpu",
		"key":             "cpu   ",
		"keyColor":        pf.colorFmt,
		"valueColor":      pf.colorFmt,
		"showPeCoreCount": false,
		"temp":            true,
	})

	modules = append(modules, map[string]interface{}{
		"type":           "gpu",
		"key":            "gpu   ",
		"keyColor":       pf.colorFmt,
		"valueColor":     pf.colorFmt,
		"driverSpecific": true,
		"temp":           true,
	})

	modules = append(modules, Module("memory", "memory", pf.colorFmt))
	modules = append(modules, Module("disk", "disk  ", pf.colorFmt))
	modules = append(modules, Module("media", "media ", pf.colorFmt))
	modules = append(modules, Module("datetime", "time ", pf.colorFmt))
	modules = append(modules, Module("version", "ver   ", pf.colorFmt))
	modules = append(modules, "separator")

	modules = append(modules, map[string]interface{}{
		"type":       "custom",
		"key":        "pokemon",
		"format":     pokemonDisplay,
		"keyColor":   pf.colorFmt,
		"valueColor": pf.colorFmt,
	})

	modules = append(modules, "break")
	modules = append(modules, "colors")

	config["modules"] = modules

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(pf.configPath, data, 0644); err != nil {
		return err
	}

	return nil
}

func loadConfig(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var resultRaw interface{}
	if err := json.Unmarshal(data, &resultRaw); err != nil {
		return nil, err
	}

	result, ok := resultRaw.(map[string]interface{})
	if !ok {
		return nil, errors.New("expected config to be a JSON object")
	}

	return result, nil
}
