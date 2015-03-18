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
// Author: Spencer Kimball (spencer.kimball@gmail.com)

package storage

import (
	"github.com/cockroachdb/cockroach/proto"
	"github.com/cockroachdb/cockroach/storage/engine"
	"github.com/cockroachdb/cockroach/util/encoding"
)

// keyRange is a helper struct for the rangeDataIterator.
type keyRange struct {
	start, end proto.EncodedKey
}

// rangeDataIterator provides a complete iteration over all key / value
// rows in a range, including all system-local metadata and user data.
// The ranges keyRange slice specifies the key ranges which comprise
// all of the range's data.
//
// A rangeDataIterator provides the same API as an Engine iterator
// with the exception of the Seek() method.
type rangeDataIterator struct {
	curIndex int
	ranges   []keyRange
	iter     engine.Iterator
}

func newRangeDataIterator(r *Range, e engine.Engine) *rangeDataIterator {
	r.RLock()
	startKey := r.Desc().StartKey
	if startKey.Equal(engine.KeyMin) {
		startKey = engine.KeyLocalMax
	}
	endKey := r.Desc().EndKey
	r.RUnlock()
	ri := &rangeDataIterator{
		ranges: []keyRange{
			{
				start: engine.MVCCEncodeKey(engine.MakeKey(engine.KeyLocalRangeIDPrefix, encoding.EncodeUvarint(nil, uint64(r.Desc().RaftID)))),
				end:   engine.MVCCEncodeKey(engine.MakeKey(engine.KeyLocalRangeIDPrefix, encoding.EncodeUvarint(nil, uint64(r.Desc().RaftID+1)))),
			},
			{
				start: engine.MVCCEncodeKey(engine.MakeKey(engine.KeyLocalRangeKeyPrefix, encoding.EncodeBytes(nil, startKey))),
				end:   engine.MVCCEncodeKey(engine.MakeKey(engine.KeyLocalRangeKeyPrefix, encoding.EncodeBytes(nil, endKey))),
			},
			{
				start: engine.MVCCEncodeKey(startKey),
				end:   engine.MVCCEncodeKey(endKey),
			},
		},
		iter: e.NewIterator(),
	}
	ri.iter.Seek(ri.ranges[ri.curIndex].start)
	ri.advance()
	return ri
}

// Close closes the underlying iterator.
func (ri *rangeDataIterator) Close() {
	ri.curIndex = len(ri.ranges)
	ri.iter.Close()
}

func (ri *rangeDataIterator) Valid() bool {
	return ri.iter.Valid()
}

// Next returns the next raw key value in the iteration, or nil if
// iteration is done.
func (ri *rangeDataIterator) Next() {
	ri.iter.Next()
	ri.advance()
}

// Key returns the current Key for the iteration if valid.
func (ri *rangeDataIterator) Key() proto.EncodedKey {
	return ri.iter.Key()
}

// Value returns the current Value for the iteration if valid.
func (ri *rangeDataIterator) Value() []byte {
	return ri.iter.Value()
}

// Error returns the Error for the iteration if applicable.
func (ri *rangeDataIterator) Error() error {
	return ri.iter.Error()
}

// advance moves the iterator forward through the ranges until a valid
// key is found or the iteration is done and the iterator becomes
// invalid.
func (ri *rangeDataIterator) advance() {
	for {
		if !ri.iter.Valid() || ri.iter.Key().Less(ri.ranges[ri.curIndex].end) {
			return
		}
		ri.curIndex++
		if ri.curIndex < len(ri.ranges) {
			ri.iter.Seek(ri.ranges[ri.curIndex].start)
		} else {
			// Otherwise, seek to end to make iterator invalid.
			ri.iter.Seek(engine.MVCCKeyMax)
			return
		}
	}
}
