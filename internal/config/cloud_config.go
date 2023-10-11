package config

type CloudConfig struct {
	SelectedCloud string
	AWSConfig
	GCPConfig
}

type AWSConfig struct {
	Profile string
	Region  string
}
type GCPConfig struct {
	CredFilePath string
}

func NewCloudConfig() CloudConfig {

	return CloudConfig{}
}
