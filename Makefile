install-kustomize:
	GO111MODULE=on go get sigs.k8s.io/kustomize/kustomize/v4

SecretsFromGCP.so:
	mkdir -p dist/kustomize/plugin/pollination.cloud/v1/secretsfromgcp
	go build -buildmode plugin -o dist/kustomize/plugin/pollination.cloud/v1/secretsfromgcp/SecretsFromGCP.so kustomize/plugin/pollination.cloud/v1/secretsfromgcp/SecretsFromGCP.go

build-plugins: install-kustomize SecretsFromGCP.so

install-plugins: build-plugins
	mkdir -p ${HOME}/.config/kustomize/plugin
	cp -R dist/kustomize ${HOME}/.config/
