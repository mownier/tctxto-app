prepare:
	git submodule init && git submodule update --init --recursive && cd libs/tctxto-client/frontend && npm install && cd ../../tctxto-proxy && git submodule init && git submodule update --init --recursive

build-macos:
	rm -rf build/macos && mkdir -p build/macos && \
	cd libs/tctxto-server && make build-macos && \
	cd ../tctxto-proxy && make build-macos && \
	cd ../tctxto-client && make build-macos && \
	cd ../../ && \
	mv libs/tctxto-server/build/macos/* build/macos/ && \
	mv libs/tctxto-proxy/build/macos/* build/macos/ && \
	mv libs/tctxto-client/build/macos/* build/macos/

build-linux:
	rm -rf build/linux && mkdir -p build/linux && \
	cd libs/tctxto-server && make build-linux && \
	cd ../tctxto-proxy && make build-linux && \
	cd ../tctxto-client && make build-linux && \
	cd ../../ && \
	mv libs/tctxto-server/build/linux/* build/linux/ && \
	mv libs/tctxto-proxy/build/linux/* build/linux/ && \
	mv libs/tctxto-client/build/linux/* build/linux/

build-windows:
	rm -rf build/windows && mkdir -p build/windows && \
	cd libs/tctxto-server && make build-windows && \
	cd ../tctxto-proxy && make build-windows && \
	cd ../tctxto-client && make build-windows && \
	cd ../../ && \
	mv libs/tctxto-server/build/windows/* build/windows/ && \
	mv libs/tctxto-proxy/build/windows/* build/windows/ && \
	mv libs/tctxto-client/build/windows/* build/windows/

build-all: build-macos build-linux build-windows

run-macos-prepare:
	cd libs/tctxto-server && make build-macos && \
	cd ../tctxto-proxy && make build-macos && \
	cd ../tctxto-client && make build-macos && \
	cd ../../ && \
	mv libs/tctxto-server/build/macos/* . && \
	mv libs/tctxto-proxy/build/macos/* . && \
	mv libs/tctxto-client/build/macos/* .

run-linux-prepare:
	cd libs/tctxto-server && make build-linux && \
	cd ../tctxto-proxy && make build-linux && \
	cd ../tctxto-client && make build-linux && \
	cd ../../ && \
	mv libs/tctxto-server/build/linux/* . && \
	mv libs/tctxto-proxy/build/linux/* . && \
	mv libs/tctxto-client/build/linux/* .

run-windows-prepare:
	cd libs/tctxto-server && make build-windows && \
	cd ../tctxto-proxy && make build-windows && \
	cd ../tctxto-client && make build-windows && \
	cd ../../ && \
	mv libs/tctxto-server/build/windows/* . && \
	mv libs/tctxto-proxy/build/windows/* . && \
	mv libs/tctxto-client/build/windows/* .

run-app:
	go build && ./tctxtoapp

run-app-w:
	go build -o tctxtoapp.exe  && ./tctxtoapp.exe

run-macos: run-macos-prepare run-app

run-linux: run-linux-prepare run-app

run-windows: run-windows-prepare run-app-w
