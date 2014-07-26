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
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package storage

import (
	"bytes"
	"sort"
	"testing"

	"github.com/golang/glog"
)

const (
	config1 = 1
	config2 = 2
	config3 = 3
	config4 = 4
)

func buildTestPrefixConfigMap() PrefixConfigMap {
	configs := []*PrefixConfig{
		{KeyMin, nil, config1},
		{Key("/db1"), nil, config2},
		{Key("/db1/table"), nil, config3},
		{Key("/db3"), nil, config4},
	}
	pcc, err := NewPrefixConfigMap(configs)
	if err != nil {
		glog.Fatalf("unexpected error building config map: %v", err)
	}
	return pcc
}

// TestPrefixEndKey verifies the end keys on prefixes.
func TestPrefixEndKey(t *testing.T) {
	testData := []struct {
		prefix, expEnd Key
	}{
		{KeyMin, KeyMax},
		{Key("0"), Key("1")},
		{Key("a"), Key("b")},
		{Key("db0"), Key("db1")},
		{Key("\xfe"), Key("\xff")},
		{KeyMax, KeyMax},
		{Key("\xff\xff"), Key("\xff\xff")},
	}

	for i, test := range testData {
		if bytes.Compare(PrefixEndKey(test.prefix), test.expEnd) != 0 {
			t.Errorf("%d: %q end key %q != %q", i, test.prefix, PrefixEndKey(test.prefix), test.expEnd)
		}
	}
}

// TestPrefixConfigSort verifies sorting of keys.
func TestPrefixConfigSort(t *testing.T) {
	keys := []Key{
		KeyMax,
		Key("c"),
		Key("a"),
		Key("b"),
		Key("aa"),
		Key("\xfe"),
		KeyMin,
	}
	expKeys := []Key{
		KeyMin,
		Key("a"),
		Key("aa"),
		Key("b"),
		Key("c"),
		Key("\xfe"),
		KeyMax,
	}
	pcc := PrefixConfigMap{}
	for _, key := range keys {
		pcc = append(pcc, &PrefixConfig{key, nil, nil})
	}
	sort.Sort(pcc)
	for i, pc := range pcc {
		if bytes.Compare(pc.Prefix, expKeys[i]) != 0 {
			t.Errorf("order for index %d incorrect; expected %q, got %q", i, expKeys[i], pc.Prefix)
		}
	}
}

// TestPrefixConfigBuild adds prefixes and verifies they're
// sorted and proper end keys are generated.
func TestPrefixConfigBuild(t *testing.T) {
	expPrefixConfigs := []PrefixConfig{
		{KeyMin, nil, config1},
		{Key("/db1"), nil, config2},
		{Key("/db1/table"), nil, config3},
		{Key("/db1/tablf"), nil, config2},
		{Key("/db2"), nil, config1},
		{Key("/db3"), nil, config4},
		{Key("/db4"), nil, config1},
	}
	pcc := buildTestPrefixConfigMap()
	if len(pcc) != len(expPrefixConfigs) {
		t.Fatalf("incorrect number of built prefix configs; expected %d, got %d",
			len(expPrefixConfigs), len(pcc))
	}
	for i, pc := range pcc {
		exp := expPrefixConfigs[i]
		if bytes.Compare(pc.Prefix, exp.Prefix) != 0 {
			t.Errorf("prefix for index %d incorrect; expected %q, got %q", i, exp.Prefix, pc.Prefix)
		}
		if pc.Config != exp.Config {
			t.Errorf("config for index %d incorrect: expected %v, got %v", i, exp.Config, pc.Config)
		}
	}
}

// TestMatchByPrefix verifies matching on longest prefix.
func TestMatchByPrefix(t *testing.T) {
	pcc := buildTestPrefixConfigMap()
	testData := []struct {
		key       Key
		expConfig interface{}
	}{
		{KeyMin, config1},
		{Key("\x01"), config1},
		{Key("/db"), config1},
		{Key("/db1"), config2},
		{Key("/db1/a"), config2},
		{Key("/db1/table1"), config3},
		{Key("/db1/table\xff"), config3},
		{Key("/db2"), config1},
		{Key("/db3"), config4},
		{Key("/db3\xff"), config4},
		{Key("/db5"), config1},
		{Key("/xfe"), config1},
		{Key("/xff"), config1},
	}
	for i, test := range testData {
		pc := pcc.MatchByPrefix(test.key)
		if test.expConfig != pc.Config {
			t.Errorf("%d: expected config %v for %q; got %v", i, test.expConfig, test.key, pc.Config)
		}
	}
}

