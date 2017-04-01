agni: *.go
	GOOS=linux go build -o agni *.go

clean:
	rm -f agni
