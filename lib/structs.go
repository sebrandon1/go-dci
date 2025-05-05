package lib

// Shared structs

type Meta struct {
	Count int `json:"count,omitempty"`
}

type Product struct {
	CreatedAt   string `json:"created_at,omitempty"`
	Description string `json:"description,omitempty"`
	Etag        string `json:"etag,omitempty"`
	ID          string `json:"id,omitempty"`
	Label       string `json:"label,omitempty"`
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type Data struct {
	Digest      []string `json:"digest"`
	DisplayName string   `json:"display_name"`
	PullURL     string   `json:"pull_url"`
	Tags        []string `json:"tags"`
	UID         string   `json:"uid"`
	URL         string   `json:"url"`
	Version     string   `json:"version"`
}

// Refactored structs

type JsonCertsuiteInfo struct {
	ID               string `json:"id"`
	CertsuiteVersion string `json:"certsuite_version"`
	OCPVersion       string `json:"ocp_version"`
}

type JsonOcpVersionCount struct {
	OcpVersion string `json:"ocp_version"`
	RunCount   int    `json:"run_count"`
}

type OcpJsonOutput struct {
	OcpVersions []JsonOcpVersionCount `json:"ocp_versions"`
}

type JobsJsonOutput struct {
	Jobs []JsonCertsuiteInfo `json:"jobs"`
}

type Components struct {
	CanonicalProjectName string   `json:"canonical_project_name"`
	CreatedAt            string   `json:"created_at,omitempty"`
	Data                 Data     `json:"data,omitempty"`
	DisplayName          string   `json:"display_name"`
	Etag                 string   `json:"etag"`
	ID                   string   `json:"id"`
	Message              string   `json:"message"`
	Name                 string   `json:"name"`
	ReleasedAt           string   `json:"released_at,omitempty"`
	State                string   `json:"state"`
	Tags                 []string `json:"tags"`
	TeamID               any      `json:"team_id"`
	Title                string   `json:"title"`
	TopicID              string   `json:"topic_id"`
	Type                 string   `json:"type"`
	UID                  string   `json:"uid"`
	UpdatedAt            string   `json:"updated_at,omitempty"`
	URL                  string   `json:"url"`
	Version              string   `json:"version"`
}

type Job struct {
	ClientVersion string       `json:"client_version"`
	Comment       string       `json:"comment"`
	Components    []Components `json:"components"`
	Configuration string       `json:"configuration"`
	CreatedAt     string       `json:"created_at,omitempty"`
	Duration      int          `json:"duration"`
	Etag          string       `json:"etag"`
	ID            string       `json:"id"`
	KeysValues    []struct {
		JobID string  `json:"job_id"`
		Key   string  `json:"key"`
		Value float64 `json:"value"`
	} `json:"keys_values"`
	Name     string `json:"name"`
	Pipeline struct {
		CreatedAt string `json:"created_at,omitempty"`
		Etag      string `json:"etag"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		State     string `json:"state"`
		TeamID    string `json:"team_id"`
		UpdatedAt string `json:"updated_at,omitempty"`
	} `json:"pipeline"`
	PipelineID    string `json:"pipeline_id"`
	PreviousJobID any    `json:"previous_job_id"`
	ProductID     string `json:"product_id"`
	Remoteci      struct {
		APISecret string `json:"api_secret"`
		CreatedAt string `json:"created_at,omitempty"`
		Data      struct {
		} `json:"data"`
		Etag      string `json:"etag"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		Public    bool   `json:"public"`
		State     string `json:"state"`
		TeamID    string `json:"team_id"`
		UpdatedAt string `json:"updated_at,omitempty"`
	} `json:"remoteci"`
	RemoteciID string `json:"remoteci_id"`
	Results    []struct {
		CreatedAt    string `json:"created_at"`
		Errors       int    `json:"errors"`
		Failures     int    `json:"failures"`
		FileID       string `json:"file_id"`
		ID           string `json:"id"`
		JobID        string `json:"job_id"`
		Name         string `json:"name"`
		Regressions  int    `json:"regressions"`
		Skips        int    `json:"skips"`
		Success      int    `json:"success"`
		Successfixes int    `json:"successfixes"`
		Time         int    `json:"time"`
		Total        int    `json:"total"`
		UpdatedAt    string `json:"updated_at"`
	} `json:"results"`
	State        string   `json:"state"`
	Status       string   `json:"status"`
	StatusReason string   `json:"status_reason"`
	Tags         []string `json:"tags"`
	Team         struct {
		Country             any    `json:"country"`
		CreatedAt           string `json:"created_at,omitempty"`
		Etag                string `json:"etag"`
		External            bool   `json:"external"`
		HasPreReleaseAccess bool   `json:"has_pre_release_access"`
		ID                  string `json:"id"`
		Name                string `json:"name"`
		State               string `json:"state"`
		UpdatedAt           string `json:"updated_at,omitempty"`
	} `json:"team"`
	TeamID string `json:"team_id"`
	Topic  struct {
		ComponentTypes         []string `json:"component_types"`
		ComponentTypesOptional []any    `json:"component_types_optional"`
		CreatedAt              string   `json:"created_at,omitempty"`
		Data                   struct {
			PullSecret struct {
				Auths struct {
					CloudOpenshiftCom struct {
						Auth  string `json:"auth"`
						Email string `json:"email"`
					} `json:"cloud.openshift.com"`
					QuayIo struct {
						Auth  string `json:"auth"`
						Email string `json:"email"`
					} `json:"quay.io"`
					RegistryCiOpenshiftOrg struct {
						Auth string `json:"auth"`
					} `json:"registry.ci.openshift.org"`
					RegistryConnectRedhatCom struct {
						Auth  string `json:"auth"`
						Email string `json:"email"`
					} `json:"registry.connect.redhat.com"`
					RegistryRedhatIo struct {
						Auth  string `json:"auth"`
						Email string `json:"email"`
					} `json:"registry.redhat.io"`
				} `json:"auths"`
			} `json:"pull_secret"`
		} `json:"data"`
		Etag          string `json:"etag"`
		ExportControl bool   `json:"export_control"`
		ID            string `json:"id"`
		Name          string `json:"name"`
		NextTopicID   string `json:"next_topic_id"`
		ProductID     string `json:"product_id"`
		State         string `json:"state"`
		UpdatedAt     string `json:"updated_at,omitempty"`
	} `json:"topic"`
	TopicID             string `json:"topic_id"`
	UpdatePreviousJobID any    `json:"update_previous_job_id"`
	UpdatedAt           string `json:"updated_at,omitempty"`
	URL                 string `json:"url"`
	UserAgent           string `json:"user_agent"`
}

type JobsResponse struct {
	Meta Meta  `json:"_meta"`
	Jobs []Job `json:"jobs"`
}

type TopicsResponse struct {
	Meta   Meta `json:"_meta,omitempty"`
	Topics []struct {
		ComponentTypes         []string `json:"component_types,omitempty"`
		ComponentTypesOptional []any    `json:"component_types_optional,omitempty"`
		CreatedAt              string   `json:"created_at,omitempty"`
		Data                   struct {
		} `json:"data,omitempty"`
		Etag          string  `json:"etag,omitempty"`
		ExportControl bool    `json:"export_control,omitempty"`
		ID            string  `json:"id,omitempty"`
		Name          string  `json:"name,omitempty"`
		NextTopic     any     `json:"next_topic,omitempty"`
		NextTopicID   any     `json:"next_topic_id,omitempty"`
		Product       Product `json:"product,omitempty"`
		ProductID     string  `json:"product_id,omitempty"`
		State         string  `json:"state,omitempty"`
		UpdatedAt     string  `json:"updated_at,omitempty"`
	} `json:"topics,omitempty"`
}
