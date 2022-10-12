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

package add_enviroment_metadata

import (
	"fmt"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/processors"
	jsprocessor "github.com/elastic/beats/v7/libbeat/processors/script/javascript/module/processor"
	"github.com/elastic/elastic-agent-autodiscover/kubernetes"
	"github.com/elastic/elastic-agent-autodiscover/kubernetes/metadata"
	agentconfig "github.com/elastic/elastic-agent-libs/config"
	"github.com/elastic/elastic-agent-libs/logp"
)

const (
	processorName  = "add_environment_metadata"
	clusterNameKey = "orchestrator.cluster.name"
)

func init() {
	processors.RegisterPlugin(processorName, New)
	jsprocessor.RegisterPlugin("AddEnvironmentMetadata", New)
}

type addEnvironmentMetadata struct {
	ClusterName string
	logger      *logp.Logger
	config      Config
}

// New constructs a new add environment metadata processor.
func New(cfg *agentconfig.C) (processors.Processor, error) {
	config := defaultConfig()

	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("fail to unpack the %s configuration %v, skipping with error: %v", processorName, cfg.FlattenedKeys(), err)
	}

	client, err := kubernetes.GetKubernetesClient("", kubernetes.KubeClientOptions{})
	if err != nil {
		return nil, err
	}

	logger := logp.NewLogger(processorName)
	clusterIdentifier, err := metadata.GetKubernetesClusterIdentifier(cfg, client)
	if err == nil {
		logger.Errorf("fail to resolve the name of the cluster, error %v", err)
	}

	p := &addEnvironmentMetadata{
		ClusterName: clusterIdentifier.Name,
		logger:      logger,
		config:      config,
	}

	return p, nil
}

// Run enriches the given event with the environment metadata
func (p *addEnvironmentMetadata) Run(event *beat.Event) (*beat.Event, error) {
	clusterName := p.ClusterName

	if clusterName != "" {
		_, err := event.PutValue(clusterNameKey, clusterName)
		if err != nil {
			return nil, fmt.Errorf("failed to add cluster name to object: %v", err)
		}
	}

	return event, nil
}

func (p *addEnvironmentMetadata) String() string {
	return fmt.Sprintf("%v=[%s=[%v]]", processorName, clusterNameKey, p.ClusterName)
}
