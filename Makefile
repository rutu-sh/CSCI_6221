build-service:
	cd src/service && \
	$(MAKE) build

clean-service:
	cd src/service && \
	$(MAKE) clean

build-alerter:
	cd src/alerter && \
	$(MAKE) build

clean-alerter:
	cd src/alerter && \
	$(MAKE) clean

build-auth:
	cd src/auth && \
	$(MAKE) build

clean-auth:
	cd src/auth && \
	$(MAKE) clean

