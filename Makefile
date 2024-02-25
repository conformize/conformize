BIN_NAME := conformize

BUILD_DIR := "build"

SUPPORTED_OS := linux darwin windows freebsd netbsd openbsd

DARWIN_SUPPORTED_ARCHS := amd64 arm64
COMMON_SUPPORTED_ARCHS := 386 amd64 arm arm64

BUILD ?= dev

BUILD_FLAGS :=
ifeq ($(BUILD), release)
	gcflags = "all=-l -B"
	ldflags = "-w -s"
	BUILD_FLAGS :=-a -gcflags=$(gcflags) -ldflags=$(ldflags)
endif

.PHONY: all tidy lint check-target-platform check-build-option build-all build clean
all: build-all

tidy:
	@echo "\nTidying up Go modules..."
	@go mod tidy

lint:
	@echo "\nInstalling golint if not already installed..."
	@which golint > /dev/null || go install golang.org/x/lint/golint@latest
	@echo "\nRunning golint..."
	@golint ./...

check-build-option:
	@if [ "$(BUILD)" != "release" ] && [ "$(BUILD)" != "dev" ]; then \
		echo "\nError: BUILD must be set to either 'release' or 'dev' - got $(BUILD)"; \
		exit 1; \
	fi

check-target-platform:
	@if [ -n "$(OS)" ] && [ -n "$(ARCH)" ]; then \
		SUPPORTED=0; \
		echo "\nChecking support for target - $(OS)/$(ARCH)..."; \
		case "$(OS)" in \
			linux) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			darwin) ARCH_LIST="$(DARWIN_SUPPORTED_ARCHS)" ;; \
			windows) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			freebsd) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			netbsd) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			openbsd) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			*) echo "\nError: Unsupported ARCH $(ARCH)"; exit 1 ;; \
		esac; \
		for arch in $$ARCH_LIST; do \
			if [ "$(ARCH)" = "$$arch" ]; then \
				SUPPORTED=1; \
				break; \
			fi; \
		done; \
		if [ "$$SUPPORTED" -ne 1 ]; then \
			echo "\nError: Unsupported OS / ARCH target - $(OS)/$(ARCH)"; \
			exit 1; \
		fi; \
	fi

build-all: check-build-option
	@for os in $(SUPPORTED_OS); do \
		case "$$os" in \
			linux) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			darwin) ARCH_LIST="$(DARWIN_SUPPORTED_ARCHS)" ;; \
			windows) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			freebsd) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			netbsd) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
			openbsd) ARCH_LIST="$(COMMON_SUPPORTED_ARCHS)" ;; \
		esac; \
		for arch in $$ARCH_LIST; do \
			$(MAKE) OS=$$os ARCH=$$arch BUILD=$(BUILD) build; \
		done; \
	done

build: check-build-option check-target-platform tidy
	@echo "\nBuilding for OS=$(OS) ARCH=$(ARCH) BUILD=$(BUILD)..."
	@mkdir -p $(BUILD_DIR)/$(BUILD)$(if $(OS),/$(OS))$(if $(ARCH),/$(ARCH))
	@BIN_EXT=""; \
	if [ "$(OS)" = "windows" ]; then BIN_EXT=".exe"; fi; \
	BIN_PATH="$(BUILD_DIR)/$(BUILD)$(if $(OS),/$(OS))$(if $(ARCH),/$(ARCH))/$(BIN_NAME)$${BIN_EXT}"; \
	GOOS=$(OS) GOARCH=$(ARCH) go build -tags $(BUILD) $(BUILD_FLAGS) -o "$${BIN_PATH}"; \
	echo "\nDone building for OS=$(OS) ARCH=$(ARCH) BUILD=$(BUILD) BUILD_FLAGS=$(BUILD_FLAGS)"; \
	echo "\nBinary is available at ./$${BIN_PATH}"

clean:
	@echo "\nCleaning up..."
	@rm -rf $(BUILD_DIR)

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)
