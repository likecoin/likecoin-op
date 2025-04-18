.PHONY: decrypt-everything
decrypt-everything:
	blackbox_postdeploy

.PHONY: operator-key-link
operator-key-link:
	ln -s ../deploy/env.operator likecoin/.env
	ln -s ../deploy/env.operator likenft/.env

.PHONY: remove-operator-key-link
remove-operator-key-link:
	rm -f likecoin/.env
	rm -f likenft/.env

.PHONY: local-node
local-node:
	(sleep 1 && make -C operation init-local-state) &
	(sleep 2 && make -C likecoin deploy-local) &
	(sleep 3 && make -C likenft deploy-local) &
	make -C operation local-node

.PHONY: abigen
abigen:
	make -C likenft build
	make -C likecoin build
	mkdir -p abi
	cp likenft/artifacts/contracts/LikeProtocol.sol/LikeProtocol.json abi/
	cp likenft/artifacts/contracts/BookNFT.sol/BookNFT.json abi/
	cp likecoin/artifacts/contracts/EkilCoin.sol/EkilCoin.json abi/
	make -C likenft-indexer abigen
	make -C migration-backend abigen

.PHONY: docker-images
docker-images:
	DOCKER_BUILD_ARGS=--push make -C migration-backend docker-image
	DOCKER_BUILD_ARGS=--push make -C signer-backend docker-image
	DOCKER_BUILD_ARGS=--push make -C likenft-indexer docker-image
	DOCKER_BUILD_ARGS=--push make -C likenft-migration docker-image
	DOCKER_BUILD_ARGS=--push make -C likecoin-migration docker-image
	DOCKER_BUILD_ARGS=--push make -C migration-admin docker-image

.PHONY: deploy
deploy: decrypt-everything
	make -C deploy deploy
