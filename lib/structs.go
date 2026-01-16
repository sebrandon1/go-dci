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

type ComponentsResponse struct {
	Meta       Meta         `json:"_meta,omitempty"`
	Components []Components `json:"components,omitempty"`
}

// IdentityResponse represents the response from the /api/v1/identity endpoint
type IdentityResponse struct {
	Identity Identity `json:"identity"`
}

// Identity represents the authenticated user/remoteci information
type Identity struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Email     string `json:"email,omitempty"`
	Etag      string `json:"etag,omitempty"`
	Fullname  string `json:"fullname,omitempty"`
	State     string `json:"state,omitempty"`
	TeamID    string `json:"team_id,omitempty"`
	TeamName  string `json:"team_name,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ComponentTypesResponse represents the response from the /api/v1/componenttypes endpoint
type ComponentTypesResponse struct {
	Meta           Meta            `json:"_meta,omitempty"`
	ComponentTypes []ComponentType `json:"componenttypes,omitempty"`
}

// ComponentType represents a single component type in DCI
type ComponentType struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Etag      string `json:"etag,omitempty"`
	State     string `json:"state,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Topic represents a single topic in DCI
type Topic struct {
	ID                     string   `json:"id,omitempty"`
	Name                   string   `json:"name,omitempty"`
	ComponentTypes         []string `json:"component_types,omitempty"`
	ComponentTypesOptional []string `json:"component_types_optional,omitempty"`
	ProductID              string   `json:"product_id,omitempty"`
	NextTopicID            string   `json:"next_topic_id,omitempty"`
	ExportControl          bool     `json:"export_control,omitempty"`
	State                  string   `json:"state,omitempty"`
	Etag                   string   `json:"etag,omitempty"`
	CreatedAt              string   `json:"created_at,omitempty"`
	UpdatedAt              string   `json:"updated_at,omitempty"`
	Product                Product  `json:"product,omitempty"`
}

// TopicResponse represents a single topic response from the API
type TopicResponse struct {
	Topic Topic `json:"topic"`
}

// CreateTopicRequest represents the request body for creating a new topic
type CreateTopicRequest struct {
	Name                   string   `json:"name"`
	ProductID              string   `json:"product_id"`
	ComponentTypes         []string `json:"component_types,omitempty"`
	ComponentTypesOptional []string `json:"component_types_optional,omitempty"`
	ExportControl          bool     `json:"export_control,omitempty"`
	NextTopicID            string   `json:"next_topic_id,omitempty"`
}

// UpdateTopicRequest represents the request body for updating a topic
type UpdateTopicRequest struct {
	Name                   string   `json:"name,omitempty"`
	ComponentTypes         []string `json:"component_types,omitempty"`
	ComponentTypesOptional []string `json:"component_types_optional,omitempty"`
	ExportControl          *bool    `json:"export_control,omitempty"`
	NextTopicID            string   `json:"next_topic_id,omitempty"`
	State                  string   `json:"state,omitempty"`
}

// ComponentTypeResponse represents a single component type response from the API
type ComponentTypeResponse struct {
	ComponentType ComponentType `json:"componenttype"`
}

// CreateComponentTypeRequest represents the request body for creating a component type
type CreateComponentTypeRequest struct {
	Name string `json:"name"`
}

// UpdateComponentTypeRequest represents the request body for updating a component type
type UpdateComponentTypeRequest struct {
	Name  string `json:"name,omitempty"`
	State string `json:"state,omitempty"`
}

// ComponentResponse represents a single component response from the API
type ComponentResponse struct {
	Component Components `json:"component"`
}

// CreateComponentRequest represents the request body for creating a new component
type CreateComponentRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	TopicID string `json:"topic_id"`
	Version string `json:"version,omitempty"`
	URL     string `json:"url,omitempty"`
	State   string `json:"state,omitempty"`
}

