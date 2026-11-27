package aws

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProfilesCombinesCredentialsAndConfig(t *testing.T) {
	dir := t.TempDir()
	awsDir := filepath.Join(dir, ".aws")
	assert.NoError(t, os.Mkdir(awsDir, 0700))
	credentialsFile := filepath.Join(awsDir, "credentials")
	configFile := filepath.Join(awsDir, "config")

	assert.NoError(t, os.WriteFile(credentialsFile, []byte("[default]\naws_access_key_id=test\n[static]\naws_access_key_id=test\n"), 0600))
	assert.NoError(t, os.WriteFile(configFile, []byte("[default]\nregion=us-east-1\n[profile assumed]\nrole_arn=arn:aws:iam::123456789012:role/test\nsource_profile=static\n[sso-session company]\nsso_region=us-east-1\n"), 0600))
	t.Setenv("HOME", dir)

	profiles, err := GetProfiles()

	assert.NoError(t, err)
	assert.Equal(t, []string{"default", "static", "assumed"}, profiles)
}

func TestGetLocalstackCfgFallsBackToStaticCredsWithNoRealCredsConfigured(t *testing.T) {
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_PROFILE")

	cfg, err := GetLocalstackCfg("us-east-1")
	assert.Nil(t, err)

	got, err := cfg.Credentials.Retrieve(context.TODO())
	assert.Nil(t, err)

	assert.Equal(t, "test", got.AccessKeyID)
	assert.Equal(t, "test", got.SecretAccessKey)
}
