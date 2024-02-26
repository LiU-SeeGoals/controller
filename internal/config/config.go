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

	// Simulation config
	Simulation ConfigSimulation

	// Game controller config
	GC ConfigGameController

	// Base station config
	BaseStation ConfigBaseStation

	// Misc config
	Misc ConfigMisc
}

type ConfigMisc struct {
	// Probably removed later (idk how we get which team we are on in the future)
	IsBlueTeam bool `env:"IS_BLUE_TEAM,required"`
}

type ConfigBaseStation struct {
	// Base station address
	Address string `env:"BASE_STATION_ADDR,required"`

	// Base station vision port
	VisionPort string `env:"BASE_STATION_VISION_PORT,required"`

	// Base station action port
	ActionPort string `env:"BASE_STATION_ACTION_PORT,required"`
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

	// Amount of cameras
	AmountOfCameras uint32 `env:"SSL_VISION_AMOUNT_OF_CAMERAS,required"`
}


type ConfigSimulation struct {
	// Simulation used
	SimulationUsed string `env:"SIMULATION_USED,required"`

	// Grsim simulation config
	Grsim GrsimSimulation
}

// Config struct for grsim
//
// Check the .env file to make changes.
// Unless you have a good reason, you
// shouldn't need to use this directly.
type GrsimSimulation struct {
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

	// Amount of Cameras
	// Used to normalize the data from the simulation
	AmountOfCameras uint32 `env:"GRSIM_AMOUNT_OF_CAMERAS,required"`
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
func GetSimCommandAddress() string {
	cfg := GetInstance()
	if cfg.Simulation.SimulationUsed == "grsim" {
		address := cfg.Simulation.Grsim.Address
		port := cfg.Simulation.Grsim.CommandPort
		return fmt.Sprintf("%s:%s", address, port)
	}
	panic("Unknown simulation used.")
}

func GetAmountOfCameras() uint32 {
	cfg := GetInstance()
	if cfg.Env == "sim" {
		return GetAmountOfCamerasSim()
	}
	return cfg.SSLVision.AmountOfCameras
}

func GetAmountOfCamerasSim() uint32 {
	cfg := GetInstance()
	if cfg.Simulation.SimulationUsed == "grsim" {
		return cfg.Simulation.Grsim.AmountOfCameras
	}
	panic("Unknown simulation used.")
}

// GetBaseStationVisionPort returns the base station vision port from the config.
func GetBaseStationVisionAddress() string {
	cfg := GetInstance()
	return fmt.Sprintf("%s:%s", cfg.BaseStation.Address, cfg.BaseStation.VisionPort)
}

// GetBaseStationVisionPort returns the base station vision port from the config.
func GetBaseStationActionAddress() string {
	cfg := GetInstance()
	return fmt.Sprintf("%s:%s", cfg.BaseStation.Address, cfg.BaseStation.ActionPort)
}

func GetIsBlueTeam() bool {
	return GetInstance().Misc.IsBlueTeam
}
