/*
Copyright 2017 the Heptio Ark contributors.

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
	"github.com/equinor/radix-velero-plugin/models"
	"github.com/sirupsen/logrus"
	veleroplugin "github.com/vmware-tanzu/velero/pkg/plugin/framework"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	kubeUtil, err := models.GetKubeUtil()
	if err != nil {
		logrus.Fatalf("cannot get Kubernetes or Radix client: %v", err)
		return
	}
	veleroplugin.NewServer().
		RegisterRestoreItemAction("equinor.com/restore-application-plugin", func(logger logrus.FieldLogger) (interface{}, error) {
			return &RestoreRadixApplicationPlugin{
				Log:      logger,
				kubeUtil: kubeUtil,
			}, nil
		}).
		RegisterRestoreItemAction("equinor.com/restore-deployment-plugin", func(logger logrus.FieldLogger) (interface{}, error) {
			return &RestoreRadixDeploymentPlugin{
				Log:      logger,
				kubeUtil: kubeUtil,
			}, nil
		}).
		RegisterRestoreItemAction("equinor.com/restore-job-plugin", func(logger logrus.FieldLogger) (interface{}, error) {
			return &RestoreRadixJobPlugin{
				Log:      logger,
				kubeUtil: kubeUtil,
			}, nil
		}).
		RegisterRestoreItemAction("equinor.com/restore-alert-plugin", func(logger logrus.FieldLogger) (interface{}, error) {
			return &RestoreAlertPlugin{
				Log:      logger,
				kubeUtil: kubeUtil,
			}, nil
		}).
		RegisterRestoreItemAction("equinor.com/restore-environment-plugin", func(logger logrus.FieldLogger) (interface{}, error) {
			return &RestoreRadixEnvironmentPlugin{
				Log:      logger,
				kubeUtil: kubeUtil,
			}, nil
		}).
		RegisterRestoreItemAction("equinor.com/restore-secret-plugin", func(logger logrus.FieldLogger) (interface{}, error) {
			return &RestoreRadixAppSecretPlugin{
				Log:      logger,
				kubeUtil: kubeUtil,
			}, nil
		}).
		RegisterRestoreItemAction("equinor.com/restore-configmap-plugin", func(logger logrus.FieldLogger) (interface{}, error) {
			return &RestoreRadixAppConfigMapPlugin{
				Log:      logger,
				kubeUtil: kubeUtil,
			}, nil
		}).
		Serve()
	logrus.Infoln("Initialized 'radix-velero-plugin'")
}
