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
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.
//
// Author: Ben Darnell

package multiraft

import (
	"reflect"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/util"
	"github.com/cockroachdb/cockroach/util/leaktest"
	"github.com/cockroachdb/cockroach/util/log"
	"github.com/coreos/etcd/raft"
	"github.com/coreos/etcd/raft/raftpb"
)

var testRand = util.NewPseudoRand()

func makeCommandID() string {
	return string(util.RandBytes(testRand, commandIDLen))
}

type testCluster struct {
	t         *testing.T
	nodes     []*state
	tickers   []*manualTicker
	events    []*eventDemux
	storages  []*BlockableStorage
	transport Transport
}

func newTestCluster(transport Transport, size int, t *testing.T) *testCluster {
	if transport == nil {
		transport = NewLocalRPCTransport()
	}
	cluster := &testCluster{t: t}
	cluster.transport = transport
	for i := 0; i < size; i++ {
		ticker := newManualTicker()
		storage := &BlockableStorage{storage: NewMemoryStorage()}
		config := &Config{
			Transport:              transport,
			Storage:                storage,
			Ticker:                 ticker,
			ElectionTimeoutTicks:   1,
			HeartbeatIntervalTicks: 1,
			TickInterval:           time.Hour, // not in use
			Strict:                 true,
		}
		mr, err := NewMultiRaft(NodeID(i+1), config)
		if err != nil {
			t.Fatal(err)
		}
		state := newState(mr)
		demux := newEventDemux(state.Events)
		demux.start()
		cluster.nodes = append(cluster.nodes, state)
		cluster.tickers = append(cluster.tickers, ticker)
		cluster.events = append(cluster.events, demux)
		cluster.storages = append(cluster.storages, storage)
	}
	cluster.start()
	return cluster
}

func (c *testCluster) start() {
	// Let all the states listen before starting any.
	for _, node := range c.nodes {
		go node.start()
	}
}

func (c *testCluster) stop() {
	for _, node := range c.nodes {
		node.Stop()
	}
	for _, demux := range c.events {
		demux.stop()
	}
	c.transport.Close()
}

// createGroup replicates a group consisting of numReplicas members,
// the first being the node at index firstNode.
func (c *testCluster) createGroup(groupID uint64, firstNode, numReplicas int) {
	var replicaIDs []uint64
	for i := 0; i < numReplicas; i++ {
		replicaIDs = append(replicaIDs, uint64(c.nodes[firstNode+i].nodeID))
	}
	for i := 0; i < numReplicas; i++ {
		gs := c.storages[firstNode+i].GroupStorage(groupID)
		memStorage := gs.(*blockableGroupStorage).s.(*raft.MemoryStorage)
		if err := memStorage.SetHardState(raftpb.HardState{
			Commit: 10,
			Term:   5,
		}); err != nil {
			c.t.Fatal(err)
		}
		if err := memStorage.ApplySnapshot(raftpb.Snapshot{
			Metadata: raftpb.SnapshotMetadata{
				ConfState: raftpb.ConfState{
					Nodes: replicaIDs,
				},
				Index: 10,
				Term:  5,
			},
		}); err != nil {
			c.t.Fatal(err)
		}

		node := c.nodes[firstNode+i]
		err := node.CreateGroup(groupID)
		if err != nil {
			c.t.Fatal(err)
		}
	}
}

func (c *testCluster) triggerElection(i int) {
	// Elections are currently triggered after ElectionTimeoutTicks+1 ticks.
	c.tickers[i].Tick()
	c.tickers[i].Tick()
}

// Trigger an election on node i and wait for it to complete.
// TODO(bdarnell): once we have better leader discovery and forwarding/queuing, remove this.
func (c *testCluster) waitForElection(i int) *EventLeaderElection {
	for {
		e := <-c.events[i].LeaderElection
		if e == nil {
			panic("got nil LeaderElection event, channel likely closed")
		}
		// Ignore events with NodeID 0; these mark elections that are in progress.
		if e.NodeID != 0 {
			return e
		}
	}
}

func TestInitialLeaderElection(t *testing.T) {
	defer leaktest.AfterTest(t)
	// Run the test three times, each time triggering a different node's election clock.
	// The node that requests an election first should win.
	for leaderIndex := 0; leaderIndex < 3; leaderIndex++ {
		log.Infof("testing leader election for node %v", leaderIndex)
		cluster := newTestCluster(nil, 3, t)
		groupID := uint64(1)
		cluster.createGroup(groupID, 0, 3)

		cluster.triggerElection(leaderIndex)
		event := cluster.waitForElection(leaderIndex)
		if event.GroupID != groupID {
			t.Fatalf("election event had incorrect group id %v", event.GroupID)
		}
		if event.NodeID != cluster.nodes[leaderIndex].nodeID {
			t.Fatalf("expected %v to win election, but was %v", cluster.nodes[leaderIndex].nodeID,
				event.NodeID)
		}
		cluster.stop()
	}
}