// TestesMatchesByPrefix verifies all matching prefixes.
func TestMatchesByPrefix(t *testing.T) {
	pcc := buildTestPrefixConfigMap()
	testData := []struct {
		key        Key
		expConfigs []interface{}
	}{
		{KeyMin, []interface{}{config1}},
		{Key("\x01"), []interface{}{config1}},
		{Key("/db"), []interface{}{config1}},
		{Key("/db1"), []interface{}{config2, config1}},
		{Key("/db1/a"), []interface{}{config2, config1}},
		{Key("/db1/table1"), []interface{}{config3, config2, config1}},
		{Key("/db1/table\xff"), []interface{}{config3, config2, config1}},
		{Key("/db2"), []interface{}{config1}},
		{Key("/db3"), []interface{}{config4, config1}},
		{Key("/db3\xff"), []interface{}{config4, config1}},
		{Key("/db5"), []interface{}{config1}},
		{Key("/xfe"), []interface{}{config1}},
		{Key("/xff"), []interface{}{config1}},
	}
	for i, test := range testData {
		pcs := pcc.MatchesByPrefix(test.key)
		if len(pcs) != len(test.expConfigs) {
			t.Errorf("%d: expected %d matches, got %d", i, len(test.expConfigs), len(pcs))
			continue
		}
		for j, pc := range pcs {
			if pc.Config != test.expConfigs[j] {
				t.Errorf("%d: expected \"%d\"th config %v for %q; got %v", i, j, test.expConfigs[j], test.key, pc.Config)
			}
		}
	}
}

// TestSplitRangeByPrefixesErrors verifies various error conditions
// for splitting ranges.
func TestSplitRangeByPrefixesError(t *testing.T) {
	pcc, err := NewPrefixConfigMap([]*PrefixConfig{})
	if err == nil {
		t.Error("expected error building config map with no default prefix")
	}
	pcc = buildTestPrefixConfigMap()
	// Key order is reversed.
	if _, err := pcc.SplitRangeByPrefixes(KeyMax, KeyMin); err == nil {
		t.Error("expected error with reversed keys")
	}
	// Same start and end keys.
	if _, err := pcc.SplitRangeByPrefixes(KeyMin, KeyMin); err != nil {
		t.Error("unexpected error with same start & end keys")
	}
}

// TestSplitRangeByPrefixes verifies splitting of a key range
// into sub-ranges based on config prefixes.
func TestSplitRangeByPrefixes(t *testing.T) {
	pcc := buildTestPrefixConfigMap()
	testData := []struct {
		start, end Key
		expRanges  []*RangeResult
	}{
		// The full range.
		{KeyMin, KeyMax, []*RangeResult{
			{KeyMin, Key("/db1"), config1},
			{Key("/db1"), Key("/db1/table"), config2},
			{Key("/db1/table"), Key("/db1/tablf"), config3},
			{Key("/db1/tablf"), Key("/db2"), config2},
			{Key("/db2"), Key("/db3"), config1},
			{Key("/db3"), Key("/db4"), config4},
			{Key("/db4"), KeyMax, config1},
		}},
		// A subrange containing all databases.
		{Key("/db"), Key("/dc"), []*RangeResult{
			{Key("/db"), Key("/db1"), config1},
			{Key("/db1"), Key("/db1/table"), config2},
			{Key("/db1/table"), Key("/db1/tablf"), config3},
			{Key("/db1/tablf"), Key("/db2"), config2},
			{Key("/db2"), Key("/db3"), config1},
			{Key("/db3"), Key("/db4"), config4},
			{Key("/db4"), Key("/dc"), config1},
		}},
		// A subrange spanning from arbitrary points within zones.
		{Key("/db1/a"), Key("/db3/b"), []*RangeResult{
			{Key("/db1/a"), Key("/db1/table"), config2},
			{Key("/db1/table"), Key("/db1/tablf"), config3},
			{Key("/db1/tablf"), Key("/db2"), config2},
			{Key("/db2"), Key("/db3"), config1},
			{Key("/db3"), Key("/db3/b"), config4},
		}},
		// A subrange containing only /db1.
		{Key("/db1"), Key("/db2"), []*RangeResult{
			{Key("/db1"), Key("/db1/table"), config2},
			{Key("/db1/table"), Key("/db1/tablf"), config3},
			{Key("/db1/tablf"), Key("/db2"), config2},
		}},
		// A subrange containing only /db1/table.
		{Key("/db1/table"), Key("/db1/tablf"), []*RangeResult{
			{Key("/db1/table"), Key("/db1/tablf"), config3},
		}},
		// A subrange within /db1/table.
		{Key("/db1/table3"), Key("/db1/table4"), []*RangeResult{
			{Key("/db1/table3"), Key("/db1/table4"), config3},
		}},
	}
	for i, test := range testData {
		results, err := pcc.SplitRangeByPrefixes(test.start, test.end)
		if err != nil {
			t.Errorf("%d: unexpected error splitting ranges: %v", i, err)
		}
		if len(results) != len(test.expRanges) {
			t.Errorf("%d: expected %d matches, got %d:", i, len(test.expRanges), len(results))
			for j, r := range results {
				t.Errorf("match %d: %q, %q, %v", j, r.start, r.end, r.config)
			}
			continue
		}
		for j, r := range results {
			if bytes.Compare(r.start, test.expRanges[j].start) != 0 {
				t.Errorf("%d: expected \"%d\"th range start key %q, got %q", i, j, test.expRanges[j].start, r.start)
			}
			if bytes.Compare(r.end, test.expRanges[j].end) != 0 {
				t.Errorf("%d: expected \"%d\"th range end key %q, got %q", i, j, test.expRanges[j].end, r.end)
			}
			if r.config != test.expRanges[j].config {
				t.Errorf("%d: expected \"%d\"th range config %v, got %v", i, j, test.expRanges[j].config, r.config)
			}
		}
	}
}
