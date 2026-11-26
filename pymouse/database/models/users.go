package models

type UsersInformations struct {
	UserID       int64                    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserName     string                   `bson:"username" json:"username" default:"nil"`
	FirstName    string                   `bson:"first_name" json:"first_name" default:"nil"`
	Language     string                   `bson:"language" json:"language" default:"en_us"`
	Away         AwayInformations         `bson:"away,omitempty" json:"away,omitempty"`
	ListenBrainz ListenBrainzInformations `bson:"listenbrainz,omitempty" json:"listenbrainz,omitempty"`
	LastFM       LastFMInformations       `bson:"lastfm,omitempty" json:"lastfm,omitempty"`
	StatsFM      StatsFMInformations      `bson:"statsfm,omitempty" json:"statsfm,omitempty"`
}
