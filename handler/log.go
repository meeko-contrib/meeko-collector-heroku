// Copyright (c) 2014 The meeko-collector-heroku AUTHORS
//
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package handler

type Logger interface {
	Infof(format string, v ...interface{})
	Warn(v ...interface{}) error
	Critical(v ...interface{}) error
}
