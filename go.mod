module github.com/onosproject/onos-ric

go 1.13

require (
	github.com/DataDog/mmh3 v0.0.0-20200316233529-f5b682d8c981
	github.com/atomix/go-client v0.1.0
	github.com/atomix/go-framework v0.0.0-20200211010924-f3f12b63db0a // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/joncalhoun/pipe v0.0.0-20170510025636-72505674a733
	github.com/onosproject/helmit v0.5.0
	github.com/onosproject/onos-lib-go v0.5.0
	github.com/onosproject/onos-test v0.0.0-20200222204048-99418ae1ce43 // indirect
	github.com/onosproject/onos-topo v0.5.0
	github.com/prometheus/client_golang v1.4.1
	github.com/renstrom/dedent v1.0.0 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/stretchr/testify v1.5.1
	golang.org/x/tools v0.0.0-20200313205530-4303120df7d8 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools v2.2.0+incompatible
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
