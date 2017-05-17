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

// +build freebsd openbsd netbsd dragonfly linux

package storage

import (
	gonemo "github.com/deepfabric/go-nemo"
)

type nemoSetEngine struct {
	db *gonemo.NEMO
}

func newNemoSetEngine(db *gonemo.NEMO) SetEngine {
	return &nemoSetEngine{
		db: db,
	}
}

func (e *nemoSetEngine) SAdd(key []byte, members ...[]byte) (int64, error) {
	// TODO: nemo must support more members
	return e.db.SAdd(key, members[0])
}

func (e *nemoSetEngine) SRem(key []byte, members ...[]byte) (int64, error) {
	// TODO: nemo must support more members
	return e.db.SRem(key, members[0])
}

func (e *nemoSetEngine) SCard(key []byte) (int64, error) {
	// TODO: nemo must return a error
	return e.db.SCard(key), nil
}

func (e *nemoSetEngine) SMembers(key []byte) ([][]byte, error) {
	return e.db.SMembers(key)
}

func (e *nemoSetEngine) SIsMember(key []byte, member []byte) (int64, error) {
	yes := e.db.SIsMember(key, member)
	var value int64
	if yes {
		value = 1
	}
	// TODO: nemo must return a error
	return value, nil
}

func (e *nemoSetEngine) SPop(key []byte) ([]byte, error) {
	return e.db.SPop(key)
}
