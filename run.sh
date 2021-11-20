#!/bin/bash
set -o allexport; source .env; set +o allexport
echo "Updating API documentation..."
swag init
go build -o gocommerce
./gocommerce
