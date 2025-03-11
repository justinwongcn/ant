module example

go 1.24.0

replace github.com/roovet/ant => ../..

require (
	github.com/justinwongcn/ant v0.0.0-20250302091633-6b8363f63d03
	github.com/roovet/ant v0.0.0-00010101000000-000000000000
)

require github.com/hashicorp/golang-lru v1.0.2 // indirect
