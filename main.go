// Copyright (c) 2014 The meeko-collector-heroku AUTHORS
//
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"github.com/meeko-contrib/meeko-collector-heroku/handler"

	"github.com/meeko-contrib/go-meeko-webhook-receiver/receiver"
	"github.com/meeko/go-meeko/agent"
)

func main() {
	var (
		logger = agent.Logging()
		pubsub = agent.PubSub()
	)
	receiver.ListenAndServe(&handler.WebhookHandler{
		logger,
		pubsub.Publish,
	})
}
