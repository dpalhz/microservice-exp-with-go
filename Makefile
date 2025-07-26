.PHONY: proto wire


proto:
	protoc --go_out=./gen/go --go_opt=paths=source_relative \
	api/proto/event/v1/user_event.proto

wire:
	cd internal/services/auth/app/di && wire
	cd internal/services/email/app/di && wire


create-kafka-topics:
	docker-compose exec kafka kafka-topics --create --topic user-events --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1

run:
	make proto
	make wire
	docker-compose -f deployments/docker-compose.yml up --build