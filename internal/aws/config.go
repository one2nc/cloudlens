package aws

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	awsV2Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/ini.v1"
)

type Profiles struct {
	Data  []string
	Error string
}
type credentialProvider struct {
	awsV2.Credentials
}

func (c credentialProvider) Retrieve() (credentials.Value, error) {
	return credentials.Value{AccessKeyID: c.AccessKeyID, SecretAccessKey: c.SecretAccessKey, SessionToken: os.Getenv("AWS_SESSION_TOKEN")}, nil
}

func (c credentialProvider) IsExpired() bool {
	return c.Expired()
}

// func GetSession(profile, region string) (*session.Session, error) {
// 	log.Info().Msg("Profiles inside GetSession:" + profile + " and region inside GetSession:" + region)
// 	cfg, err := awsV2Config.LoadDefaultConfig(context.Background(),
// 		awsV2Config.WithSharedConfigProfile(profile),
// 		awsV2Config.WithRegion(region),
// 	)
// 	if err != nil {
// 		fmt.Printf("failed to load config")
// 		return nil, err
// 	}
// 	creds, err := cfg.Credentials.Retrieve(context.Background())
// 	if err != nil {
// 		log.Info().Msg("Failed to read credentials.")
// 		return nil, err
// 	}
// 	log.Info().Msg("Access key id is:" + creds.AccessKeyID)
// 	log.Info().Msg("Secret Key id is:" + creds.SecretAccessKey)
// 	credentialProvider := credentialProvider{Credentials: creds}
// 	if credentialProvider.IsExpired() {
// 		fmt.Println("Credentials have expired")
// 		return nil, errors.New("AWS Credentials expired")
// 	}

// 	// create session
// 	sess, err := session.NewSessionWithOptions(session.Options{Config: aws.Config{
// 		//TODO: remove hardcoded enpoint
// 		//Endpoint:         aws.String(localstackEndpoint),
// 		Credentials:      credentials.NewCredentials(credentialProvider),
// 		Region:           aws.String(region),
// 		S3ForcePathStyle: aws.Bool(true),
// 	},
// 		Profile: profile})
// 	if err != nil {
// 		fmt.Println("Failed to create session", err)
// 		os.Exit(1)
// 	}
// 	return sess, nil
// }

func GetSession(profile, region string) (*session.Session, error) {
	cfg, err := awsV2Config.LoadDefaultConfig(
		context.TODO(),
		awsV2Config.WithSharedConfigProfile(profile),
		awsV2Config.WithRegion(region),
	)
	if err != nil {
		fmt.Printf("failed to load config")
		return nil, err
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Printf("failed to read credentials")
		return nil, err
	}

	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		fmt.Println("Credentials have expired")
		return nil, errors.New("AWS Credentials expired")
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

	// token, err := cfg.BearerAuthTokenProvider.RetrieveBearerToken(context.Background())
	// if err != nil {
	// 	return nil, err
	// }
	// log.Info().Msg("Token is: " + token.Value)

	// credsS := credentials.Value{
	// 	AccessKeyID:     creds.AccessKeyID,
	// 	SecretAccessKey: creds.SecretAccessKey,
	// 	SessionToken:    token.Value,
	// }

	// credential := credentials.NewStaticCredentialsFromCreds(credsS)

	// config := aws.NewConfig().WithCredentials(credential).WithRegion("ap-south-1")
	// sess1, err := session.NewSession(config)
	// if err != nil {
	// 	panic(err)
	// }

	// sssoStartURL := "https://my-sso-portal.awsapps.com/start"
	// ssoRegion := "us-east-1"
	// ssoAccountID := "123456789012"

	// sess, err = session.NewSessionWithOptions(session.Options{
	// 	Config: aws.Config{
	// 		Region: aws.String(ssoRegion),
	// 		Credentials: credentials.NewCredentials(&credentials.{
	// 			StartURL:  ssoStartURL,
	// 			AccountID: ssoAccountID,
	// 		}),
	// 	},
	// })

	// if err != nil {
	// 	fmt.Println("Error creating session: ", err)
	// 	return
	// }

	return sess, nil
}

func GetCfg(profile, region string) (awsV2.Config, error) {
	cfg, err := awsV2Config.LoadDefaultConfig(
		context.TODO(),
		awsV2Config.WithSharedConfigProfile(profile),
		awsV2Config.WithRegion(region),
	)
	if err != nil {
		fmt.Printf("failed to load config")
		return awsV2.Config{}, err
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Printf("failed to read credentials")
		return awsV2.Config{}, err
	}

	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		fmt.Println("Credentials have expired")
		return awsV2.Config{}, errors.New("AWS Credentials expired")
	}
	return cfg, err
}

func GetSessionUsingEnvVariables(region, profile string) (*session.Session, error) {
	akid := aws.String(os.Getenv(AWS_ACCESS_KEY_ID))
	secKey := aws.String(os.Getenv(AWS_SECRET_ACCESS_KEY))
	//~/.aws/config and ~/.aws/credentials file are not present and even the env variables are not set.
	if *akid == "" || *secKey == "" {
		return nil, errors.New("Cannot find AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY")
	}
	creds := awsV2.Credentials{AccessKeyID: *akid, SecretAccessKey: *secKey}
	credentialProvider := credentialProvider{Credentials: creds}
	if credentialProvider.IsExpired() {
		fmt.Println("Credentials have expired")
		return nil, errors.New("AWS Credentials expired")
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(*akid, *secKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func GetProfiles() (profiles []string, err error) {
	fpCred := defaults.SharedCredentialsFilename()
	_, errCred := os.Stat(fpCred)
	fpConf := defaults.SharedConfigFilename()
	_, errConf := os.Stat(fpConf)
	if os.IsNotExist(errCred) && os.IsNotExist(errConf) {
		return nil, errConf
	}
	var ret []string
	defaultReturn := &Profiles{Data: nil, Error: ""}
	fp := defaults.SharedCredentialsFilename()
	_, err = os.Stat(fp)
	if os.IsNotExist(err) {
		fp = defaults.SharedConfigFilename()
	}
	f, err := ini.Load(fp) // Load ini file
	if err != nil {
		defaultReturn.Error = err.Error()
	} else {
		arr := []string{}
		for _, v := range f.Sections() {
			if len(v.Keys()) != 0 {
				arr = append(arr, v.Name())
			}
		}
		defaultReturn.Data = arr
	}
	for i := 0; i < len(defaultReturn.Data); i++ {
		spltiArr := strings.Split(defaultReturn.Data[i], " ")
		if len(spltiArr) == 1 {
			ret = append(ret, spltiArr[len(spltiArr)-1])
		} else if len(spltiArr) > 1 && spltiArr[0] == "profile" {
			ret = append(ret, spltiArr[len(spltiArr)-1])
		}
	}
	return ret, nil
}
