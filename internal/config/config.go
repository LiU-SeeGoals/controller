package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config wraps around specific config object.
//
// Check the .env file to make changes.
// Unless you have a good reason, you shouldn't need to use this directly.
// It's better to use the helper methods in this file instead.
type Config struct {
	// Environment - e.g. "docker"
	Env string `env:"ENVIRONMENT,required"`

	// SSL vision config
	SSLVision ConfigSSLVision

	// Grsim config
	Grsim ConfigGrsim

	// Game controller config
	GC ConfigGameController
}

// Config struct for SSL Vision.
//
// Check the .env file to make changes.
// Unless you have a good reason, you
// shouldn't need to use this directly.
type ConfigSSLVision struct {
	// Multicast address.
	Address string `env:"SSL_VISION_MULTICAST_ADDR,required"`

	// Tracker, detection, and geometry packets.
	Port string `env:"SSL_VISION_MAIN_PORT,required"`

	// Visualization packets.
	VizPort string `env:"SSL_VISION_VIZ_PORT,required"`
}

// Config struct for grsim
//
// Check the .env file to make changes.
// Unless you have a good reason, you
// shouldn't need to use this directly.
type ConfigGrsim struct {
	// Grsim address
	Address string `env:"GRSIM_ADDR,required"`

	// Command listen port.
	// Accepts robots commands.
	CommandPort string `env:"GRSIM_COMMAND_LISTEN_PORT,required"`

	// Blue team status listen port.
	// Use unknown.
	BlueStatusPort string `env:"GRSIM_BLUE_STATUS_SEND_PORT,required"`

	// Yellow team status listen port.
	// Use unknown.
	YellowStatusPort string `env:"GRSIM_YELLOW_STATUS_SEND_PORT,required"`

	// Simulation controller send port.
	// Use unknown.
	SimControllerPort string `env:"GRSIM_SIM_CONTROLLER_SEND_PORT,required"`

	// Blue team controller listen port.
	// Use unknown.
	BlueControllerPort string `env:"GRSIM_BLUE_CONTROLLER_LISTEN_PORT,required"`

	// Yellow team controller listen port.
	// Use unknown.
	YellowControllerPort string `env:"GRSIM_YELLOW_CONTROLLER_LISTEN_PORT,required"`
}

// Config struct for game controller (GC)
//
// Check the .env file to make changes.
// Unless you have a good reason, you
// shouldn't need to use this directly.
type ConfigGameController struct {
	// GC multicast publish address
	Address string `env:"GC_PUBLISH_ADDR,required"`

	// GC publish port
	Port string `env:"GC_PUBLISH_PORT,required"`
}

var (
	// Config instance
	instance *Config

	// Init helper
	once sync.Once
)

// Loads config from .env file.
// Config object saved as global object in this file.
func loadConfig() {
	// Determine the path to the config file.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("loadConfig: unable to determine current directory")
	}
	envPath := filepath.Join(filepath.Dir(filename), "../../.env")
	envPath = filepath.Clean(envPath)

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Unable to load .env file: %e", err)
	}

	cfg := &Config{}
	err = env.Parse(cfg)
	if err != nil {
		log.Fatalf("Unable to parse config: %e", err)
	}

	instance = cfg
}

// GetInstance returns the singleton instance of the Config object,
// initializing it if necessary.
func GetInstance() *Config {
	once.Do(loadConfig)
	return instance
}

// GetSSLClientAddress returns the SSL client address from the config.
func GetSSLClientAddress() string {
	cfg := GetInstance()
	return fmt.Sprintf("%s:%s", cfg.SSLVision.Address, cfg.SSLVision.Port)
}

// GetGrSimAddress returns the GrSim address from the config.
func GetGrSimAddress() string {
	cfg := GetInstance()
	return fmt.Sprintf("%s:%s", cfg.Grsim.Address, cfg.Grsim.CommandPort)
}
