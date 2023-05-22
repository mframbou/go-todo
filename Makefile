all: run

backend:
	cd backend && go run main.go

frontend:
	cd frontend && npm install && http-server -p 8080

run:
	make frontend & make backend

stop:
	pkill -f "go run main.go" && pkill -f "http-server -p 8080"

.PHONY: backend frontend run stop