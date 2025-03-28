module app

go 1.23.1

require github.com/wegmarken2006/gwaui/gwasrv v1.0.0

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/webview/webview_go v0.0.0-20240831120633-6173450d4dd6 // indirect
	golang.org/x/sys v0.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//require github.com/wegmarken2006/gwaui/gwasrv v0.0.0-20250221070208-32883a805718

replace github.com/wegmarken2006/gwaui/gwasrv v1.0.0 => ../../gwasrv
