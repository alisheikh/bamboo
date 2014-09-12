package api


import (
	"github.com/QubitProducts/bamboo/configuration"
	eb "github.com/QubitProducts/bamboo/services/event_bus"
	"net/http"
	"io"
	"log"
	"encoding/json"
)

type EventSubscriptions struct {
	Conf *configuration.Configuration
	EventBus *eb.EventBus
}

func (sub *EventSubscriptions) Callback(w http.ResponseWriter, r *http.Request) {
	var event eb.MarathonEvent
	payload := make([]byte, r.ContentLength)
	r.Body.Read(payload)
	defer r.Body.Close()
	err := json.Unmarshal(payload, &event)
	if err != nil {
		log.Fatal("Unable to decode JSON Marathon Event request")
	}

	sub.EventBus.Publish(event)
	io.WriteString(w, "Got it!")
}