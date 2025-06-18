MINING_SERVICE_DIR=apps/mining-service
SHAPING_SERVICE_DIR=apps/shaping-service

# Service
install-mining-service:
	cd $(MINING_SERVICE_DIR) && go mod tidy
install-shaping-service:
	if [ ! -d "$(SHAPING_SERVICE_DIR)/.venv" ]; then \
		cd $(SHAPING_SERVICE_DIR) && uv venv; \
	fi
	cd $(SHAPING_SERVICE_DIR) && source .venv/bin/activate
	cd $(SHAPING_SERVICE_DIR) && uv pip install -r requirements.txt
dev-mining-service:
	cd $(MINING_SERVICE_DIR) && go run cmd/main.go
dev-shaping-service:
	$(SHAPING_SERVICE_DIR)/.venv/bin/python -m uvicorn apps.shaping-service.main:app --reload --host 0.0.0.0 --port 4001
build-mining-service:
	cd $(MINING_SERVICE_DIR) && go build -o main cmd/main.go
docker-build-mining-service:
	cd $(MINING_SERVICE_DIR) && docker build -t mining-service:latest .
docker-build-shaping-service:
	cd $(SHAPING_SERVICE_DIR) && docker build -t shaping-service:latest .

# all
# dev:
# 	$(MAKE) dev-mining-service &
# 	$(MAKE) dev-shaping-service
dev: dev-mining-service dev-shaping-service # make -j2 dev
build: build-mining-service
docker-build: docker-build-mining-service docker-build-shaping-service

.PHONY: install-mining-service install-shaping-service dev-mining-service dev-shaping-service build-mining-service build dev docker-build-mining-service docker-build-shaping-service docker-build