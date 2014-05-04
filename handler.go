// Copyright (c) 2014 The cider-collector-heroku AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/cider/go-cider/cider/services/logging"
)

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

type HerokuWebhookHandler struct {
	logger  *logging.Service
	forward func(eventType string, eventObject interface{}) error
}

func (handler *HerokuWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		handler.logger.Warn("POST from %v: %v", r.RemoteAddr, err)
		return
	}

	// Publish the event.
	if err := handler.forward("heroku.deploy", event); err != nil {
		http.Error(w, "Event Not Published", http.StatusInternalServerError)
		handler.logger.Critical(err)
		return
	}

	handler.logger.Infof("POST from %v: Forwarding heroku.deploy", r.URL)
	w.WriteHeader(http.StatusAccepted)
}
