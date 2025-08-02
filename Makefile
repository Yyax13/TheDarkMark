.PHONY: run build clean help

OS := $(shell uname -s 2>/dev/null || echo Windows_NT)

ifeq ($(OS),Windows_NT)
    RM = powershell -Command "Remove-Item -Recurse -Force"
    RMBUILD = build\*
else
    RM = rm
    RMBUILD = -rf ./build/
endif

run: ## Run the project without compiling
	go run ./src

build: ## Compile and sabe execlutable into ./build
ifeq ($(OS),Windows_NT)
	powershell -ExecutionPolicy Bypass -File scripts\build.ps1
else
	chmod +x scripts/build.sh
	./scripts/build.sh
endif

clean: ## Delete the build directory
	$(RM) $(RMBUILD)

help: ## Shows that help menu
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-13s\033[0m %s\n", $$1, $$2}'
