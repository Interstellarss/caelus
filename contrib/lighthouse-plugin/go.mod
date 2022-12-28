module github.com/Interstellarss/lighthouse-plugin

go 1.14

require (
	github.com/containerd/containerd v1.6.12
	github.com/coreos/go-systemd/v22 v22.3.2
	github.com/docker/docker v20.10.1+incompatible
	github.com/evanphx/json-patch v4.12.0+incompatible
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/json-iterator/go v1.1.12
	github.com/mYmNeo/version v0.0.0-20200424030557-30e59e77cc3e
	github.com/opencontainers/runc v1.1.4
	github.com/opencontainers/selinux v1.10.2 // indirect
	//github.com/opencontainers/runc v1.1.2
	//github.com/opencontainers/selinux v1.10.2 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	google.golang.org/grpc v1.51.0 // indirect
	k8s.io/api v0.24.2
	k8s.io/apimachinery v0.24.2
	k8s.io/apiserver v0.24.2
	k8s.io/client-go v0.24.2
	k8s.io/component-base v0.24.2
	//k8s.io/api v0.22.5
	//k8s.io/apimachinery v0.22.5
	//k8s.io/apiserver v0.22.5
	//k8s.io/client-go v0.22.5
	//k8s.io/component-base v0.22.5
	k8s.io/klog v1.0.0
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
)

replace (
	github.com/Sirupsen/logrus v1.7.0 => github.com/sirupsen/logrus v1.7.0
	github.com/sirupsen/logrus v1.7.0 => github.com/Sirupsen/logrus v1.7.0
)
