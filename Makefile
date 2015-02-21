BINARY_IMAGE_NAME=resm-go-binary:v1.0
BINARY_CONTAINER_NAME=resm-go-binary
DEB_BUILDER_IMAGE_NAME=resm-deb-builder:v1.0
DEB_BUILDER_CONTAINER_NAME=resm-deb-builder
SUPERVISOR_IMAGE_NAME=resm-supervisor:v1.0
SUPERVISOR_CONTAINER_NAME=resm-supervisor

run: build
	${GOPATH}/bin/resm -limit=3

run_bolt: build
	${GOPATH}/bin/resm -limit=3 -file=bolt.db

build:
	go install github.com/nordicdyno/resm-sketch/resm

test:
	go test github.com/nordicdyno/resm-sketch/store/inmemory
	go test github.com/nordicdyno/resm-sketch/store/inbolt
	go test github.com/nordicdyno/resm-sketch/resm

clean:
	rm -rf .src
	find . -name '*.db' -exec rm {} \;

docker_clean:
	# Remove all untagged images
	./docker/rmi_clean.sh

docker_build_bin:
	-docker rm -f ${BINARY_CONTAINER_NAME}
	#resm-debian-runner resm-build-deb
	docker build --rm --tag=${BINARY_IMAGE_NAME} -f docker/debian_binary_build.Dockerfile .
	docker run -v /src --name=${BINARY_CONTAINER_NAME} ${BINARY_IMAGE_NAME}
	# docker cp /src/bin/resm bin/resm

docker_supervisor: docker_build_bin
	-docker rm -f ${SUPERVISOR_CONTAINER_NAME}
	docker build -rm --tag=${SUPERVISOR_IMAGE_NAME} -f docker/debian_supervisord.Dockerfile docker/
	docker run -it --net=host --volumes-from ${BINARY_CONTAINER_NAME} --name=${SUPERVISOR_CONTAINER_NAME} ${SUPERVISOR_IMAGE_NAME}

docker_build_deb: docker_build_bin
	-docker rm -f ${DEB_BUILDER_CONTAINER_NAME}
	docker build --tag=${DEB_BUILDER_IMAGE_NAME} -f docker/debian_fmp_deb.Dockerfile docker/
	docker run --volumes-from ${BINARY_CONTAINER_NAME} --name=${DEB_BUILDER_CONTAINER_NAME} ${DEB_BUILDER_IMAGE_NAME}
	# INFO:
	# "Now you can copy deb package from resm-fpm-deb-builder "
	# " steps depends on your environment, but final ommand would be same: "
	# "docker cp ${DEB_BUILDER_CONTAINER_NAME}:/root/resm/resm-go_1.0_amd64.deb ./"
