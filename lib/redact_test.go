package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactJob(t *testing.T) {
	job := &Job{
		ID:   "job-123",
		Name: "test-job",
	}
	job.Remoteci.APISecret = "secret-api-key"
	job.Topic.Data.PullSecret.Auths.CloudOpenshiftCom.Auth = "cloud-auth"
	job.Topic.Data.PullSecret.Auths.QuayIo.Auth = "quay-auth"
	job.Topic.Data.PullSecret.Auths.RegistryCiOpenshiftOrg.Auth = "registry-ci-auth"
	job.Topic.Data.PullSecret.Auths.RegistryConnectRedhatCom.Auth = "registry-connect-auth"
	job.Topic.Data.PullSecret.Auths.RegistryRedhatIo.Auth = "registry-redhat-auth"

	RedactJob(job)

	assert.Equal(t, redactedPlaceholder, job.Remoteci.APISecret)
	assert.Equal(t, redactedPlaceholder, job.Topic.Data.PullSecret.Auths.CloudOpenshiftCom.Auth)
	assert.Equal(t, redactedPlaceholder, job.Topic.Data.PullSecret.Auths.QuayIo.Auth)
	assert.Equal(t, redactedPlaceholder, job.Topic.Data.PullSecret.Auths.RegistryCiOpenshiftOrg.Auth)
	assert.Equal(t, redactedPlaceholder, job.Topic.Data.PullSecret.Auths.RegistryConnectRedhatCom.Auth)
	assert.Equal(t, redactedPlaceholder, job.Topic.Data.PullSecret.Auths.RegistryRedhatIo.Auth)

	assert.Equal(t, "job-123", job.ID)
	assert.Equal(t, "test-job", job.Name)
}

func TestRedactRemoteCI(t *testing.T) {
	rci := &RemoteCI{
		ID:        "rci-123",
		Name:      "test-rci",
		APISecret: "secret-key",
	}

	RedactRemoteCI(rci)

	assert.Equal(t, redactedPlaceholder, rci.APISecret)
	assert.Equal(t, "rci-123", rci.ID)
	assert.Equal(t, "test-rci", rci.Name)
}
