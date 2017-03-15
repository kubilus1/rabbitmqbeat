BEATNAME=rabbitmqbeat
BEAT_DIR=github.com/bernielomax
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

# Path to the libbeat Makefile
-include $(ES_BEATS)/metricbeat/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	make create-metricset
	make collect

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	rm -rf vendor/github.com/elastic/beats/.git

.PHONY: update-deps
update-deps:
	glide update --no-recursive --strip-vcs

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

deb_pkg: rabbitmqbeat
	mkdir -p package/usr/share/rabbitmqbeat/bin
	mkdir -p package/etc/rabbitmqbeat
	mkdir -p package/var/lib/rabbitmqbeat
	mkdir -p package/var/log/rabbitmqbeat
	mkdir -p package/etc/systemd/system
	cp rabbitmqbeat package/usr/share/rabbitmqbeat/bin/.
	cp rabbitmqbeat.yml package/etc/rabbitmqbeat/.
	cp rabbitmqbeat.template.json package/etc/rabbitmqbeat/.
	cp rabbitmqbeat.service package/etc/systemd/system/.
	cp -r DEBIAN package/.
	sudo dpkg-deb --build package
	mv package.deb rabbitmqbeat.deb
