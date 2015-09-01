// REST API Handlers for /events

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	v "github.com/asaskevich/govalidator"
	"github.com/zenazn/goji/web"
	"gopkg.in/redis.v3"
)

const (
	CURRENT_KEY    string = "fm:player:current"
	START_TIME_KEY string = "fm:player:start_time"
)

type playEvent struct {
	Start string `json:"start" valid:"iso8601,required"`
	Uri   string `json:"uri" valid:"required"`
	User  string `json:"user" valid:"required"`
}

// Handle sending a play event (POST /events/play)
func playHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	event := &playEvent{}

	err := decoder.Decode(&event)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Custom Validator
	v.TagMap["iso8601"] = v.Validator(func(str string) bool {
		_, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return false
		}

		return true
	})

	// Validate
	res, err := v.ValidateStruct(event)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, http.StatusText(422), 422)
		fmt.Println(res)
		return
	}

	// Redis Client
	red := c.Env["REDIS"].(*redis.Client)

	// Save Current Track
	err = red.Set(CURRENT_KEY, string(event.Uri), 0).Err()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Save Start Time
	err = red.Set(START_TIME_KEY, string(event.Start), 0).Err()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// We got to the end - everything went fine!
	w.WriteHeader(201)
}

// Handle sending a end event (POST /events/end)
func endHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("end"))
}

// Handle sendiong a volume change event (POST /events/volume)
func volumeHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("vol"))
}

// Handle sendiong a mute change event (POST /events/mute)
func muteHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("mute"))
}
