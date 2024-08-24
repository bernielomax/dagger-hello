.PHONY: test
test:
	dagger call test --source=.

.PHONY: build
build:
	dagger call build --source=.

.PHONY: publish
publish:
	dagger call publish --source=.
