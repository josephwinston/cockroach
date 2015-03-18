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
// Author: Tobias Schottdorf (tobias.schottdorf@gmail.com)

// Package resource embeds into the Cockroach certain data such as web html
// and stylesheets.
package resource

// When changing files included here, you may add the `-debug` flag, which will
// avoid embedding any files (but instead passes them through from disk). This
// avoids having to recompile before testing changes.
//go:generate go-bindata -pkg resource -mode 0644 -modtime 1400000000 -o ./embedded.go ./ui/...
// `go vet` complains about the go file created by go-bindata, so we groom it.
//go:generate goimports -w embedded.go
