package models

type ListenBrainzInformations struct {
	Username string `bson:"username,omitempty" json:"username,omitempty" default:"nil"`
}
