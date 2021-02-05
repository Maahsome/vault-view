.PHONY: check_release_version
check_release_version:
ifeq (,$(RELEASE_VERSION))
	$(error "RELEASE_VERSION must be set to a release tag")
endif

.PHONY: changelog
changelog: check_release_version ## Generate the changelog.
	@mkdir -p changelog/releases && rm -f changelog/releases/$(RELEASE_VERSION).md
	go run ./release/changelog/gen-changelog.go -tag=$(RELEASE_VERSION) -changelog=changelog/releases/$(RELEASE_VERSION).md
	rm -f ./changelog/fragments/!(00-template.yaml)
