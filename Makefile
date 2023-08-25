init:
	if [ ! -d ./bin ]; then mkdir './bin';fi

run: build
	@echo '---------------------< Executing >---------------------'
	@./bin/jproj
	
build: init
	go build -o './bin/jproj'

clean:
	rm -rf './bin'

test-newproj: build
	mkdir 'createproj'
	./bin/jproj createproj -b './createproj' -n 'newproject'

clean-test-newproj:
	rm -rf './createproj'
