build:
	@go build -o bin/hola main.go
test:
	go run main.go examples/single.http
	go run main.go examples/users.http -n 1
	X_API_SECRET=GIMME-ACCESS X_API_KEY=THIS-IS-MY-KEY go run main.go examples/users.http -n 2
flake:
	nix build .#hola

