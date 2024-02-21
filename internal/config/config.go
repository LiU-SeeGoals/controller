package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
)

type Config struct {
	SSLClientAddress *string `json:"sslClientAddress"`
	GrSimAddress     *string `json:"grSimAddress"`
	BaseStationAddress *string `json:"baseStationAddress"`
	BaseStationVisionPort *int `json:"baseStationVisionPort"`
	BaseStationActionPort *int `json:"baseStationActionPort"`
	AmountOfCameras  *uint32    `json:"amountOfCameras"`
	IsBlueTeam *bool `json:"isBlueTeam"`
}

var (
	instance *Config
	once     sync.Once
)

func loadConfig() {
	// Determine the path to the config file.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("loadConfig: unable to determine current directory")
	}
	configPath := filepath.Join(filepath.Dir(filename), "../../config.json")
	configPath = filepath.Clean(configPath)

	// Open the config file.
	file, err := os.Open(configPath)
	if err != nil {
		panic(fmt.Errorf("loadConfig: %w", err))
	}
	defer file.Close()

	// Deserialize the json.
	instance = &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(instance); err != nil {
		panic(fmt.Errorf("loadConfig: %w", err))
	}

	// Validate required fields using reflection
	val := reflect.ValueOf(instance).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			panic(fmt.Sprintf("loadConfig: missing required field '%s'", val.Type().Field(i).Name))
		}
	}
}

// GetInstance returns the singleton instance of the Config object, initializing it if necessary.
func GetInstance() *Config {
	once.Do(loadConfig)
	return instance
}

// GetSSLClientAddress returns the SSL client address from the config.
func GetSSLClientAddress() string {
	return *GetInstance().SSLClientAddress
}

// GetGrSimAddress returns the GrSim address from the config.
func GetGrSimAddress() string {
	return *GetInstance().GrSimAddress
}

// GetAmountOfCameras returns the amount of cameras from the config.
func GetAmountOfCameras() uint32 {
	return *GetInstance().AmountOfCameras
}

// GetBaseStationClientAddress returns the base station client address from the config.
func GetBaseStationAddress() string {
	return *GetInstance().BaseStationAddress
}

// GetBaseStationVisionPort returns the base station vision port from the config.
func GetBaseStationVisionPort() int {
	return *GetInstance().BaseStationVisionPort
}

// GetBaseStationActionPort returns the base station action port from the config.
func GetBaseStationActionPort() int {
	return *GetInstance().BaseStationActionPort
}

func GetIsBlueTeam() bool {
	return *GetInstance().IsBlueTeam
}
