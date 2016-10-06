/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"sort"

	"golang.org/x/net/context"
	rpb "k8s.io/helm/pkg/proto/hapi/release"
	tpb "k8s.io/helm/pkg/proto/hapi/services"
)

func (s *releaseServer) GetHistory(ctx context.Context, req *tpb.GetHistoryRequest) (*tpb.GetHistoryResponse, error) {
	if !checkClientVersion(ctx) {
		return nil, errIncompatibleVersion
	}

	h, err := s.env.Releases.History(req.Name)
	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(byRev(h)))

	var resp tpb.GetHistoryResponse
	for i := 0; i < min(len(h), int(req.Max)); i++ {
		resp.Releases = append(resp.Releases, h[i])
	}

	return &resp, nil
}

type byRev []*rpb.Release

func (s byRev) Len() int           { return len(s) }
func (s byRev) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byRev) Less(i, j int) bool { return s[i].Version < s[j].Version }

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
