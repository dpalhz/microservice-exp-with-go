.PHONY: proto wire

# Menghasilkan kode Go dari file.proto
proto:
	protoc --go_out=./gen/go --go_opt=paths=source_relative \
	api/proto/event/v1/user_event.proto

# Menghasilkan file injeksi dependensi untuk semua layanan
wire:
	cd internal/services/auth/app/di && wire
	cd internal/services/email/app/di && wire

# Perintah untuk membuat topik Kafka secara manual jika diperlukan
create-kafka-topics:
	docker-compose exec kafka kafka-topics --create --topic user-events --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1

# Menjalankan semuanya
run:
	make proto
	make wire
	docker-compose -f deployments/docker-compose.yml up --build