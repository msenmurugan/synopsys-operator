/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownershia. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package plugins

import (
	"errors"
	"fmt"
	"time"

	"github.com/blackducksoftware/perceptor-protoform/pkg/hub/hubutils"
	"github.com/sirupsen/logrus"

	hubv1 "github.com/blackducksoftware/perceptor-protoform/pkg/api/hub/v1"
	hubclient "github.com/blackducksoftware/perceptor-protoform/pkg/hub/client/clientset/versioned"
	"github.com/blackducksoftware/perceptor-protoform/pkg/protoform"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// InitDatabaseUpdater will hold the configuration to initialize the Postgres database
type HubStatusUpdater struct {
	Config     *protoform.Config
	KubeClient *kubernetes.Clientset
	HubClient  *hubclient.Clientset
	Namespace  string
}

func (i *HubStatusUpdater) update() {
	hL, _ := i.HubClient.SynopsysV1().Hubs("").List(metav1.ListOptions{})
	for _, hub := range hL.Items {
		podList, _ := i.KubeClient.Core().Pods("").List(metav1.ListOptions{})
		hisstorg := map[string]string{}
		for _, pod := range podList.Items {
			if pod.Status.Phase != v1.PodRunning {
				hisstorg[pod.Name] = fmt.Sprintf("%v", pod.Status.Phase)
			}
		}

		// if any entreis non running, its status is busted ...
		if len(hisstorg) > 0 {
			logrus.Warnf("Warning: Hub %v is down  %v", hub.GetNamespace, hisstorg)
			hubutils.UpdateState(i.HubClient, "", "", errors.New(fmt.Sprintf("%v", hisstorg)), &hub)
		}
	}
}

// Run is a BLOCKING function which should be run by the framework .
func (i *HubStatusUpdater) Run(ch <-chan struct{}) {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return i.KubeClient.Core().Pods(i.Namespace).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return i.KubeClient.Core().Pods(i.Namespace).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&hubv1.Hub{},
		2*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				i.update()
			},
			UpdateFunc: func(obj interface{}, obj2 interface{}) {
				i.update()
			},
			DeleteFunc: func(obj interface{}) {
				i.update()
			},
		},
	)
	ctrl.Run(ch)
	// Wait until we're told to stop
	<-ch
}