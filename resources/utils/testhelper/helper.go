// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package testhelper

func CollectResources[T any](ch chan T) []T {
	var numResources = len(ch)
	var results []T
	for i := 0; i < numResources; i++ {
		results = append(results, <-ch)
	}

	return results
}

func CreateMockClients[T any](client T, regions []string) map[string]T {
	var m = make(map[string]T, 0)
	for _, clientRegion := range regions {
		m[clientRegion] = client
	}

	return m
}
