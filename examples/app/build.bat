cd ..\..\gwacln
tinygo build --no-debug -target=wasm -o main.wasm .
copy main.wasm ..\examples\app\static
cd ..\examples\app
go build -ldflags="-H windowsgui"