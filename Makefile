user_service:
	cd user-service && go run cmd/main.go

auth_service:
	cd auth-service && go run cmd/main.go

order_service:
	cd order-service && go run cmd/main.go

payment_service:
	cd payment-service && go run cmd/main.go

gateway:
	cd gateway-service && go run cmd/main.go