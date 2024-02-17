package application

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mcuadros/go-defaults"
	"github.com/shibukawa/configdir"
	"gopkg.in/yaml.v2"
)

const (
	vendorName       = "aero"
	appName          = "adsb1090"
	settingsFilename = "settings.yaml"
)

type settingsContext struct{}

// Settings is a configuration file.
type Settings struct {
	basePath                 string
	AircraftDatabaseFilename string `default:"aircrafts.json.gz" json:"acDbFilename" yaml:"acDbFilename"`
}

// NewSettings creates a new configuration struct.
func NewSettings() *Settings {
	output := &Settings{}
	defaults.SetDefaults(output)

	return output
}

// UserSettings reads the user configuration.
func UserSettings() (*Settings, error) {
	config := NewSettings()

	configDirs := configdir.New(vendorName, appName)

	configDirs.LocalPath, _ = filepath.Abs(".")

	folder := configDirs.QueryFolderContainsFile(settingsFilename)

	if folder != nil {
		config.basePath = folder.Path

		data, _ := folder.ReadFile(settingsFilename)

		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}

		return config, nil
	}

	// Stores to user folder
	folders := configDirs.QueryFolders(configdir.Global)
	if len(folders) > 0 {
		config.basePath = folders[0].Path

		data, err := yaml.Marshal(config)
		if err != nil {
			return nil, err
		}

		if err := folders[0].WriteFile(settingsFilename, data); err != nil {
			return nil, err
		}

		return config, nil
	}

	return nil, os.ErrNotExist
}

// AircraftDatabaseFile is the path of the aircraft database.
func (s Settings) AircraftDatabaseFile() string {
	return filepath.Join(s.basePath, s.AircraftDatabaseFilename)
}

// WithSettings sets the settings in the context.
func (s *Settings) WithSettings(ctx context.Context) context.Context {
	return context.WithValue(ctx, settingsContext{}, s)
}

// SettingsFromContext gets the settings from the context.
func SettingsFromContext(ctx context.Context) (*Settings, error) {
	if settings, found := ctx.Value(settingsContext{}).(*Settings); found {
		return settings, nil
	}

	return nil, os.ErrNotExist
}
