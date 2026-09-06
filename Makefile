# ==============================================================================
# LatinaServer — Makefile
# ==============================================================================

.PHONY: all build build-web build-server test test-web clean run dev-web

all: build

# 1. Build All: Static assets with Hugo Extended, then compile Go server
build: build-web build-server

# 2. Build Web Frontend (Hugo static site generator into web/dist)
build-web:
	hugo --minify -s web -d dist

# 3. Build Go Server Binary
build-server:
	go build -v -o latinaserver ./cmd/latinaserver

# 4. Run all unit and integration tests
test:
	go test ./...

# 5. Run Web Service Tests
test-web:
	go test ./internal/service/web/... -v

# 6. Run Hugo Local Development Server with Hot Reload
dev-web:
	hugo server -s web -D

# 7. Clean build artifacts
clean:
	rm -f latinaserver latinaserver.exe
