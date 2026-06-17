package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type FastxMeta struct {
	Version    string            `toml:"version"`
	Registry   string            `toml:"registry"`
	Skipped    []string          `toml:"skipped"`
	Components map[string]string `toml:"components"`
}

func ParsePyproject(path string) (*FastxMeta, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("fail to read file: %w", err)
	}

	var wrapper struct {
		Tool struct {
			Fastx *FastxMeta `toml:"fastx"`
		} `toml:"tool"`
	}

	if _, err := toml.Decode(string(file), &wrapper); err != nil {
		return nil, fmt.Errorf("fail to decode toml: %w", err)
	}

	return wrapper.Tool.Fastx, nil
}

func WriteFastxMeta(path string, meta *FastxMeta) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("fail to read file: %w", err)
	}

	var doc map[string]any

	if _, err := toml.Decode(string(file), &doc); err != nil {
		return fmt.Errorf("fail to decode toml: %w", err)
	}

	if _, ok := doc["tool"]; !ok {
		doc["tool"] = make(map[string]any)
	}

	toolMap, ok := doc["tool"].(map[string]any)
	if !ok {
		return fmt.Errorf("tool is invalid format")
	}

	toolMap["fastx"] = meta

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("fail to create file: %w", err)
	}
	defer f.Close()

	err= toml.NewEncoder(f).Encode(&doc)

	return nil
}
