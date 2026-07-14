package lib

const redactedPlaceholder = "**REDACTED**"

func RedactJob(job *Job) {
	job.Remoteci.APISecret = redactedPlaceholder
	job.Topic.Data.PullSecret.Auths.CloudOpenshiftCom.Auth = redactedPlaceholder
	job.Topic.Data.PullSecret.Auths.QuayIo.Auth = redactedPlaceholder
	job.Topic.Data.PullSecret.Auths.RegistryCiOpenshiftOrg.Auth = redactedPlaceholder
	job.Topic.Data.PullSecret.Auths.RegistryConnectRedhatCom.Auth = redactedPlaceholder
	job.Topic.Data.PullSecret.Auths.RegistryRedhatIo.Auth = redactedPlaceholder
}

func RedactRemoteCI(rci *RemoteCI) {
	rci.APISecret = redactedPlaceholder
}
