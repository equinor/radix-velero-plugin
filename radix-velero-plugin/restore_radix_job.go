/*
Copyright 2018 the Heptio Ark contributors.

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
	"encoding/json"
	"fmt"

	"github.com/equinor/radix-operator/pkg/apis/kube"
	"github.com/equinor/radix-velero-plugin/models"

	radixv1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

// RestoreRadixJobPlugin is a restore item action plugin for Velero
type RestoreRadixJobPlugin struct {
	Log      logrus.FieldLogger
	kubeUtil *models.Kube
}

// AppliesTo returns information about which resources this action should be invoked for.
// A RestoreItemAction's Execute function will only be invoked on items that match the returned
// selector. A zero-valued ResourceSelector matches all resources.g
func (p *RestoreRadixJobPlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"radixjobs.radix.equinor.com"},
	}, nil
}

// Execute allows the RestorePlugin to perform arbitrary logic with the item being restored,
// in this case, setting a custom annotation on the item being restored.
func (p *RestoreRadixJobPlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.Log.Info("Radix Job RestorePlugin!")

	metadata, err := meta.Accessor(input.Item)
	if err != nil {
		return &velero.RestoreItemActionExecuteOutput{}, err
	}

	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	var rj radixv1.RadixJob
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(input.ItemFromBackup.UnstructuredContent(), &rj); err != nil {
		return nil, fmt.Errorf("unable to convert unstructured item to Radix Deployment: %w", err)
	}

	restoredStatus, err := json.Marshal(rj.Status)

	if err != nil {
		return &velero.RestoreItemActionExecuteOutput{}, err
	}

	radixAppName := rj.Labels[kube.RadixAppLabel]
	rrExists, err := p.kubeUtil.ExistsRadixRegistration(radixAppName)
	if err != nil {
		return &velero.RestoreItemActionExecuteOutput{}, err
	}
	if rrExists {
		annotations["equinor.com/velero-restored-status"] = string(restoredStatus)
		metadata.SetAnnotations(annotations)
		return velero.NewRestoreItemActionExecuteOutput(input.Item), nil
	}

	p.Log.Infof("RadixRegistration %s does not exists - skip restoring RadixJob", radixAppName)
	return &velero.RestoreItemActionExecuteOutput{
		SkipRestore: true,
	}, nil
}