func TestLeaderElectionEvent(t *testing.T) {
	defer leaktest.AfterTest(t)
	// Leader election events are fired when the leader commits an entry, not when it
	// issues a call for votes.
	cluster := newTestCluster(nil, 3, t)
	defer cluster.stop()
	groupID := uint64(1)
	cluster.createGroup(groupID, 0, 3)

	// Send a Ready with a new leader but no new commits.
	// This happens while an election is in progress.
	cluster.nodes[1].handleRaftReady(map[uint64]raft.Ready{
		groupID: {
			SoftState: &raft.SoftState{
				Lead: 3,
			},
		},
	})

	// No events are sent.
	select {
	case e := <-cluster.events[1].LeaderElection:
		t.Fatalf("got unexpected event %v", e)
	case <-time.After(time.Millisecond):
	}

	// Now there are new committed entries. A new leader always commits an entry
	// to conclude the election.
	entry := raftpb.Entry{
		Index: 42,
		Term:  42,
	}
	cluster.nodes[1].handleRaftReady(map[uint64]raft.Ready{
		groupID: {
			Entries:          []raftpb.Entry{entry},
			CommittedEntries: []raftpb.Entry{entry},
		},
	})

	// Now we get an event.
	select {
	case e := <-cluster.events[1].LeaderElection:
		if !reflect.DeepEqual(e, &EventLeaderElection{
			GroupID: groupID,
			NodeID:  3,
			Term:    42,
		}) {
			t.Errorf("election event did not match expectations: %+v", e)
		}
	case <-time.After(time.Millisecond):
		t.Fatal("didn't get expected event")
	}
}

func TestCommand(t *testing.T) {
	defer leaktest.AfterTest(t)
	cluster := newTestCluster(nil, 3, t)
	defer cluster.stop()
	groupID := uint64(1)
	cluster.createGroup(groupID, 0, 3)
	cluster.triggerElection(0)
	cluster.waitForElection(0)

	// Submit a command to the leader
	cluster.nodes[0].SubmitCommand(groupID, makeCommandID(), []byte("command"))

	// The command will be committed on each node.
	for i, events := range cluster.events {
		log.Infof("waiting for event to be committed on node %v", i)
		commit := <-events.CommandCommitted
		if string(commit.Command) != "command" {
			t.Errorf("unexpected value in committed command: %v", commit.Command)
		}
	}
}

func TestSlowStorage(t *testing.T) {
	defer leaktest.AfterTest(t)
	cluster := newTestCluster(nil, 3, t)
	defer cluster.stop()
	groupID := uint64(1)
	cluster.createGroup(groupID, 0, 3)

	cluster.triggerElection(0)
	cluster.waitForElection(0)

	// Block the storage on the last node.
	// TODO(bdarnell): there appear to still be issues if the storage is blocked during
	// the election.
	cluster.storages[2].Block()

	// Submit a command to the leader
	cluster.nodes[0].SubmitCommand(groupID, makeCommandID(), []byte("command"))

	// Even with the third node blocked, the other nodes can make progress.
	for i := 0; i < 2; i++ {
		events := cluster.events[i]
		log.Infof("waiting for event to be commited on node %v", i)
		commit := <-events.CommandCommitted
		if string(commit.Command) != "command" {
			t.Errorf("unexpected value in committed command: %v", commit.Command)
		}
	}

	// Ensure that node 2 is in fact blocked.
	time.Sleep(time.Millisecond)
	select {
	case commit := <-cluster.events[2].CommandCommitted:
		t.Errorf("didn't expect commits on node 2 but got %v", commit)
	default:
	}

	// After unblocking the third node, it will catch up.
	cluster.storages[2].Unblock()
	cluster.tickers[0].Tick()
	log.Infof("waiting for event to be commited on node 2")
	commit := <-cluster.events[2].CommandCommitted
	if string(commit.Command) != "command" {
		t.Errorf("unexpected value in committed command: %v", commit.Command)
	}
}

func TestMembershipChange(t *testing.T) {
	defer leaktest.AfterTest(t)
	cluster := newTestCluster(nil, 4, t)
	defer cluster.stop()

	// Create a group with a single member, cluster.nodes[0].
	groupID := uint64(1)
	cluster.createGroup(groupID, 0, 1)
	cluster.triggerElection(0)
	cluster.waitForElection(0)

	// Consume and apply the membership change events.
	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				e, ok := <-cluster.events[i].MembershipChangeCommitted
				if !ok {
					return
				}
				e.Callback(nil)
			}
		}(i)
	}

	// Add each of the other three nodes to the cluster.
	for i := 1; i < 4; i++ {
		ch := cluster.nodes[0].ChangeGroupMembership(groupID, makeCommandID(),
			raftpb.ConfChangeAddNode,
			cluster.nodes[i].nodeID, nil)
		<-ch
	}

	// TODO(bdarnell): verify that the channel events are sent out correctly.
	/*
		for i := 0; i < 10; i++ {
			log.Infof("tick %d", i)
			cluster.tickers[0].Tick()
			time.Sleep(5 * time.Millisecond)
		}

		// Each node is notified of each other node's joining.
		for i := 0; i < 4; i++ {
			for j := 1; j < 4; j++ {
				select {
				case e := <-cluster.events[i].MembershipChangeCommitted:
					if e.NodeID != cluster.nodes[j].nodeID {
						t.Errorf("node %d expected event for %d, got %d", i, j, e.NodeID)
					}
				default:
					t.Errorf("node %d did not get expected event for %d", i, j)
				}
			}
		}*/
}
