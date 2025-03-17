cd ../../gwacln
tinygo build --no-debug -o main.wasm .
cp main.wasm ../examples/app/static