compile:
	echo "Compiling for every OS and Platform"
	mkdir -p bin
	GOOS=linux 	 GOARCH=arm   go build -C ./src -o ../bin/deraph-linux-arm  			 main.go
	GOOS=linux   GOARCH=arm64 go build -C ./src -o ../bin/deraph-linux-arm64  		 main.go
	GOOS=windows GOARCH=arm   go build -C ./src -o ../bin/deraph-windows-arm.exe   main.go
	GOOS=windows GOARCH=arm64 go build -C ./src -o ../bin/deraph-windows-arm64.exe main.go
	GOOS=windows GOARCH=amd64 go build -C ./src -o ../bin/deraph-windows-amd64.exe main.go
	GOOS=darwin  GOARCH=arm64 go build -C ./src -o ../bin/deraph-darwin-arm64  		 main.go
	GOOS=darwin  GOARCH=amd64 go build -C ./src -o ../bin/deraph-darwin-amd64  		 main.go
