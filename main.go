package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
)

func main() {
	schema, err := graphql.NewSchema(defineSchema())
	if err != nil {
		log.Panic("Error when creating the graphQL schema", err)
	}

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.Handle("/graphql", h)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic("Error when starting the http server", err)
	}
	log.Print("Server is started")
}
