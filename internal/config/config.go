// Package config holds the configuration.
package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/landru29/adsb1090/internal/serialize/nmea"
	"github.com/landru29/adsb1090/internal/transport/net"
	"github.com/mcuadros/go-defaults"
	"github.com/shibukawa/configdir"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

const (
	vendorName                            = "aero"
	appName                               = "adsb1090"
	settingsFilename                      = "settings.yaml"
	defaultNMEAmid                        = 226
	defaultFrequency                      = 1090000000
	defaultDatabaseLifetime time.Duration = time.Minute
)

// Config is the application configuration.
type Config struct {
	basePath string

	FixturesFilename         string             `default:""                                               json:"fixturesFilename"         yaml:"fixturesFilename"`         //nolint: lll
	FixtureLoop              bool               `default:"false"                                          json:"fixtureLoop"              yaml:"fixtureLoop"`              //nolint: lll
	DeviceIndex              uint32             `default:"0"                                              json:"deviceIndex"              yaml:"deviceIndex"`              //nolint: lll
	Frequency                uint32             `default:"1090000000"                                     json:"frequency"                yaml:"frequency"`                //nolint: lll
	Gain                     float64            `default:"0"                                              json:"gain"                     yaml:"gain"`                     //nolint: lll
	EnableAGC                bool               `default:"false"                                          json:"enableAgc"                yaml:"enableAgc"`                //nolint: lll
	DatabaseLifetime         time.Duration      `default:"0"                                              json:"databaseLifetime"         yaml:"databaseLifetime"`         //nolint: lll
	RefAircraftDatabaseURL   string             `default:"https://opensky-network.org/datasets/metadata/" json:"refAircraftDatabaseUrl"   yaml:"refAircraftDatabaseUrl"`   //nolint: lll
	UDPConf                  net.ProtocolConfig `default:""                                               json:"udpConf"                  yaml:"udpConf"`                  //nolint: lll
	TCPConf                  net.ProtocolConfig `default:""                                               json:"tcpConf"                  yaml:"tcpConf"`                  //nolint: lll
	HTTPConf                 HTTPConfig         `default:""                                               json:"httpConf"                 yaml:"httpConf"`                 //nolint: lll
	TransportScreen          string             `default:""                                               json:"transportScreen"          yaml:"transportScreen"`          //nolint: lll
	NmeaVessel               Vessel             `default:""                                               json:"nmeaVessel"               yaml:"nmeaVessel"`               //nolint: lll
	NmeaMid                  uint16             `default:"226"                                            json:"nmeaMid"                  yaml:"nmeaMid"`                  //nolint: lll
	TransportFile            string             `default:""                                               json:"transportFile"            yaml:"transportFile"`            //nolint: lll
	AircraftDatabaseFilename string             `default:"aircrafts.json.gz"                              json:"aircraftDatabaseFilename" yaml:"aircraftDatabaseFilename"` //nolint: lll
}

func newConfig(flags *pflag.FlagSet) *Config { //nolint: funlen
	output := &Config{
		DatabaseLifetime: defaultDatabaseLifetime,
		UDPConf:          net.NewProtocol("udp"),
		TCPConf:          net.NewProtocol("tcp"),
		NmeaVessel:       nmea.VesselTypeAircraft,
	}
	if flags != nil {
		flags.StringVarP(
			&output.RefAircraftDatabaseURL,
			"url",
			"u",
			"https://opensky-network.org/datasets/metadata/",
			"URL to aircraft database",
		)

		flags.StringVarP(
			&output.FixturesFilename,
			"fixture-file",
			"",
			"",
			"Filename of the fixture data file",
		)

		flags.Uint32VarP(
			&output.DeviceIndex,
			"device",
			"d",
			0,
			"Device index",
		)

		flags.BoolVarP(
			&output.EnableAGC,
			"enable-agc",
			"a",
			false,
			"Enable AGC",
		)

		flags.Uint32VarP(
			&output.Frequency,
			"frequency",
			"f",
			defaultFrequency,
			"frequency in Hz",
		)

		flags.DurationVarP(
			&output.DatabaseLifetime,
			"db-lifetime",
			"",
			defaultDatabaseLifetime,
			"lifetime of elements in the AC database",
		)

		flags.Float64VarP(
			&output.Gain,
			"gain",
			"g",
			0,
			"gain valid values are: 1.5, 4, 6.5, 9, 11.5, 14, 16.5, 19, 21.5, 24, 29, 34, 42, 43, 45, 47, 49",
		)

		flags.VarP(
			&output.UDPConf,
			"udp",
			"",
			"transmit data over udp (syntax: 'direction>format@host:port'; ie: --udp dial>json@192.168.1.10:8000)",
		)

		flags.VarP(
			&output.TCPConf,
			"tcp",
			"",
			"transmit data over tcp (syntax: 'direction>format@host:port'; ie: --tcp bind>json@192.168.1.10:8000)",
		)

		flags.VarP(
			&output.HTTPConf,
			"http",
			"",
			"transmit data over http (syntax: 'host:port/path'; ie: --http 0.0.0.0:8080/api)",
		)

		flags.StringVarP(
			&output.TransportScreen,
			"screen",
			"",
			"",
			"format to display output on the screen (json|nmea|text|none)",
		)

		flags.VarP(
			&output.NmeaVessel,
			"nmea-vessel",
			"",
			"MMSI vessel (aircraft|helicopter)",
		)

		flags.Uint16VarP(
			&output.NmeaMid,
			"nmea-mid",
			"",
			defaultNMEAmid,
			"MID (command 'mid' to list)",
		)

		flags.BoolVarP(
			&output.FixtureLoop,
			"loop",
			"",
			false,
			"With --fixture-file, read the same file in a loop",
		)

		flags.StringVarP(
			&output.TransportFile,
			"out-file",
			"",
			"",
			"format to display output on a file; ie --out-file nmea@/tmp/foo.txt",
		)
	}

	defaults.SetDefaults(output)

	return output
}

// UserSettings reads the user configuration.
func UserSettings(flags *pflag.FlagSet) (*Config, error) {
	config := newConfig(flags)

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

// Path is the path of the config file.
func (s Config) Path() string {
	return filepath.Join(s.basePath, settingsFilename)
}

// AircraftDatabaseFile is the path of the aircraft database.
func (s Config) AircraftDatabaseFile() string {
	return filepath.Join(s.basePath, s.AircraftDatabaseFilename)
}
