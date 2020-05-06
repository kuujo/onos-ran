module github.com/onosproject/onos-ric

go 1.14

require (
	github.com/DataDog/mmh3 v0.0.0-20200316233529-f5b682d8c981
	github.com/atomix/go-client v0.1.0
	github.com/go-playground/overalls v0.0.0-20191218162659-7df9f728c018 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.3.1
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/joncalhoun/pipe v0.0.0-20170510025636-72505674a733
	github.com/mattn/goveralls v0.0.5 // indirect
	github.com/onosproject/helmit v0.6.0
	github.com/onosproject/onos-lib-go v0.6.1
	github.com/onosproject/onos-topo v0.6.0
	github.com/prometheus/client_golang v1.4.1
	github.com/spf13/cobra v0.0.6
	github.com/stretchr/testify v1.5.1
	github.com/yookoala/realpath v1.0.0 // indirect
	golang.org/x/tools v0.0.0-20200313205530-4303120df7d8 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools v2.2.0+incompatible
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
