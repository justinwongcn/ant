module github.com/justinwongcn/ant/examples/router

go 1.24.0

toolchain go1.24.2

require (
	github.com/hashicorp/golang-lru v1.0.2
	github.com/justinwongcn/ant v0.0.1
)

require github.com/patrickmn/go-cache v2.1.0+incompatible // indirect

replace github.com/justinwongcn/ant => ../..
