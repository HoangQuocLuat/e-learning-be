go get -d github.com/99designs/gqlgen@v0.17.44
GOTOOLCHAIN=go1.22.3 go run github.com/99designs/gqlgen generate --config ./user.gqlgen.yml

go get -d github.com/99designs/gqlgen@v0.17.44
GOTOOLCHAIN=go1.22.3 go run github.com/99designs/gqlgen generate --config ./admin.gqlgen.yml