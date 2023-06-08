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

package systemtest

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	apmqueue "github.com/elastic/apm-queue"
)

// ProvisionInfraFunc is a function returned by Init* functions for
// provisioning infrastructure.
type ProvisionInfraFunc func(context.Context) error

// DestroyInfraFunc is a function returned by Init* functions for
// destroying infrastructure.
type DestroyInfraFunc func(context.Context) error

var persistentSuffix string

func init() {
	rand.Seed(time.Now().Unix())
	persistentSuffix = RandomSuffix()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

// RandomSuffix generates a lowercase alphabetic 8 character random string
func RandomSuffix() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// SuffixTopics suffixes the received topics with a random suffix.
func SuffixTopics(topics ...apmqueue.Topic) []apmqueue.Topic {
	suffix := RandomSuffix()
	suffixed := make([]apmqueue.Topic, len(topics))
	for i := range suffixed {
		topic := fmt.Sprintf("%s.%s", strings.ToLower(string(topics[i])), suffix)
		topic = strings.ReplaceAll(topic, "_", "-")
		topic = strings.ReplaceAll(topic, "/", "-")
		suffixed[i] = apmqueue.Topic(topic)
	}
	return suffixed
}

// SingleSubscribers returns a single apmqueue.Subscription for each
// topic, with the subscriber having the same name as the topic.
func SingleSubscribers(topics ...apmqueue.Topic) []apmqueue.Subscription {
	out := make([]apmqueue.Subscription, len(topics))
	for i, topic := range topics {
		out[i] = apmqueue.Subscription{
			Name:  string(topic),
			Topic: topic,
		}
	}
	return out
}
