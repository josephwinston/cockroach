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
// Author: Tobias Schottdorf (tobias.schottdorf@gmail.com)

package storage

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/util"
)

func ensureRangeEqual(t *testing.T, sortedKeys []string, keyMap map[string][]byte, keyvals []KeyValue) {
	if len(keyvals) != len(sortedKeys) {
		t.Errorf("length mismatch. expected %s, got %s", sortedKeys, keyvals)
	}
	t.Log("---")
	for i, kv := range keyvals {
		t.Logf("index: %d\tk: %q\tv: %q\n", i, kv.Key, kv.Value.Bytes)
		if sortedKeys[i] != string(kv.Key) {
			t.Errorf("key mismatch at index %d: expected %q, got %q", i, sortedKeys[i], kv.Key)
		}
		if !bytes.Equal(keyMap[sortedKeys[i]], kv.Value.Bytes) {
			t.Errorf("value mismatch at index %d: expected %q, got %q", i, keyMap[sortedKeys[i]], kv.Value.Bytes)
		}
	}
}

// runWithAllEngines creates a new engine of each supported type and
// invokes the supplied test func with each instance.
func runWithAllEngines(test func(e Engine, t *testing.T), t *testing.T) {
	inMem := NewInMem(Attributes{}, 10<<20)

	loc := fmt.Sprintf("%s/data_%d", os.TempDir(), time.Now().UnixNano())
	rocksdb, err := NewRocksDB(Attributes([]string{"ssd"}), loc)
	if err != nil {
		t.Fatalf("could not create new rocksdb db instance at %s: %v", loc, err)
	}
	defer func(t *testing.T) {
		rocksdb.close()
		if err := rocksdb.destroy(); err != nil {
			t.Errorf("could not delete rocksdb db at %s: %v", loc, err)
		}
	}(t)

	test(inMem, t)
	test(rocksdb, t)
}

// TestEngineWriteBatch writes a batch containing 50K rows (all the
// same key) and concurrently attempts to read the value in a tight
// loop. The test verifies that either there is no value for the key
// or it contains the final value, but never a value in between.
func TestEngineWriteBatch(t *testing.T) {
	numWrites := 50000
	key := Key("a")
	finalVal := []byte(strconv.Itoa(numWrites - 1))

	runWithAllEngines(func(e Engine, t *testing.T) {
		// Start a concurrent read operation in a busy loop.
		readsBegun := make(chan struct{})
		readsDone := make(chan struct{})
		writesDone := make(chan struct{})
		go func() {
			for i := 0; ; i++ {
				select {
				case <-writesDone:
					close(readsDone)
					return
				default:
					val, err := e.get(key)
					if err != nil {
						t.Fatal(err)
					}
					if val.Bytes != nil && bytes.Compare(val.Bytes, finalVal) != 0 {
						close(readsDone)
						t.Fatalf("key value should be empty or %q; got %q", string(finalVal), string(val.Bytes))
					}
					if i == 0 {
						close(readsBegun)
					}
				}
			}
		}()
		// Wait until we've succeeded with first read.
		<-readsBegun

		// Create key/values and put them in a batch to engine.
		puts := make([]interface{}, numWrites, numWrites)
		for i := 0; i < numWrites; i++ {
			puts[i] = BatchPut{Key: key, Value: Value{Bytes: []byte(strconv.Itoa(i))}}
		}
		if err := e.writeBatch(puts); err != nil {
			t.Fatal(err)
		}
		close(writesDone)
		<-readsDone
	}, t)
}

