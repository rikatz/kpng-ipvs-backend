/*
Copyright 2021 The Kubernetes Authors.

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
	"github.com/rikatz/kpng-ipvs-backend/pkg/ipvs"
	"k8s.io/klog"
	"sigs.k8s.io/kpng/pkg/client"
)

func main() {
	err := ipvs.PreRun()
	if err != nil {
		klog.Fatalf("Error starting kpng ipvs: %s", err)
	}
	client.RunCh(ipvs.Callback, ipvs.BindFlags)
}
