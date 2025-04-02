module github.com/salaboy/devex-on-k8s/chapter03/end/go/appointments

go 1.21

toolchain go1.21.0

require (
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/chi/v5 v5.0.10
	github.com/google/uuid v1.3.1
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cucumber/godog => github.com/laurazard/godog v0.0.0-20220922095256-4c4b17abdae7
