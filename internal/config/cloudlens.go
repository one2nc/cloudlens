package config

type Active struct {
	Profile string `yaml:"profile"`
	Region  string `yaml:"region"`
	View    string `yaml:"view"`
}

type Cloudlens struct {
	EnableMouse bool    `yaml:"enableMouse"`
	Headless    bool    `yaml:"headless"`
	Logoless    bool    `yaml:"logoless"`
	Crumbsless  bool    `yaml:"crumbsless"`
	Active      *Active `yaml:"active"`
}

// NewCloudlens create a new Cloudlens configuration.
func NewCloudlens() *Cloudlens {
	return &Cloudlens{}
}
