run-s:
	cd tctxto-server && make run-me && cd ..

run-p:
	cd tctxto-proxy && make run AO=http://localhost:2323 && cd ..
	
run-c:
	cd tctxto-client && make run && cd ..

run-all:
	make run-s & make run-p & make run-c

clone-s:
	mkdir tctxto-server && git clone https://github.com/mownier/tctxto-server.git tctxto-server/

clone-p:
	mkdir tctxto-proxy && git clone https://github.com/mownier/tctxto-proxy.git tctxto-proxy/ && cd tctxto-proxy/ && make clone && cd ..

clone-c:
	mkdir tctxto-client && git clone https://github.com/mownier/tctxto-client.git tctxto-client/

setup-p:
	cd tctxto-proxy/grpc-web/go/grpcwebproxy/ && openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -subj "/CN=localhost" && openssl rsa -in key.pem -out decrypted_key.pem && cd ../../../../

init:
	make clone-s && make clone-p && make clone-c && make setup-p
