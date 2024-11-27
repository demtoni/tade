all: app manager

vet:
	go vet ./...

app:
	cd webapp && npm run build
	CGO_ENABLED=0 go build -o app cmd/app/main.go

manager:
	CGO_ENABLED=0 go build -o manager cmd/manager/main.go
