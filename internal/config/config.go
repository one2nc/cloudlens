package config

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/adrg/xdg"
	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	awsV2Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
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

var config Config

func GetSession(profile, region string) (*session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{Config: aws.Config{
		//TODO: remove hardcoded enpoint
		Endpoint:         aws.String("http://localhost:4566"),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	},
		Profile: profile})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return nil, err
	}
	return sess, nil
}

func Get() (Config, error) {
	emptyCfg := Config{}
	if reflect.DeepEqual(emptyCfg, config) {
		profiles, err := GetProfiles()
		if err != nil {
			return emptyCfg, err
		}
		config.Profiles = profiles
		if LookupForValue(config.Profiles, "default") {
			// Load the Shared AWS Configuration (~/.aws/config)
			awsLocalCfg, err := awsV2Config.LoadDefaultConfig(context.TODO())
			if err != nil {
				return emptyCfg, err
			}
			config.AwsConfig = awsLocalCfg
		}
	}
	return config, nil
}

func GetProfiles() (profiles []string, err error) {
	filepath := defaults.SharedCredentialsFilename()
	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return profiles, err
	}
	lines := strings.Split(string(fileContent), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			profile := line[1 : len(line)-1]
			profiles = append(profiles, profile)
		}
	}
	if len(profiles) < 1 {
		err = errors.New("NO PROFILES FOUND")
		return nil, err
	}

	return profiles, nil
}
