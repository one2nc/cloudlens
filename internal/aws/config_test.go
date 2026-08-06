package aws

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
