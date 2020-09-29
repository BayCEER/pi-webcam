version=1.0.0
project=pi-webcam
os=linux
arch=arm
arm=6

.DEFAULT_GOAL := package

clean:
	rm -rf target
	
prepare: clean
	mkdir target
	cp -r deb target/deb
	sed -i 's/\[\[version\]\]/${version}/g' target/deb/DEBIAN/control
	
build: prepare
	env GOOS=$(os) GOARCH=$(arch) GOARM=$(arm) go build -ldflags="-s -w" -o target/$(project) main.go

compress: build
	upx -o target/deb/usr/bin/$(project) target/$(project)

package: compress
	# Fails on WSL if not mounted with metadata, please add metaoption to /etc/wsl.conf
	dpkg-deb -b target/deb target/$(project)-$(version)_$(arch).deb

install: package
	dpkg -i target/$(project)-$(version)_$(arch).deb
	systemctl restart pi-webcam.service


