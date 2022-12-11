package model

type Access struct {
	UserID     int64    `json:"user_id,omitempty"`
	Code       string   `json:"code,omitempty"`
	Additional []string `json:"additional,omitempty"`

	UserName      string `json:"user_name,omitempty"`
	UserFirstName string `json:"user_first_name,omitempty"`
	UserLastName  string `json:"user_last_name,omitempty"`
}
