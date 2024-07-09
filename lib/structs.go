package lib

type JsonTNFInfo struct {
	ID         string `json:"id"`
	TNFVersion string `json:"tnf_version"`
}

type JobsJsonOutput struct {
	Jobs []JsonTNFInfo `json:"jobs"`
}

type JobsResponse struct {
	Meta struct {
		Count int `json:"count"`
	} `json:"_meta"`
	Jobs []struct {
		ClientVersion string `json:"client_version"`
		Comment       string `json:"comment"`
		Components    []struct {
			CanonicalProjectName string `json:"canonical_project_name"`
			CreatedAt            string `json:"created_at,omitempty"`
			Data                 struct {
				// Created     time.Time `json:"created,omitempty"`
				Digest      []string `json:"digest"`
				DisplayName string   `json:"display_name"`
				PullURL     string   `json:"pull_url"`
				Tags        []string `json:"tags"`
				UID         string   `json:"uid"`
				URL         string   `json:"url"`
				Version     string   `json:"version"`
			} `json:"data,omitempty"`
			DisplayName string   `json:"display_name"`
			Etag        string   `json:"etag"`
			ID          string   `json:"id"`
			Message     string   `json:"message"`
			Name        string   `json:"name"`
			ReleasedAt  string   `json:"released_at,omitempty"`
			State       string   `json:"state"`
			Tags        []string `json:"tags"`
			TeamID      any      `json:"team_id"`
			Title       string   `json:"title"`
			TopicID     string   `json:"topic_id"`
			Type        string   `json:"type"`
			UID         string   `json:"uid"`
			UpdatedAt   string   `json:"updated_at,omitempty"`
			URL         string   `json:"url"`
			Version     string   `json:"version"`
		} `json:"components"`
		Configuration string `json:"configuration"`
		CreatedAt     string `json:"created_at,omitempty"`
		Duration      int    `json:"duration"`
		Etag          string `json:"etag"`
		ID            string `json:"id"`
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
		RemoteciID   string   `json:"remoteci_id"`
		Results      []any    `json:"results"`
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
	} `json:"jobs"`
}

