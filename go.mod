module github.com/onosproject/onos-ric

go 1.13

require (
	github.com/atomix/api v0.0.0-20200219005318-0350f11bfcde
	github.com/atomix/go-client v0.0.0-20200218200323-6fd69e684d05
	github.com/atomix/go-framework v0.0.0-20200211010924-f3f12b63db0a
	github.com/atomix/go-local v0.0.0-20200211010611-c99e53e4c653
	github.com/docker/docker v1.13.1
	github.com/gogo/protobuf v1.3.1
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onosproject/onos-test v0.0.0-20200222204048-99418ae1ce43
	github.com/onosproject/onos-topo v0.0.0-20200218171206-55029b503689
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.27.1
	gotest.tools v2.2.0+incompatible
	k8s.io/klog v1.0.0
)

replace github.com/onosproject/onos-ric => ../onos-ric

replace github.com/onosproject/onos-test => ../onos-test
