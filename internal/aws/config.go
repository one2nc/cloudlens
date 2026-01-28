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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/one2nc/cloudlens/internal"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
)

type Profiles struct {
	Data  []string
	Error string
}
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
	files, err := profileFiles()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	seen := map[string]struct{}{}
	for _, fp := range files {
		if err = appendProfilesFromFile(fp, seen, &ret); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func profileFiles() ([]string, error) {
	fpCred := defaults.SharedCredentialsFilename()
	_, errCred := os.Stat(fpCred)
	fpConf := defaults.SharedConfigFilename()
	_, errConf := os.Stat(fpConf)
	if os.IsNotExist(errCred) && os.IsNotExist(errConf) {
		return nil, errConf
	}
	files := make([]string, 0, 2)
	if !os.IsNotExist(errCred) {
		files = append(files, fpCred)
	}
	if !os.IsNotExist(errConf) {
		files = append(files, fpConf)
	}
	return files, nil
}

func appendProfilesFromFile(fp string, seen map[string]struct{}, ret *[]string) error {
	f, err := ini.Load(fp)
	if err != nil {
		return err
	}
	for _, v := range f.Sections() {
		if len(v.Keys()) == 0 {
			continue
		}
		if name, ok := profileNameFromSection(v.Name()); ok {
			addProfile(name, seen, ret)
		}
	}
	return nil
}

func profileNameFromSection(section string) (string, bool) {
	parts := strings.Split(section, " ")
	if len(parts) == 1 {
		return parts[0], true
	}
	if len(parts) > 1 && parts[0] == "profile" {
		return parts[len(parts)-1], true
	}
	return "", false
}

func addProfile(name string, seen map[string]struct{}, ret *[]string) {
	if _, ok := seen[name]; ok {
		return
	}
	seen[name] = struct{}{}
	*ret = append(*ret, name)
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
	)
	if err != nil {
		log.Fatal().Err(err)
	}
	return awsLSCfg, nil
}

func GetLocastackEndpoint() string {

	port := os.Getenv(internal.LOCALSTACK_PORT)

	return fmt.Sprintf("http://localhost:%v", port)
}
