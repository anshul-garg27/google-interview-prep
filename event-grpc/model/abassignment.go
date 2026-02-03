package model

import (
	"time"
)

type AbAssignment struct {
	Timestamp   time.Time
	UserId      int
	BBDeviceId  string
	NewUser     int
	Checksum    string
	Ver         int64
	Assignments JSONB
}
