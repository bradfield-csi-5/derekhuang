### WebAssembly

Implement a compiler that generates [WebAssembly text format (WAT)](https://developer.mozilla.org/en-US/docs/WebAssembly/Understanding_the_text_format), starting from a Lox AST (or a Go AST, if you're feeling more ambitious).

You can use [wat2wasm](https://github.com/WebAssembly/wabt) (local) or [WebAssembly Explorer](https://mbebenita.github.io/WasmExplorer/) (online) to compile your module from a `.wat` text format file to a `.wasm` binary format file.

You can then use your compiled `.wasm` file directly from Node.js; for example, if `tmp.wasm` has a `main` function (implemented on the WebAssembly side) that imports and calls `print`, then you can run `main` on the JavaScript side as follows:

```javascript
const fs = require('fs').promises;

async function run() {
  async function createWebAssembly(bytes) {
    const env = {
      print: x => console.log(x)
    };
    return WebAssembly.instantiate(bytes, { env });
  }

  const result = await createWebAssembly(new Uint8Array(await fs.readFile('tmp.wasm')));
  console.log(result.instance.exports.main());
}
run();
```
