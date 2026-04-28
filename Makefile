APP=diwe-web
TEMPL=go run github.com/a-h/templ/cmd/templ@v0.3.819

templ:
	$(TEMPL) generate

dev: templ
	go run ./cmd/web

build: templ
	go build -o bin/$(APP) ./cmd/web
