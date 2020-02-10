SRCDIR := $(shell pwd)
DSTDIR := dist

export index_html
override define index_html
<!DOCTYPE html>
<script src="wasm_exec.js"></script><script>
(async () => {
	const resp = await fetch('main.wasm');
	if (!resp.ok) {
		const pre = document.createElement('pre');
		pre.innerText = await resp.text();
		document.body.appendChild(pre);
		return;
	}
	const src = await resp.arrayBuffer();
	const go = new Go();
	const result = await WebAssembly.instantiate(src, go.importObject);
	go.run(result.instance);
})();
</script>
endef

run:
	wasmserve

deploy:
	@echo "deploy: $(SRCDIR) -> $(DSTDIR)"
	mkdir -p $(DSTDIR)
	GO111MODULE=on GOOS=js GOARCH=wasm go build -o $(DSTDIR)/main.wasm .
	echo "$$index_html" > $(DSTDIR)/index.html
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $(DSTDIR)/wasm_exec.js
	cp -f contents.md $(DSTDIR)/
	cp -Rf css $(DSTDIR)/
