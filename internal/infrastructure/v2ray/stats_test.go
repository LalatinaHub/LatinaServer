package v2ray

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTransportCredentials(t *testing.T) {
	// Clean env
	os.Unsetenv("V2RAY_API_TLS")
	os.Unsetenv("V2RAY_API_CERT")
	os.Unsetenv("V2RAY_API_TLS_INSECURE")

	// Default should be insecure
	creds := getTransportCredentials()
	assert.NotNil(t, creds)
	assert.Equal(t, "insecure", creds.Info().SecurityProtocol)

	// TLS enabled
	os.Setenv("V2RAY_API_TLS", "true")
	credsTLS := getTransportCredentials()
	assert.NotNil(t, credsTLS)
	assert.Equal(t, "tls", credsTLS.Info().SecurityProtocol)
	os.Unsetenv("V2RAY_API_TLS")
}

func TestGetUsersStatsBatch_Empty(t *testing.T) {
	usages := GetUsersStatsBatch([]string{})
	assert.Empty(t, usages)
}

func TestClose_NoActiveConn(t *testing.T) {
	err := Close()
	assert.NoError(t, err)
}
