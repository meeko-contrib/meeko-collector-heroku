// Copyright (c) 2014 The meeko-collector-heroku AUTHORS
//
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package handler

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const EventType = "heroku.deployment"

const (
	statusUnprocessableEntity = 422
	maxBodySize               = int64(10 << 20)
)

type HerokuEvent struct {
	App      string `codec:"app"`
	User     string `codec:"user"`
	URL      string `codec:"url"`
	Head     string `codec:"head"`
	HeadLong string `codec:"head_long"`
	GitLog   string `codec:"git_log"`
}

func (event *HerokuEvent) Validate() error {
	ev := reflect.Indirect(reflect.ValueOf(event))
	et := ev.Type()
	for i := 0; i < et.NumField(); i++ {
		fv := ev.Field(i)
		ft := et.Field(i)
		if fv.Interface().(string) == "" {
			variable := strings.Split(string(ft.Tag), "\"")[1]
			return fmt.Errorf("%v variable is not set", variable)
		}
	}
	return nil
}

type WebhookHandler struct {
	Logger  Logger
	Forward func(eventType string, eventObject interface{}) error
}

func (handler *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the event object.
	event := &HerokuEvent{
		App:      r.FormValue("app"),
		User:     r.FormValue("user"),
		URL:      r.FormValue("url"),
		Head:     r.FormValue("head"),
		HeadLong: r.FormValue("head_long"),
		GitLog:   r.FormValue("git_log"),
	}
	if err := event.Validate(); err != nil {
		http.Error(w, err.Error(), statusUnprocessableEntity)
		handler.Logger.Warn("POST from %v: %v", r.RemoteAddr, err)
		return
	}

	// Publish the event.
	if err := handler.Forward(EventType, event); err != nil {
		http.Error(w, "Event Not Published", http.StatusInternalServerError)
		handler.Logger.Critical(err)
		return
	}

	handler.Logger.Infof("POST from %v: Forwarding %v", r.URL, EventType)
	w.WriteHeader(http.StatusAccepted)
}
