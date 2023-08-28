init:
	if [ ! -d ./bin ]; then mkdir './bin';fi

run: build
	@echo '---------------------< Executing >---------------------'
	@./bin/jproj
	
build: init
	go build -o './bin/jproj'

clean:
	rm -rf './bin'

test-newproj: build clean-test-newproj
	mkdir 'test_playground'
	./bin/jproj createproject -b './test_playground' -n 'newproject'
	if [ ! -d './test_playground' ];then exit 1;fi
	@printf "class Program \n{\tpublic static void main(String[] args) {\n\t\tSystem.out.println(\"Hello World!\");\n\t}\n}" > './test_playground/newproject/src/App.java'

clean-test-newproj:
	rm -rf './test_playground'

test-build: test-newproj
	./bin/jproj build -d './test_playground/newproject'

