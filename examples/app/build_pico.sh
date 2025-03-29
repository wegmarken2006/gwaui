cd ../../gwacln
tinygo build -tags pico --no-debug -o main.wasm .
cp main.wasm ../examples/app/static
