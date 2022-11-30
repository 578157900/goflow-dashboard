module github.com/s8sg/goflow-dashboard

go 1.13

require (
	github.com/gorilla/mux v1.8.0
	github.com/onsi/ginkgo v1.14.1 // indirect
	github.com/onsi/gomega v1.10.2 // indirect; indirectxw
	github.com/rs/xid v1.2.1
	github.com/s8sg/faas-flow v0.6.2
	github.com/s8sg/goflow v0.0.9-0.20200927104759-0563d6d6dc7a
	gopkg.in/redis.v5 v5.2.9
)

replace github.com/s8sg/goflow v0.0.9-0.20200927104759-0563d6d6dc7a => gitlab.zenlayer.net/gia/go-flow v0.1.1-0.20221107093412-5f4fe5fd8d97

replace github.com/faasflow/sdk v1.0.0 => gitlab.zenlayer.net/Otis/faasflow-sdk v1.0.4

replace github.com/faasflow/faas-flow-redis-datastore v1.0.1-0.20200718081732-431d3cc7894a => github.com/578157900/faas-flow-redis-datastore v1.0.1-0.20201228095629-7bd1ddb5fba0

replace github.com/faasflow/faas-flow-redis-statestore v1.0.1-0.20200718082116-d90985fdbde1 => github.com/578157900/faas-flow-redis-statestore v1.0.1-0.20201228095149-84bf0e8f268e
