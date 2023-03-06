GOOS=linux GOARCH=amd64 go build -o bin/create-orders/main src/create-orders/createOrders.go
GOOS=linux GOARCH=amd64 go build -o bin/create-payments/main src/create-payments/createPayments.go
GOOS=linux GOARCH=amd64 go build -o bin/update-orders/main src/update-orders/updateOrders.go
GOOS=linux GOARCH=amd64 go build -o bin/update-payments/main src/update-payments/updatePayments.go