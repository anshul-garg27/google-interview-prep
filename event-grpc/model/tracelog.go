package model

import "time"

type TraceLogEvent struct {
	Id             string
	HostName       string
	ServiceName    string
	Timestamp      time.Time
	TimeTaken      int64
	Method         string
	URI            string
	Headers        JSONB `sql:"type:jsonb"`
	ResponseStatus int64
	RemoteAddress  string
	OnlyRequest    bool
	RequestBody    string
	ResponseBody   string
}
