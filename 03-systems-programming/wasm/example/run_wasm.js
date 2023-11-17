// JS functions callable from within WASM code
// 2 layers of nesting for importObject[module][field]
var importObject = {
  console: {
    log: function(arg) {
      console.log(arg);
    }
  }
};

// Boilerplate to fetch the binary .wasm file, instantiate it as a
// WebAssembly module, then run the callback on the module
function loadAndRun(modulePath, callback) {
  fetch(modulePath).then(response =>
    response.arrayBuffer()
  ).then(bytes =>
    WebAssembly.instantiate(bytes, importObject)
  ).then(module => {
    callback(module.instance.exports)
  });
}

loadAndRun("fib.wasm", function(exports) {
  console.log(exports.fib(10));
});
