.PHONY: proto wire


proto:
	@echo "Generating protobuf code..."
	rm -rf gen/
	mkdir -p gen/go/event/v1
	protoc -I api --go_out=./gen/go --go_opt=paths=source_relative \
	  --go-grpc_out=./gen/go --go-grpc_opt=paths=source_relative \
	  api/proto/event/v1/email.proto

wire:
	cd internal/services/auth/app/di && wire
	cd internal/services/email/app/di && wire


create-kafka-topics:
	docker-compose exec kafka kafka-topics --create --topic user-events --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1

run:
	make proto
	make wire
	docker-compose -f deployments/docker-compose.yml up --build