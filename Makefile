# Binary names (Windows targets only)
EXE_NAME=go-cookie-monster.exe
DLL_NAME=go-cookie-monster.dll

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Determine platform-specific variables
ifeq ($(OS),Windows_NT)
    # Windows-specific settings
    SHELL := powershell.exe
    .SHELLFLAGS := -NoProfile -Command
    RM_F := Remove-Item -Force -Recurse -ErrorAction Ignore
    # For Windows builds
    BUILD_CMD_EXE = $$env:CGO_ENABLED=1; $$env:GOARCH='amd64'; $$env:GOOS='windows'; $(GOBUILD)
    BUILD_CMD_DLL = $$env:CGO_ENABLED=1; $$env:GOARCH='amd64'; $$env:GOOS='windows'; $(GOBUILD) -buildmode=c-shared
else
    # Unix-like system settings (Linux/MacOS)
    RM_F := rm -f
    # Cross-compilation settings
    BUILD_CMD_EXE = GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(GOBUILD)
    BUILD_CMD_DLL = GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(GOBUILD) -buildmode=c-shared
endif

all: build-exe build-dll

build-exe exe:
	$(BUILD_CMD_EXE) -o $(EXE_NAME) .

build-dll dll:
	$(BUILD_CMD_DLL) -o $(DLL_NAME) ./dll

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	$(RM_F) $(EXE_NAME)
	$(RM_F) $(DLL_NAME)
	$(RM_F) $(DLL_NAME:.dll=.h)

run: build-exe
	./$(EXE_NAME)

deps:
	$(GOGET) ./...

.PHONY: all build-exe exe build-dll dll test clean run deps