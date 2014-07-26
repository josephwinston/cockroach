// Copyright 2014 The Cockroach Authors.
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
// implied.  See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Levon Lloyd (levon.lloyd@gmail.com)
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package storage

import (
	"reflect"
	"testing"
)

var testConfig = ZoneConfig{
	Replicas: []Attributes{
		Attributes([]string{"a", "ssd"}),
		Attributes([]string{"a", "hdd"}),
		Attributes([]string{"b", "ssd"}),
		Attributes([]string{"b", "hdd"}),
	},
	RangeMinBytes: 1 << 20,
	RangeMaxBytes: 64 << 20,
}

var yamlConfig = `
replicas:
  - [a, ssd]
  - [a, hdd]
  - [b, ssd]
  - [b, hdd]
range_min_bytes: 1048576
range_max_bytes: 67108864
`

func TestZoneConfigRoundTrip(t *testing.T) {
	yaml, err := testConfig.ToYAML()
	if err != nil {
		t.Errorf("failed converting to yaml: %v", err)
	}
	parsedZoneConfig, err := ParseZoneConfig(yaml)
	if err != nil {
		t.Errorf("failed parsing config: %v", err)
	}
	if !reflect.DeepEqual(testConfig, *parsedZoneConfig) {
		t.Errorf("yaml round trip configs differ.\nOriginal: %+v\nParse: %+v\n", testConfig, parsedZoneConfig)
	}
}

func TestZoneConfigYAMLParsing(t *testing.T) {
	parsedZoneConfig, err := ParseZoneConfig([]byte(yamlConfig))
	if err != nil {
		t.Errorf("failed parsing config: %v", err)
	}
	if !reflect.DeepEqual(testConfig, *parsedZoneConfig) {
		t.Errorf("yaml round trip configs differ.\nOriginal: %+v\nParse: %+v\n", testConfig, parsedZoneConfig)
	}
}

func TestIsSubset(t *testing.T) {
	a := Attributes([]string{"a", "b", "c"})
	b := Attributes([]string{"a", "b"})
	c := Attributes([]string{"a"})
	if !b.IsSubset(a) {
		t.Errorf("expected %+v to be a subset of %+v", b, a)
	}
	if !c.IsSubset(a) {
		t.Errorf("expected %+v to be a subset of %+v", c, a)
	}
	if !c.IsSubset(b) {
		t.Errorf("expected %+v to be a subset of %+v", c, b)
	}
	if a.IsSubset(b) {
		t.Errorf("%+v should not be a subset of %+v", a, b)
	}
	if a.IsSubset(c) {
		t.Errorf("%+v should not be a subset of %+v", a, c)
	}
	if b.IsSubset(c) {
		t.Errorf("%+v should not be a subset of %+v", b, c)
	}
}

func TestSortedString(t *testing.T) {
	a := Attributes([]string{"a", "b", "c"})
	if a.SortedString() != "a,b,c" {
		t.Errorf("sorted string of %+v (%s) != \"a,b,c\"", a, a.SortedString())
	}
	b := Attributes([]string{"c", "a", "b"})
	if b.SortedString() != "a,b,c" {
		t.Errorf("sorted string of %+v (%s) != \"a,b,c\"", b, b.SortedString())
	}
	// Duplicates.
	c := Attributes([]string{"c", "c", "a", "a", "b", "b"})
	if c.SortedString() != "a,b,c" {
		t.Errorf("sorted string of %+v (%s) != \"a,b,c\"", c, c.SortedString())
	}
}

func TestPermConfig(t *testing.T) {
	p := &PermConfig{
		Read:  []string{"foo", "bar", "baz"},
		Write: []string{"foo", "baz"},
	}
	for _, u := range p.Read {
		if !p.CanRead(u) {
			t.Errorf("expected read permission for %q", u)
		}
	}
	if p.CanRead("bad") {
		t.Errorf("unexpected read access for user \"bad\"")
	}
	for _, u := range p.Write {
		if !p.CanWrite(u) {
			t.Errorf("expected read permission for %q", u)
		}
	}
	if p.CanWrite("bar") {
		t.Errorf("unexpected read access for user \"bar\"")
	}
}