func TestEngineBatch(t *testing.T) {
	runWithAllEngines(func(engine Engine, t *testing.T) {
		numShuffles := 100
		key := Key("a")
		// Those are randomized below.
		batch := []interface{}{
			BatchPut{Key: key, Value: Value{Bytes: []byte("~ockroachDB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("C~ckroachDB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("Co~kroachDB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("Coc~roachDB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("Cock~oachDB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("Cockr~achDB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("Cockro~chDB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("Cockroa~hDB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("Cockroac~DB")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("Cockroach~B")}},
			BatchPut{Key: key, Value: Value{Bytes: []byte("CockroachD~")}},
			BatchDelete(key),
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender("C"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender(" o"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender("  c"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender(" k"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender("r"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender(" o"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender("  a"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender(" c"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender("h"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender(" D"))}},
			BatchMerge{Key: key, Value: Value{Bytes: util.GobEncodeOrDie(Appender("  B"))}},
		}

		for i := 0; i < numShuffles; i++ {
			// In each run, create an array of shuffled operations.
			shuffledIndices := rand.Perm(len(batch))
			currentBatch := make([]interface{}, len(batch))
			for k := range currentBatch {
				currentBatch[k] = batch[shuffledIndices[k]]
			}
			// Reset the key
			engine.del(key)
			// Run it once with individual operations and remember the result.
			for _, op := range currentBatch {
				if err := engine.writeBatch([]interface{}{op}); err != nil {
					t.Errorf("batch test: %d: op %v: %v", i, op, err)
					continue
				}
			}
			correctValue, _ := engine.get(key)
			// Run the whole thing as a batch and compare.
			if err := engine.writeBatch(currentBatch); err != nil {
				t.Errorf("batch test: %d: %v", i, err)
				continue
			}
			actualValue, _ := engine.get(key)
			if !bytes.Equal(actualValue.Bytes, correctValue.Bytes) {
				t.Errorf("batch test: %d: result inconsistent")
			}
		}
	}, t)
}

func TestEnginePutGetDelete(t *testing.T) {
	runWithAllEngines(func(engine Engine, t *testing.T) {
		// Test for correct handling of empty keys, which should produce errors.
		for _, err := range []error{
			engine.put([]byte(""), Value{}),
			engine.put(nil, Value{}),
			func() error {
				_, err := engine.get([]byte(""))
				return err
			}(),
			engine.del(nil),
			func() error {
				_, err := engine.get(nil)
				return err
			}(),
			engine.del(nil),
			engine.del([]byte("")),
		} {
			if err == nil {
				t.Fatalf("illegal handling of empty key")
			}
		}

		// Test for allowed keys, which should go through.
		testCases := []struct {
			key, value []byte
		}{
			{[]byte("dog"), []byte("woof")},
			{[]byte("cat"), []byte("meow")},
			{[]byte("emptyval"), nil},
			{[]byte("emptyval2"), []byte("")},
			{[]byte("server"), []byte("42")},
		}
		for _, c := range testCases {
			val, err := engine.get(c.key)
			if err != nil {
				t.Errorf("get: expected no error, but got %s", err)
			}
			if len(val.Bytes) != 0 {
				t.Errorf("expected key %q value.Bytes to be nil: got %+v", c.key, val)
			}
			if err := engine.put(c.key, Value{Bytes: c.value}); err != nil {
				t.Errorf("put: expected no error, but got %s", err)
			}
			val, err = engine.get(c.key)
			if err != nil {
				t.Errorf("get: expected no error, but got %s", err)
			}
			if !bytes.Equal(val.Bytes, c.value) {
				t.Errorf("expected key value %s to be %+v: got %+v", c.key, c.value, val)
			}
			if err := engine.del(c.key); err != nil {
				t.Errorf("delete: expected no error, but got %s", err)
			}
			val, err = engine.get(c.key)
			if err != nil {
				t.Errorf("get: expected no error, but got %s", err)
			}
			if len(val.Bytes) != 0 {
				t.Errorf("expected key %s value.Bytes to be nil: got %+v", c.key, val)
			}
		}
	}, t)
}

// TestEngineMerge tests that the passing through of engine merge operations
// to the goMerge function works as expected. The semantics are tested more
// exhaustively in the merge tests themselves.
func TestEngineMerge(t *testing.T) {
	runWithAllEngines(func(engine Engine, t *testing.T) {
		testKey := Key("haste not in life")
		merges := []Value{
			Value{Bytes: []byte(util.GobEncodeOrDie(Appender("x")))},
			Value{Bytes: []byte(util.GobEncodeOrDie(Appender("y")))},
			Value{Bytes: []byte(util.GobEncodeOrDie(Appender("z")))},
		}
		for i, update := range merges {
			if err := engine.merge(testKey, update); err != nil {
				t.Fatalf("%d: %v", i, err)
			}
		}
		result, _ := engine.get(testKey)
		if !bytes.Equal(util.GobDecodeOrDie(result.Bytes).(Appender), Appender("xyz")) {
			t.Errorf("unexpected append-merge result")
		}
	}, t)
}

func TestEngineScan1(t *testing.T) {
	runWithAllEngines(func(engine Engine, t *testing.T) {
		testCases := []struct {
			key, value []byte
		}{
			{[]byte("dog"), []byte("woof")},
			{[]byte("cat"), []byte("meow")},
			{[]byte("server"), []byte("42")},
			{[]byte("french"), []byte("Allô?")},
			{[]byte("german"), []byte("hallo")},
			{[]byte("chinese"), []byte("你好")},
		}
		keyMap := map[string][]byte{}
		for _, c := range testCases {
			if err := engine.put(c.key, Value{Bytes: c.value}); err != nil {
				t.Errorf("could not put key %q: %v", c.key, err)
			}
			keyMap[string(c.key)] = c.value
		}
		sortedKeys := make([]string, len(testCases))
		for i, t := range testCases {
			sortedKeys[i] = string(t.key)
		}
		sort.Strings(sortedKeys)

		keyvals, err := engine.scan([]byte("chinese"), []byte("german"), 0)
		if err != nil {
			t.Fatalf("could not run scan: %v", err)
		}
		ensureRangeEqual(t, sortedKeys[1:4], keyMap, keyvals)

		// Check an end of range which does not equal an existing key.
		keyvals, err = engine.scan([]byte("chinese"), []byte("german1"), 0)
		if err != nil {
			t.Fatalf("could not run scan: %v", err)
		}
		ensureRangeEqual(t, sortedKeys[1:5], keyMap, keyvals)

		keyvals, err = engine.scan([]byte("chinese"), []byte("german"), 2)
		if err != nil {
			t.Fatalf("could not run scan: %v", err)
		}
		ensureRangeEqual(t, sortedKeys[1:3], keyMap, keyvals)

		// Should return all key/value pairs in lexicographic order.
		// Note that []byte("") is the lowest key possible and is
		// a special case in engine.scan, that's why we test it here.
		startKeys := []Key{Key("cat"), Key("")}
		for _, startKey := range startKeys {
			keyvals, err := engine.scan(startKey, KeyMax, 0)
			if err != nil {
				t.Fatalf("could not run scan: %v", err)
			}
			ensureRangeEqual(t, sortedKeys, keyMap, keyvals)
		}
	}, t)
}

func TestEngineIncrement(t *testing.T) {
	runWithAllEngines(func(engine Engine, t *testing.T) {
		// Start with increment of an empty key.
		val, err := increment(engine, Key("a"), 1, 0)
		if err != nil {
			t.Fatal(err)
		}
		if val != 1 {
			t.Errorf("expected increment to be %d; got %d", 1, val)
		}
		// Increment same key by 1.
		if val, err = increment(engine, Key("a"), 1, 0); err != nil {
			t.Fatal(err)
		}
		if val != 2 {
			t.Errorf("expected increment to be %d; got %d", 2, val)
		}
		// Increment same key by 2.
		if val, err = increment(engine, Key("a"), 2, 0); err != nil {
			t.Fatal(err)
		}
		if val != 4 {
			t.Errorf("expected increment to be %d; got %d", 4, val)
		}
		// Decrement same key by -1.
		if val, err = increment(engine, Key("a"), -1, 0); err != nil {
			t.Fatal(err)
		}
		if val != 3 {
			t.Errorf("expected increment to be %d; got %d", 3, val)
		}
		// Increment same key by max int64 value to cause overflow; should return error.
		if val, err = increment(engine, Key("a"), math.MaxInt64, 0); err == nil {
			t.Error("expected an overflow error")
		}
		if val, err = increment(engine, Key("a"), 0, 0); err != nil {
			t.Fatal(err)
		}
		if val != 3 {
			t.Errorf("expected increment to be %d; got %d", 3, val)
		}
	}, t)
}

func verifyScan(start, end Key, max int64, expKeys []Key, engine Engine, t *testing.T) {
	kvs, err := engine.scan(start, end, max)
	if err != nil {
		t.Errorf("scan %q-%q: expected no error, but got %s", string(start), string(end), err)
	}
	if len(kvs) != len(expKeys) {
		t.Errorf("scan %q-%q: expected scanned keys mismatch %d != %d: %v",
			start, end, len(kvs), len(expKeys), kvs)
	}
	for i, kv := range kvs {
		if !bytes.Equal(kv.Key, expKeys[i]) {
			t.Errorf("scan %q-%q: expected keys equal %q != %q", string(start), string(end),
				string(kv.Key), string(expKeys[i]))
		}
	}
}

func TestEngineScan2(t *testing.T) {
	// TODO(Tobias): Merge this with TestEngineScan1 and remove
	// either verifyScan or the other helper function.
	runWithAllEngines(func(engine Engine, t *testing.T) {
		keys := []Key{
			Key("a"),
			Key("aa"),
			Key("aaa"),
			Key("ab"),
			Key("abc"),
			KeyMax,
		}

		// Add keys to store in random order (make sure they sort!).
		order := rand.Perm(len(keys))
		for idx := range order {
			if err := engine.put(keys[idx], Value{Bytes: []byte("value")}); err != nil {
				t.Errorf("put: expected no error, but got %s", err)
			}
		}

		// Scan all keys (non-inclusive of final key).
		verifyScan(KeyMin, KeyMax, 10, keys[0:5], engine, t)
		verifyScan(Key("a"), KeyMax, 10, keys[0:5], engine, t)

		// Scan sub range.
		verifyScan(Key("aab"), Key("abcc"), 10, keys[3:5], engine, t)
		verifyScan(Key("aa0"), Key("abcc"), 10, keys[2:5], engine, t)

		// Scan with max values.
		verifyScan(KeyMin, KeyMax, 3, keys[0:3], engine, t)
		verifyScan(Key("a0"), KeyMax, 3, keys[1:4], engine, t)

		// Scan with max value 0 gets all values.
		verifyScan(KeyMin, KeyMax, 0, keys[0:5], engine, t)
	}, t)
}
