package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	awsV2Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Config struct {
	// List of profiles in (~/.aws/credentials)
	Profiles  []string
	AwsConfig awsV2.Config
}

var config Config

func GetSession(profile, region string, awsCfg awsV2.Config) (*session.Session, error) {
	crds, _ := awsCfg.Credentials.Retrieve(context.TODO())
	sess, err := session.NewSessionWithOptions(session.Options{Config: aws.Config{
		//TODO: remove hardcoded enpoint
		Endpoint:    aws.String("http://localhost:4566"),
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(crds.AccessKeyID, crds.SecretAccessKey, ""),
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
		// Load the Shared AWS Configuration (~/.aws/config)
		awsLocalCfg, err := awsV2Config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return emptyCfg, err
		}
		config.AwsConfig = awsLocalCfg
		creds, err := awsLocalCfg.Credentials.Retrieve(context.TODO())
		if err != nil {
			return emptyCfg, err
		}
		config.Profiles, err = GetProfiles(strings.Split(creds.Source, " ")[1])
		if err != nil {
			return emptyCfg, err
		}
	}
	return config, nil
}

func GetProfiles(filepath string) ([]string, error) {
	profiles := []string{}
	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(fileContent), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			profile := line[1 : len(line)-1]
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}
