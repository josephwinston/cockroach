// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Ben Darnell

package multiraft

import (
	"fmt"

	"github.com/cockroachdb/cockroach/util/log"
	"github.com/coreos/etcd/raft"
)

// init installs an adapter to use glog for all log messages from raft.
func init() {
	raft.SetLogger(&glogLogger{})
}

// *glogLogger implements the raft.Logger interface. Note that all methods
// must be defined on the pointer type rather than the value type because
// (at least in the go 1.4 compiler), methods on a value type called through
// an interface pointer go through an additional layer of indirection that
// appears on the stack, and would make all our stack frame offsets incorrect.
type glogLogger struct{}

func (*glogLogger) Debug(v ...interface{}) {
	if log.V(1) {
		log.InfoDepth(1, v...)
	}
}

func (*glogLogger) Debugf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if log.V(1) {
		log.InfoDepth(1, s)
	}
}

func (*glogLogger) Info(v ...interface{}) {
	log.InfoDepth(1, v...)
}

func (*glogLogger) Infof(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	log.InfoDepth(1, s)
}

func (*glogLogger) Warning(v ...interface{}) {
	log.WarningDepth(1, v...)
}

func (*glogLogger) Warningf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	log.WarningDepth(1, s)
}

func (*glogLogger) Error(v ...interface{}) {
	log.ErrorDepth(1, v...)
}

func (*glogLogger) Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	log.ErrorDepth(1, s)
}

func (*glogLogger) Fatal(v ...interface{}) {
	log.FatalDepth(1, v...)
}

func (*glogLogger) Fatalf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	log.FatalDepth(1, s)
}

func (*glogLogger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	log.ErrorDepth(1, s)
	panic(s)
}

func (*glogLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	log.ErrorDepth(1, s)
	panic(s)
}
