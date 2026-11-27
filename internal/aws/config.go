package aws

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsV2Config "github.com/aws/aws-sdk-go-v2/config"
	creds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/one2nc/cloudlens/internal"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
)

type credentialProvider struct {
	awsV2.Credentials
}

type AWSConfigInput struct {
	Profile, Region string
	UseLocalStack   bool
	UseEnvVariables bool
}

func (c credentialProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey, SessionToken: os.Getenv("AWS_SESSION_TOKEN")}, nil
}

func (c credentialProvider) IsExpired() bool {
	return c.Expired()
}

func GetCfg(cfgInput AWSConfigInput) (awsV2.Config, error) {

	var cfg awsV2.Config
	var err error

	if cfgInput.UseLocalStack {
		cfg, err = GetLocalstackCfg(cfgInput.Region)
	} else if cfgInput.UseEnvVariables {
		cfg, err = GetCfgUsingEnvVariables(cfgInput.Profile, cfgInput.Region)
	} else {
		cfg, err = awsV2Config.LoadDefaultConfig(
			context.TODO(),
			awsV2Config.WithSharedConfigProfile(cfgInput.Profile),
			awsV2Config.WithRegion(cfgInput.Region),
		)
	}

	if err != nil {
		log.Print("failed to load config")
		return awsV2.Config{}, err
	}
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		log.Print("failed to read credentials ", err)
		return awsV2.Config{}, err
	}

	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		log.Print("Credentials have expired")
		return awsV2.Config{}, errors.New("AWS Credentials expired")
	}
	return cfg, err
}

func GetCfgUsingEnvVariables(profile, region string) (awsV2.Config, error) {
	akid := aws.String(os.Getenv(AWS_ACCESS_KEY_ID))
	secKey := aws.String(os.Getenv(AWS_SECRET_ACCESS_KEY))
	cfg, err := awsV2Config.LoadDefaultConfig(
		context.TODO(),
		awsV2Config.WithSharedConfigProfile(profile),
		awsV2Config.WithRegion(region),
		config.WithCredentialsProvider(
			creds.NewStaticCredentialsProvider(*akid, *secKey, ""),
		),
	)
	if err != nil {
		log.Print("failed to load config")
		return awsV2.Config{}, err
	}
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		log.Print("failed to read credentials ", err)
		return awsV2.Config{}, err
	}

	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		log.Print("Credentials have expired")
		return awsV2.Config{}, errors.New("AWS Credentials expired")
	}
	return cfg, err
}

func GetProfiles() (profiles []string, err error) {
	files := []struct {
		path         string
		sharedConfig bool
	}{
		{path: defaults.SharedCredentialsFilename()},
		{path: defaults.SharedConfigFilename(), sharedConfig: true},
	}

	seen := make(map[string]struct{})
	foundFile := false
	for _, file := range files {
		if _, statErr := os.Stat(file.path); statErr != nil {
			if os.IsNotExist(statErr) {
				continue
			}
			return nil, statErr
		}
		foundFile = true

		f, loadErr := ini.Load(file.path)
		if loadErr != nil {
			return nil, loadErr
		}
		for _, section := range f.Sections() {
			if len(section.Keys()) == 0 {
				continue
			}

			name := section.Name()
			if file.sharedConfig && name != "default" {
				const profilePrefix = "profile "
				if !strings.HasPrefix(name, profilePrefix) || len(name) == len(profilePrefix) {
					continue
				}
				name = name[len(profilePrefix):]
			}
			if _, exists := seen[name]; !exists {
				profiles = append(profiles, name)
				seen[name] = struct{}{}
			}
		}
	}

	if !foundFile {
		return nil, os.ErrNotExist
	}
	return profiles, nil
}

func GetLocalstackCfg(region string) (awsV2.Config, error) {
	customResolver := awsV2.EndpointResolverFunc(func(service, region string) (awsV2.Endpoint, error) {
		return awsV2.Endpoint{
			URL:           GetLocastackEndpoint(),
			SigningRegion: region,
		}, nil
	})

	awsLSCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithEndpointResolver(customResolver),
		config.WithCredentialsProvider(localstackCredentialsProvider(region)),
	)
	if err != nil {
		log.Fatal().Err(err)
	}
	return awsLSCfg, nil
}

func localstackCredentialsProvider(region string) awsV2.CredentialsProviderFunc {
	return func(ctx context.Context) (awsV2.Credentials, error) {
		real, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithEC2IMDSClientEnableState(imds.ClientDisabled),
		)
		if err == nil {
			if realCreds, err := real.Credentials.Retrieve(ctx); err == nil {
				return realCreds, nil
			}
		}
		return creds.NewStaticCredentialsProvider("test", "test", "").Retrieve(ctx)
	}
}

func GetLocastackEndpoint() string {

	port := os.Getenv(internal.LOCALSTACK_PORT)

	return fmt.Sprintf("http://localhost:%v", port)
}
