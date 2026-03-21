package config

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ThisaruGuruge/bestow/internal/log"
)

//go:embed defaults/default-config.yaml
var defaultConfigTemplate string

var DefaultIgnoreList []string = []string{".git", ".gitignore", "README.md", "LICENSE"}

func GetDefaultConfigTemplate(source, destination string) (string, error) {
	tmpl, err := template.New("config").Parse(defaultConfigTemplate)
	if err != nil {
		return "", err
	}
	data := struct {
		Source      string
		Destination string
	}{
		Source:      source,
		Destination: destination,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func setDefaultSource(config *Config) error {
	log.Debug("checking source config")
	if config.Source != "" {
		log.Debug("source is set by configs", "source", config.Source)
		return nil
	}
	log.Debug("no source provided, setting default source")
	home, err := os.UserHomeDir()
	if err != nil {
		return &ConfigError{
			Message:    "failed to read home directory",
			ConfigName: "source",
			Value:      "$HOME",
			Cause:      err,
		}
	}
	config.Source = filepath.Join(home, "dotfiles")
	log.Debug("default value is set for source", "source", config.Source)
	return nil
}

func setDefaultDestination(config *Config) error {
	log.Debug("checking destination config")
	if config.Destination != "" {
		log.Debug("destination is set by configs", "destination", config.Destination)
		return nil
	}
	log.Debug("no destination provided, setting default destination")
	home, err := os.UserHomeDir()
	if err != nil {
		return &ConfigError{
			Message:    "failed to load default config",
			ConfigName: "destination",
			Value:      "user home directory",
			Cause:      err,
		}
	}
	config.Destination = home
	log.Debug("default value is set for destination", "destination", config.Destination)
	return nil
}
