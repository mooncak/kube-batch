/*
Copyright 2017 The Kubernetes Authors.

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

package cache

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
)

func nodeInfoEqual(l, r *NodeInfo) bool {
	if !reflect.DeepEqual(l, r) {
		return false
	}

	return true
}

func TestNodeInfo_AddPod(t *testing.T) {
	// case1
	node := buildNode("n1", buildResourceList("8000m", "10G"))
	pod1 := buildPod("c1", "p1", "n1", v1.PodRunning, buildResourceList("1000m", "1G"))
	pod2 := buildPod("c1", "p2", "n1", v1.PodRunning, buildResourceList("2000m", "2G"))

	tests := []struct {
		name     string
		node     *v1.Node
		pods     []*v1.Pod
		expected *NodeInfo
	}{
		{
			name: "add 2 running non-owner pod",
			node: node,
			pods: []*v1.Pod{pod1, pod2},
			expected: &NodeInfo{
				Name:        "n1",
				Node:        node,
				Idle:        buildResource("5000m", "7G"),
				Used:        buildResource("3000m", "3G"),
				Allocatable: buildResource("8000m", "10G"),
				Capability:  buildResource("8000m", "10G"),
				Pods: map[string]*PodInfo{
					"c1/p1": NewPodInfo(pod1),
					"c1/p2": NewPodInfo(pod2),
				},
			},
		},
	}

	for i, test := range tests {
		ni := NewNodeInfo(test.node)

		for _, pod := range test.pods {
			pi := NewPodInfo(pod)
			ni.AddPod(pi)
		}

		if !nodeInfoEqual(ni, test.expected) {
			t.Errorf("node info %d: \n expected %v, \n got %v \n",
				i, test.expected, ni)
		}
	}
}

func TestNodeInfo_RemovePod(t *testing.T) {
	// case1
	node := buildNode("n1", buildResourceList("8000m", "10G"))
	pod1 := buildPod("c1", "p1", "n1", v1.PodRunning, buildResourceList("1000m", "1G"))
	pod2 := buildPod("c1", "p2", "n1", v1.PodRunning, buildResourceList("2000m", "2G"))
	pod3 := buildPod("c1", "p3", "n1", v1.PodRunning, buildResourceList("3000m", "3G"))

	tests := []struct {
		name     string
		node     *v1.Node
		pods     []*v1.Pod
		rmPods   []*v1.Pod
		expected *NodeInfo
	}{
		{
			name:   "add 3 running non-owner pod, remove 1 running non-owner pod",
			node:   node,
			pods:   []*v1.Pod{pod1, pod2, pod3},
			rmPods: []*v1.Pod{pod2},
			expected: &NodeInfo{
				Name:        "n1",
				Node:        node,
				Idle:        buildResource("4000m", "6G"),
				Used:        buildResource("4000m", "4G"),
				Allocatable: buildResource("8000m", "10G"),
				Capability:  buildResource("8000m", "10G"),
				Pods: map[string]*PodInfo{
					"c1/p1": NewPodInfo(pod1),
					"c1/p3": NewPodInfo(pod3),
				},
			},
		},
	}

	for i, test := range tests {
		ni := NewNodeInfo(test.node)

		for _, pod := range test.pods {
			pi := NewPodInfo(pod)
			ni.AddPod(pi)
		}

		for _, pod := range test.rmPods {
			pi := NewPodInfo(pod)
			ni.RemovePod(pi)
		}

		if !nodeInfoEqual(ni, test.expected) {
			t.Errorf("node info %d: \n expected %v, \n got %v \n",
				i, test.expected, ni)
		}
	}
}