// UpdateComponentRequest represents the request body for updating a component
type UpdateComponentRequest struct {
	Name    string   `json:"name,omitempty"`
	State   string   `json:"state,omitempty"`
	URL     string   `json:"url,omitempty"`
	Version string   `json:"version,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// JobResponse represents a single job response from the API
type JobResponse struct {
	Job Job `json:"job"`
}

// UpdateJobRequest represents the request body for updating a job
type UpdateJobRequest struct {
	Comment string   `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// ScheduleJobRequest represents the request body for scheduling a job
type ScheduleJobRequest struct {
	TopicID string `json:"topic_id"`
}

// FilesResponse represents the response from getting files
type FilesResponse struct {
	Meta  Meta   `json:"_meta,omitempty"`
	Files []File `json:"files,omitempty"`
}

// File represents a file in DCI
type File struct {
	ID        string `json:"id"`
	JobID     string `json:"job_id"`
	Name      string `json:"name"`
	Mime      string `json:"mime"`
	Size      int64  `json:"size"`
	Etag      string `json:"etag,omitempty"`
	State     string `json:"state,omitempty"`
	TeamID    string `json:"team_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// POST request/response structs

// CreateJobRequest represents the request body for creating a new job
type CreateJobRequest struct {
	TopicID    string   `json:"topic_id"`
	Components []string `json:"components,omitempty"`
	Comment    string   `json:"comment,omitempty"`
}

// CreateJobResponse represents the response from creating a new job
type CreateJobResponse struct {
	Job Job `json:"job"`
}

// JobState represents the valid job states
type JobState string

const (
	JobStateNew     JobState = "new"
	JobStatePreRun  JobState = "pre-run"
	JobStateRunning JobState = "running"
	JobStatePostRun JobState = "post-run"
	JobStateSuccess JobState = "success"
	JobStateFailure JobState = "failure"
	JobStateKilled  JobState = "killed"
	JobStateError   JobState = "error"
)

// UpdateJobStateRequest represents the request body for updating a job's state
type UpdateJobStateRequest struct {
	JobID   string `json:"job_id"`
	Status  string `json:"status"`
	Comment string `json:"comment,omitempty"`
}

// JobStateResponse represents the response from updating a job's state
type JobStateResponse struct {
	JobState struct {
		ID        string `json:"id"`
		JobID     string `json:"job_id"`
		Status    string `json:"status"`
		Comment   string `json:"comment,omitempty"`
		CreatedAt string `json:"created_at,omitempty"`
	} `json:"jobstate"`
}

// JobStatesResponse represents the response from getting job states
type JobStatesResponse struct {
	Meta      Meta       `json:"_meta,omitempty"`
	JobStates []JobState2 `json:"jobstates,omitempty"`
}

// JobState2 represents a single job state entry (different from JobState which is a string type)
type JobState2 struct {
	ID        string `json:"id"`
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	Comment   string `json:"comment,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// UploadFileResponse represents the response from uploading a file
type UploadFileResponse struct {
	File struct {
		ID        string `json:"id"`
		JobID     string `json:"job_id"`
		Name      string `json:"name"`
		Mime      string `json:"mime"`
		Size      string `json:"size"`
		Etag      string `json:"etag,omitempty"`
		State     string `json:"state,omitempty"`
		TeamID    string `json:"team_id,omitempty"`
		CreatedAt string `json:"created_at,omitempty"`
		UpdatedAt string `json:"updated_at,omitempty"`
	} `json:"file"`
}

// RemoteCI represents a remote CI in DCI
type RemoteCI struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APISecret string `json:"api_secret,omitempty"`
	TeamID    string `json:"team_id,omitempty"`
	State     string `json:"state,omitempty"`
	Public    bool   `json:"public,omitempty"`
	Etag      string `json:"etag,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// RemoteCIsResponse represents the response from getting remote CIs
type RemoteCIsResponse struct {
	Meta      Meta       `json:"_meta,omitempty"`
	RemoteCIs []RemoteCI `json:"remotecis,omitempty"`
}

// RemoteCIResponse represents a single remote CI response
type RemoteCIResponse struct {
	RemoteCI RemoteCI `json:"remoteci"`
}

// CreateRemoteCIRequest represents the request body for creating a remote CI
type CreateRemoteCIRequest struct {
	Name   string `json:"name"`
	TeamID string `json:"team_id"`
}

// UpdateRemoteCIRequest represents the request body for updating a remote CI
type UpdateRemoteCIRequest struct {
	Name  string `json:"name,omitempty"`
	State string `json:"state,omitempty"`
}

// Team represents a team in DCI
type Team struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Country             string `json:"country,omitempty"`
	External            bool   `json:"external,omitempty"`
	HasPreReleaseAccess bool   `json:"has_pre_release_access,omitempty"`
	State               string `json:"state,omitempty"`
	Etag                string `json:"etag,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
	UpdatedAt           string `json:"updated_at,omitempty"`
}

// TeamsResponse represents the response from getting teams
type TeamsResponse struct {
	Meta  Meta   `json:"_meta,omitempty"`
	Teams []Team `json:"teams,omitempty"`
}

// TeamResponse represents a single team response
type TeamResponse struct {
	Team Team `json:"team"`
}

// CreateTeamRequest represents the request body for creating a team
type CreateTeamRequest struct {
	Name    string `json:"name"`
	Country string `json:"country,omitempty"`
}

// UpdateTeamRequest represents the request body for updating a team
type UpdateTeamRequest struct {
	Name    string `json:"name,omitempty"`
	Country string `json:"country,omitempty"`
	State   string `json:"state,omitempty"`
}

// User represents a user in DCI
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Fullname  string `json:"fullname,omitempty"`
	Email     string `json:"email,omitempty"`
	TeamID    string `json:"team_id,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	State     string `json:"state,omitempty"`
	Etag      string `json:"etag,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// UsersResponse represents the response from getting users
type UsersResponse struct {
	Meta  Meta   `json:"_meta,omitempty"`
	Users []User `json:"users,omitempty"`
}

// UserResponse represents a single user response
type UserResponse struct {
	User User `json:"user"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
	TeamID   string `json:"team_id"`
	Password string `json:"password"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	State    string `json:"state,omitempty"`
}

// ProductsResponse represents the response from getting products
type ProductsResponse struct {
	Meta     Meta      `json:"_meta,omitempty"`
	Products []Product `json:"products,omitempty"`
}

// ProductResponse represents a single product response
type ProductResponse struct {
	Product Product `json:"product"`
}
