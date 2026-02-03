package jobtracker

import (
	"coffee/core/client"
)

type JobResponse struct {
	Status client.SFStatus `json:"status"`
	Data   []JobEntry      `json:"data"`
}
