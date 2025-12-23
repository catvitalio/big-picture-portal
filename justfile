default: build

build: generate-rsrc
	GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o BigPicturePortal.exe
	@echo "Build complete: BigPicturePortal.exe"

generate-rsrc:
	@if [ ! -f rsrc.syso ] || [ app.manifest -nt rsrc.syso ]; then \
		echo "Generating resource file..."; \
		export PATH="$(go env GOPATH)/bin:$$PATH"; \
		if [ -f assets/icon.ico ]; then \
			rsrc -manifest app.manifest -ico assets/icon.ico -o rsrc.syso 2>/dev/null || \
			(echo "rsrc tool not found, installing..." && \
			 go install github.com/akavel/rsrc@latest && \
			 rsrc -manifest app.manifest -ico assets/icon.ico -o rsrc.syso); \
		else \
			rsrc -manifest app.manifest -o rsrc.syso 2>/dev/null || \
			(echo "rsrc tool not found, installing..." && \
			 go install github.com/akavel/rsrc@latest && \
			 rsrc -manifest app.manifest -o rsrc.syso); \
		fi; \
	fi

clean:
	rm -f BigPicturePortal.exe rsrc.syso
	@echo "Clean complete."

install-tools:
	go install github.com/akavel/rsrc@latest
	@echo "Tools installed."

rebuild: clean build
