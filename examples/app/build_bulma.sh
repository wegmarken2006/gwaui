cd ../../gwacln
tinygo build -tags bulma --no-debug -o main.wasm .
cp main.wasm ../examples/app/static
