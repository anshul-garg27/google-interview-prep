package event

import log "github.com/sirupsen/logrus"

type EventTarget struct {
	TargetType string
	Properties map[string]string
}

type Event struct {
	EventType  string
	Target     EventTarget
	Headers    map[string]interface{}
	Properties map[string]interface{}
	Payload    interface{}
}

func (e *Event) Apply() {
	log.Error("Applying event - {}", e)
}
