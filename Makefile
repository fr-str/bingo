run-go-test:
	BINGO_DB_DIR=:memory: go run ./cmd/bingo/main.go
test:
	hurl --test ./hurl

tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/a-h/templ/cmd/templ@latest
	echo "for tests you need to install hurl, https://github.com/Orange-OpenSource/hurl?tab=readme-ov-file#installation"
