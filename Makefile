run-s:
	cd tctxto-server && make run-me && cd ..

run-p:
	cd tctxto-proxy && make run AO=http://localhost:2323 && cd ..
	
run-c:
	cd tctxto-client && make run && cd ..

run-all:
	make run-s & make run-p & make run-c
