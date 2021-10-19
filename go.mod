module github.com/equinor/radix-velero-plugin

go 1.16

require (
	github.com/equinor/radix-operator v1.16.3
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/vmware-tanzu/velero v1.6.2
	k8s.io/api v0.19.12
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v12.0.0+incompatible
)

replace k8s.io/client-go => k8s.io/client-go v0.19.9
