module github.com/onosproject/onos-ric

go 1.14

require (
	github.com/DataDog/mmh3 v0.0.0-20200316233529-f5b682d8c981
	github.com/atomix/go-client v0.2.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.3
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/joncalhoun/pipe v0.0.0-20170510025636-72505674a733
	github.com/onosproject/helmit v0.6.4
	github.com/onosproject/onos-lib-go v0.6.5
	github.com/onosproject/onos-topo v0.6.3
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/prometheus/client_golang v1.4.1
	github.com/spf13/cobra v0.0.6
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/testify v1.5.1
	go.uber.org/multierr v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20200212091648-12a6c2dcc1e4 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20200313205530-4303120df7d8 // indirect
	google.golang.org/grpc v1.29.1
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools v2.2.0+incompatible
)

replace github.com/onosproject/onos-lib-go => ../onos-lib-go
replace github.com/atomix/go-client => ../atomix-go-client

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