type TopicsResponse struct {
	Meta struct {
		Count int `json:"count,omitempty"`
	} `json:"_meta,omitempty"`
	Topics []struct {
		ComponentTypes         []string `json:"component_types,omitempty"`
		ComponentTypesOptional []any    `json:"component_types_optional,omitempty"`
		CreatedAt              string   `json:"created_at,omitempty"`
		Data                   struct {
		} `json:"data,omitempty"`
		Etag          string `json:"etag,omitempty"`
		ExportControl bool   `json:"export_control,omitempty"`
		ID            string `json:"id,omitempty"`
		Name          string `json:"name,omitempty"`
		NextTopic     any    `json:"next_topic,omitempty"`
		NextTopicID   any    `json:"next_topic_id,omitempty"`
		Product       struct {
			CreatedAt   string `json:"created_at,omitempty"`
			Description string `json:"description,omitempty"`
			Etag        string `json:"etag,omitempty"`
			ID          string `json:"id,omitempty"`
			Label       string `json:"label,omitempty"`
			Name        string `json:"name,omitempty"`
			State       string `json:"state,omitempty"`
			UpdatedAt   string `json:"updated_at,omitempty"`
		} `json:"product,omitempty"`
		ProductID string `json:"product_id,omitempty"`
		State     string `json:"state,omitempty"`
		UpdatedAt string `json:"updated_at,omitempty"`
		Data0     struct {
			PullSecret struct {
				Auths struct {
					CloudOpenshiftCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"cloud.openshift.com,omitempty"`
					QuayIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"quay.io,omitempty"`
					RegistryCiOpenshiftOrg struct {
						Auth string `json:"auth,omitempty"`
					} `json:"registry.ci.openshift.org,omitempty"`
					RegistryConnectRedhatCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.connect.redhat.com,omitempty"`
					RegistryRedhatIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.redhat.io,omitempty"`
				} `json:"auths,omitempty"`
			} `json:"pull_secret,omitempty"`
		} `json:"data,omitempty"`
		Data1 struct {
			PullSecret struct {
				Auths struct {
					CloudOpenshiftCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"cloud.openshift.com,omitempty"`
					QuayIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"quay.io,omitempty"`
					RegistryCiOpenshiftOrg struct {
						Auth string `json:"auth,omitempty"`
					} `json:"registry.ci.openshift.org,omitempty"`
					RegistryConnectRedhatCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.connect.redhat.com,omitempty"`
					RegistryRedhatIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.redhat.io,omitempty"`
				} `json:"auths,omitempty"`
			} `json:"pull_secret,omitempty"`
		} `json:"data,omitempty"`
		Data2 struct {
			PullSecret struct {
				Auths struct {
					CloudOpenshiftCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"cloud.openshift.com,omitempty"`
					QuayIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"quay.io,omitempty"`
					RegistryCiOpenshiftOrg struct {
						Auth string `json:"auth,omitempty"`
					} `json:"registry.ci.openshift.org,omitempty"`
					RegistryConnectRedhatCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.connect.redhat.com,omitempty"`
					RegistryRedhatIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.redhat.io,omitempty"`
				} `json:"auths,omitempty"`
			} `json:"pull_secret,omitempty"`
		} `json:"data,omitempty"`
		Data3 struct {
			CkiJobURL string `json:"cki_job_url,omitempty"`
		} `json:"data,omitempty"`
		Data4 struct {
			PullSecret struct {
				Auths struct {
					CloudOpenshiftCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"cloud.openshift.com,omitempty"`
					QuayIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"quay.io,omitempty"`
					RegistryCiOpenshiftOrg struct {
						Auth string `json:"auth,omitempty"`
					} `json:"registry.ci.openshift.org,omitempty"`
					RegistryConnectRedhatCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.connect.redhat.com,omitempty"`
					RegistryRedhatIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.redhat.io,omitempty"`
				} `json:"auths,omitempty"`
			} `json:"pull_secret,omitempty"`
		} `json:"data,omitempty"`
		Data5 struct {
			CkiJobURL string `json:"cki_job_url,omitempty"`
		} `json:"data,omitempty"`
		Data6 struct {
			Registry struct {
				Login    string `json:"login,omitempty"`
				Password string `json:"password,omitempty"`
			} `json:"registry,omitempty"`
			Releasename string `json:"releasename,omitempty"`
		} `json:"data,omitempty"`
		Data7 struct {
			PullSecret struct {
				Auths struct {
					CloudOpenshiftCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"cloud.openshift.com,omitempty"`
					QuayIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"quay.io,omitempty"`
					RegistryCiOpenshiftOrg struct {
						Auth string `json:"auth,omitempty"`
					} `json:"registry.ci.openshift.org,omitempty"`
					RegistryConnectRedhatCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.connect.redhat.com,omitempty"`
					RegistryRedhatIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.redhat.io,omitempty"`
				} `json:"auths,omitempty"`
			} `json:"pull_secret,omitempty"`
		} `json:"data,omitempty"`
		Data8 struct {
			Registry struct {
				Login    string `json:"login,omitempty"`
				Password string `json:"password,omitempty"`
			} `json:"registry,omitempty"`
		} `json:"data,omitempty"`
		Data9 struct {
			PullSecret struct {
				Auths struct {
					CloudOpenshiftCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"cloud.openshift.com,omitempty"`
					QuayIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"quay.io,omitempty"`
					RegistryCiOpenshiftOrg struct {
						Auth string `json:"auth,omitempty"`
					} `json:"registry.ci.openshift.org,omitempty"`
					RegistryConnectRedhatCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.connect.redhat.com,omitempty"`
					RegistryRedhatIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.redhat.io,omitempty"`
				} `json:"auths,omitempty"`
			} `json:"pull_secret,omitempty"`
		} `json:"data,omitempty"`
		Data10 struct {
			PullSecret struct {
				Auths struct {
					CloudOpenshiftCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"cloud.openshift.com,omitempty"`
					QuayIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"quay.io,omitempty"`
					RegistryCiOpenshiftOrg struct {
						Auth string `json:"auth,omitempty"`
					} `json:"registry.ci.openshift.org,omitempty"`
					RegistryConnectRedhatCom struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.connect.redhat.com,omitempty"`
					RegistryRedhatIo struct {
						Auth  string `json:"auth,omitempty"`
						Email string `json:"email,omitempty"`
					} `json:"registry.redhat.io,omitempty"`
				} `json:"auths,omitempty"`
			} `json:"pull_secret,omitempty"`
		} `json:"data,omitempty"`
	} `json:"topics,omitempty"`
}
