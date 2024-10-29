package modeldb

import "time"

type AwayInformations struct {
	IsAway     bool      `bson:"is_away" json:"is_away" default:"false"`
	AwayTime   time.Time `bson:"away_time" json:"away_time"`
	AwayReason string    `bson:"away_reason,omitempty" json:"away_reason,omitempty" default:"nil"`
}
