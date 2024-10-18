module app

go 1.23.1

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require github.com/xxx/gwasrv v1.0.0

replace github.com/xxx/gwasrv v1.0.0 => ../../gwasrv
