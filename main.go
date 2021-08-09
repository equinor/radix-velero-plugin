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
	"github.com/equinor/radix-velero-plugin/velero-plugins/status"
	"github.com/sirupsen/logrus"
	veleroplugin "github.com/vmware-tanzu/velero/pkg/plugin/framework"
)

func main() {
	veleroplugin.NewServer().
		RegisterRestoreItemAction("equinor.com/restore-deployment-plugin", newDeploymentRestorePlugin).
		RegisterRestoreItemAction("equinor.com/restore-job-plugin", newJobRestorePlugin).
		Serve()
}

func newDeploymentRestorePlugin(logger logrus.FieldLogger) (interface{}, error) {
	return &status.RestoreDeploymentPlugin{Log: logger}, nil
}

func newJobRestorePlugin(logger logrus.FieldLogger) (interface{}, error) {
	return &status.RestoreJobPlugin{Log: logger}, nil
}
