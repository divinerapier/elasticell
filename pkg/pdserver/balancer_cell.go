// Copyright 2016 DeepFabric, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package pdserver

import (
	"github.com/deepfabric/elasticell/pkg/pb/metapb"
)

type balanceCellScheduler struct {
	cfg      *Cfg
	cache    *idCache
	limit    uint64
	selector Selector
}

func (s *balanceCellScheduler) GetName() string {
	return "balance-cell-scheduler"
}

func (s *balanceCellScheduler) GetResourceKind() ResourceKind {
	return cellKind
}

func (s *balanceCellScheduler) GetResourceLimit() uint64 {
	return minUint64(s.limit, s.cfg.Schedule.CellScheduleLimit)
}

func (s *balanceCellScheduler) Prepare(cache *cache) error { return nil }

func (s *balanceCellScheduler) Cleanup(cache *cache) {}

func (s *balanceCellScheduler) Schedule(cache *cache) Operator {
	// Select a peer from the store with most cells.
	cell, oldPeer := scheduleRemovePeer(cache, s.selector)
	if cell == nil {
		return nil
	}

	// We don't schedule cell with abnormal number of replicas.
	if len(cell.getPeers()) != int(s.cfg.getMaxReplicas()) {
		return nil
	}

	op := s.transferPeer(cache, cell, oldPeer)
	if op == nil {
		// We can't transfer peer from this store now, so we add it to the cache
		// and skip it for a while.
		s.cache.set(oldPeer.StoreID)
	}
	return op
}

func (s *balanceCellScheduler) transferPeer(cache *cache, cell *cellRuntimeInfo, oldPeer *metapb.Peer) Operator {
	// scoreGuard guarantees that the distinct score will not decrease.
	stores := cache.getCellStores(cell)
	source := cache.getStore(oldPeer.StoreID)
	scoreGuard := newDistinctScoreFilter(s.cfg, stores, source)

	checker := newReplicaChecker(s.cfg, cache)
	newPeer, _ := checker.selectBestPeer(cell, scoreGuard)
	if newPeer == nil {
		return nil
	}

	target := cache.getStore(newPeer.StoreID)
	if !shouldBalance(source, target, s.GetResourceKind()) {
		return nil
	}
	s.limit = adjustBalanceLimit(cache, s.GetResourceKind())

	return newTransferPeerAggregationOp(cell, oldPeer, newPeer)
}
