module github.com/onosproject/onos-ric

go 1.13

require (
	github.com/atomix/go-client v0.0.0-20200307025134-f638fa3fb644
	github.com/gogo/protobuf v1.3.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onosproject/helmit v0.0.0-20200327211207-6ee099c52d08
	github.com/onosproject/onos-lib-go v0.0.0-20200311221003-fad88142208e
	github.com/onosproject/onos-topo v0.0.0-20200218171206-55029b503689
	github.com/prometheus/client_golang v1.4.1
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.27.1
	gotest.tools v2.2.0+incompatible
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
