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
	kube "github.com/equinor/radix-operator/pkg/apis/kube"
	"github.com/equinor/radix-velero-plugin/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// RestoreRadixAppConfigMapPlugin is a restore item action plugin for Velero
type RestoreRadixAppConfigMapPlugin struct {
	Log      logrus.FieldLogger
	kubeUtil *models.Kube
}

// AppliesTo returns information about which resources this action should be invoked for.
// A RestoreItemAction's Execute function will only be invoked on items that match the returned
// selector. A zero-valued ResourceSelector matches all resources.g
func (p *RestoreRadixAppConfigMapPlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"configmaps"},
	}, nil
}

// Execute allows the RestorePlugin to perform arbitrary logic with the item being restored,
// in this case, setting a custom annotation on the item being restored.
func (p *RestoreRadixAppConfigMapPlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.Log.Info("Radix RadixAppConfigMap RestorePlugin!")

	var configMap corev1.ConfigMap
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(input.ItemFromBackup.UnstructuredContent(), &configMap); err != nil {
		return nil, errors.Wrap(err, "unable to convert unstructured item to ConfigMap for RadixApp")
	}

	radixAppName := configMap.Labels[kube.RadixAppLabel]
	if radixAppName == "" {
		return velero.NewRestoreItemActionExecuteOutput(input.Item), nil
	}

	rrExists, err := p.kubeUtil.ExistsRadixRegistration(radixAppName)
	if err != nil {
		return &velero.RestoreItemActionExecuteOutput{}, err
	}
	if rrExists {
		return velero.NewRestoreItemActionExecuteOutput(input.Item), nil
	}

	p.Log.Infof("RadixRegistration %s does not exists - skip restoring ConfigMap for RadixApp", radixAppName)
	return &velero.RestoreItemActionExecuteOutput{
		SkipRestore: true,
	}, nil
}
