module github.com/valaparthvi/highlander-tests

go 1.19

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/epinio/epinio v1.10.0
	github.com/pkg/errors v0.9.1
	k8s.io/utils v0.0.0-20230505201702-9f6742963106
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.12.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

require (
	github.com/onsi/ginkgo/v2 v2.13.0
	github.com/onsi/gomega v1.30.0
	k8s.io/apimachinery v0.27.6
)

replace k8s.io/api => k8s.io/api v0.27.5

replace k8s.io/apimachinery => k8s.io/apimachinery v0.27.5

replace github.com/rancher/rancher/pkg/client => github.com/rancher/rancher/pkg/client v0.0.0-20231113162426-5b42ca504753

replace github.com/rancher/rancher/pkg/apis => github.com/rancher/rancher/pkg/apis v0.0.0-20231113162426-5b42ca504753

replace k8s.io/client-go => github.com/rancher/client-go v1.27.4-rancher1
