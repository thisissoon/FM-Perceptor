// Events Package Constants

package events

// Redis Keys
const (
	eventsChannel   string = "fm:events"
	currentTrackKey string = "fm:player:current"
	startTimeKey    string = "fm:player:start_time"
)

// Event payload to publish
type publishEventPayload struct {
	event string `json:"event"`
	track string `json:"uri"`
	user  string `json:"user"`
}
