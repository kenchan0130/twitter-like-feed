package models

import (
	"time"
)

type Tweet struct {
	ID               string
	Text             string
	CreatedAt        time.Time
	AuthorID         string
	AuthorScreenName string
}
