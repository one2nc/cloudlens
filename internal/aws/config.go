package aws

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	awsV2Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
)

type credentialProvider struct {
	awsV2.Credentials
}

func (c credentialProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey, SessionToken: os.Getenv("AWS_SESSION_TOKEN")}, nil
}

func (c credentialProvider) IsExpired() bool {
	return c.Expired()
}

func GetSession(profile, region string) (*session.Session, error) {
	cfg, err := awsV2Config.LoadDefaultConfig(context.TODO(),
		awsV2Config.WithSharedConfigProfile(profile),
		awsV2Config.WithRegion(region),
	)
	if err != nil {
		fmt.Printf("failed to load config")
		return nil, err
	}
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		fmt.Println("Credentials have expired")
		return nil, errors.New("AWS Credentials expired")
	}
	if err != nil {
		fmt.Printf("failed to read credentials")
		return nil, err
	}
	if err != nil {
		panic(fmt.Sprintf("failed to load SDK configuration, %v", err))
	}

	// create session
	sess, err := session.NewSessionWithOptions(session.Options{Config: aws.Config{
		//TODO: remove hardcoded enpoint
		//Endpoint:         aws.String(localstackEndpoint),
		Credentials:      credentials.NewCredentials(credentialProvider),
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
