init:
	if [ ! -d ./bin ]; then mkdir './bin';fi

build: init
	go build -o './bin/jproj'

test:
	go test ./utils
	go test ./configuration
	go test ./actions/jar
	go test ./templates

clean:
	rm -rf './bin'

test-newproj: build clean-test-newproj
	mkdir 'test_playground'
	./bin/jproj createproject -b './test_playground' -n 'newproject'
	if [ ! -d './test_playground' ];then exit 1;fi

clean-test-newproj:
	rm -rf './test_playground'

test-build: test-newproj
	./bin/jproj build -d './test_playground/newproject'

install: build
	cp ./bin/jproj ~/.local/bin
