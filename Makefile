test: test_php test_python test_golang

test_php:
	@printf "\n\033[92m\033[1mPHP\033[0m\033[0m\n"
	./vendor/bin/phpunit --colors=always --verbose --bootstrap ./tests/php/phpunit.php ./tests/php/

test_python:
	@echo "\n\033[92m\033[1mPYTHON\033[0m\033[0m"
	python3 -m unittest discover ./tests/python/ -p '*_test.py' -v

test_golang:
	@echo "\n\033[92m\033[1mGOLANG\033[0m\033[0m"
	rm -f go.mod
	go mod init example
	go clean -testcache
	go test ./golang

.PHONY: test test_php test_python test_golang