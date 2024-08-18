.PHONY: default
default:
	@echo "type an action to perform"

VERSION=v0.0.1
.PHONY: tag
tag:
	cd cmd && sed 's/v[0-9]*\.[0-9]*\.[0-9]*/$(VERSION)/g' ver.go	
	git pull && git add . && git commit -m "release $(VERSION)" && git push
	git tag -a $(VERSION) -m "release $(VERSION)"
	git push --tags