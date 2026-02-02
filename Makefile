user_service:
	cd user-service && go run cmd/main.go

auth_service:
	cd auth-service && go run cmd/main.go

gateway:
	cd gateway-service && go run cmd/main.go