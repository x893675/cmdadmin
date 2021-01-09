build-all:
	mkdir -pv dist && KUBE_VERBOSE=2 bash hack/make-rules/build.sh


# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...