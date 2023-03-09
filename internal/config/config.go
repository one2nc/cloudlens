package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// CloudlensConfig represents Cloudlens configuration dir env var.
const CloudlensConfig = "CLOUDLENSCONFIG"

var (
	//CloudlensConfigFile represents config file location.
	CloudlensConfigFile = filepath.Join(CloudlensHome(), "config.yml")
)

type Config struct {
	Cloudlens *Cloudlens `yaml:"cloudlens"`
	// List of profiles in (~/.aws/credentials)
	Profiles  []string
	AwsConfig awsV2.Config
}

// CloudlensHome returns Cloudlens configs home directory.
func CloudlensHome() string {
	if env := os.Getenv(CloudlensConfig); env != "" {
		//log.Debug().Msg("env CL: " + env)
		return env
	}
	xdgCLHome, err := xdg.ConfigFile("cloudlens")
	//log.Debug().Msg("xdgsclhome: " + xdgCLHome)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create configuration directory for cloudlens")
	}

	return xdgCLHome
}

// Load cloudlens configuration from file.
func (c *Config) Load(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	c.Cloudlens = NewCloudlens()

	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return err
	}
	if cfg.Cloudlens != nil {
		c.Cloudlens = cfg.Cloudlens
	}
	return nil
}

// Save configuration to disk.
func (c *Config) Save() error {
	//c.Validate()

	return c.SaveFile(CloudlensConfigFile)
}

// SaveFile K9s configuration to disk.
func (c *Config) SaveFile(path string) error {
	EnsurePath(path, DefaultDirMod)
	cfg, err := yaml.Marshal(c)
	if err != nil {
		log.Error().Msgf("[Config] Unable to save cloudlens config file: %v", err)
		return err
	}
	log.Info().Msg(fmt.Sprintf("Config Path: %v", path))
	return os.WriteFile(path, cfg, 0644)
}

// Unsed for now

// var config Config

// func Get() (Config, error) {
// 	emptyCfg := Config{}
// 	if reflect.DeepEqual(emptyCfg, config) {
// 		profiles, err := GetProfiles()
// 		if err != nil {
// 			return emptyCfg, err
// 		}
// 		config.Profiles = profiles
// 		if LookupForValue(config.Profiles, "default") {
// 			// Load the Shared AWS Configuration (~/.aws/config)
// 			awsLocalCfg, err := awsV2Config.LoadDefaultConfig(context.TODO())
// 			if err != nil {
// 				return emptyCfg, err
// 			}
// 			config.AwsConfig = awsLocalCfg
// 		}
// 	}
// 	return config, nil
// }
