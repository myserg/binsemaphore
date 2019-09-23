test:
	go test -v -count=1 .

testrace:
	go test -v -count=1 . -race

.PHONY: test testrace
