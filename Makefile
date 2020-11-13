JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

LDFLAGS		+= -s -w
LDFLAGS		+= -X github.com/lutomas/go-project-stub/types/version.appVersion=$(VERSION)
LDFLAGS		+= -X github.com/lutomas/go-project-stub/types/version.commit=$(GIT_REVISION)
LDFLAGS		+= -X github.com/lutomas/go-project-stub/types/version.buildTime=$(JOBDATE)

LDFLAGS_LINUX		+= -linkmode external -extldflags -static

# ###########################
# BUILD
# ###########################

install-linux:
	@echo "++ Building MAIN-APP binary (linux)"
	GOOS=linux CGO_ENABLED=1 go install -ldflags "$(LDFLAGS_LINUX) $(LDFLAGS)" github.com/lutomas/go-project-stub/cmd/main-app

install:
	@echo "++ Building MAIN-APP binary (<current-os>)"
	go install -ldflags "$(LDFLAGS)" github.com/lutomas/go-project-stub/cmd/main-app

install-cli:
	@echo "++ Building MAIN-APP-CLI binary (<current-os>)"
	go install -ldflags "$(LDFLAGS)" github.com/lutomas/go-project-stub/cmd/main-app-cli

image-main-app:
	docker build -t main-app-server:latest -f build/deploy/main-app-server.Dockerfile .

generate-main-app-api-types:
	@echo "++ Generating main-app API types:"
	oapi-codegen -package types -generate types api/open-api/main-app-open-api.yml > types/main_app_api.go

# ###########################
# DEV
# ###########################

go-mod-get:
	@echo "Get all dependencies"
	cd cmd/main-app && go get .
	cd cmd/main-app-cli && go get .

go-mod-vendor:
	@echo "Prepare and clean dependencies"
	go mod vendor && go mod tidy