// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var processorLock = &sync.Mutex{}

func reconcileUnscheduledPods(interval int, done chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Second):
			err := schedulePods()
			if err != nil {
				log.Println(err)
			}
		case <-done:
			wg.Done()
			log.Println("Stopped reconciliation loop.")
			return
		}
	}
}

func monitorUnscheduledPods(done chan struct{}, wg *sync.WaitGroup) {
	pods, errc := watchUnscheduledPods()
	//fmt.Println("monitorUnscheduledPods", "<-pods")
	//processorLock.Lock()
	//time.Sleep(2 * time.Second)
	//schedulePodUsingSolver()
	//processorLock.Unlock()

	for {
		select {
		case err := <-errc:
			log.Println(err)
		case <-pods:
			processorLock.Lock()
			time.Sleep(2 * time.Second)
			schedulePodUsingSolver()
			processorLock.Unlock()
		case <-done:
			wg.Done()
			log.Println("Stopped scheduler.")
			return
		}
	}
}

func schedulePod(pod *Pod) error {
	nodes, err := fit(pod)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("Unable to schedule pod (%s) failed to fit in any node", pod.Metadata.Name)
	}
	node, err := bestPrice(nodes)
	if err != nil {
		return err
	}
	//printPod(*pod)
	err = bind(*pod, node)
	//printNode(node)
	if err != nil {
		return err
	}
	return nil
}

func schedulePods() error {
	fmt.Println("===func schedulePods()")
	processorLock.Lock()
	defer processorLock.Unlock()
	pods, err := getUnscheduledPods()
	if err != nil {
		return err
	}

	if len(pods) != 0 {
		schedulePodUsingSolver()
	}
	return nil
}
