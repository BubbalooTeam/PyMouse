package models

type LastFMInformations struct {
	Username string `bson:"username,omitempty" json:"username,omitempty" default:"nil"`
}
