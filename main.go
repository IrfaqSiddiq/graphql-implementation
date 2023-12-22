package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Book struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Isbn     string `json:"isbn"`
	Quantity int    `json:"quantity"`
}

var books = []Book{
	{ID: 1, Title: "Book 1", Author: "Author 1", Isbn: "1234567890", Quantity: 5},
	{ID: 2, Title: "Book 2", Author: "Author 2", Isbn: "0987654321", Quantity: 3},
	// Add more books as needed
}

func setupHelloGraphQL() *handler.Handler {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
}

func setupBookByIDGraphQL() *handler.Handler {
	bookType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Book",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.Int},
			"title":    &graphql.Field{Type: graphql.String},
			"author":   &graphql.Field{Type: graphql.String},
			"isbn":     &graphql.Field{Type: graphql.String},
			"quantity": &graphql.Field{Type: graphql.Int},
		},
	})

	fields := graphql.Fields{
		"bookByID": &graphql.Field{
			Type: bookType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, ok := params.Args["id"].(int)
				if !ok {
					return nil, nil
				}
				for _, book := range books {
					if book.ID == id {
						return book, nil
					}
				}
				return nil, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
}

func main() {
	r := gin.Default()

	// Endpoint for 'hello' query
	r.POST("/hello/graphql", func(c *gin.Context) {
		h := setupHelloGraphQL()
		h.ServeHTTP(c.Writer, c.Request)
	})

	// Endpoint for 'bookByID' query
	r.POST("/bookByID/graphql", func(c *gin.Context) {
		h := setupBookByIDGraphQL()
		h.ServeHTTP(c.Writer, c.Request)
	})

	log.Println("Server started at http://localhost:8081")
	log.Fatal(r.Run(":8081"))
}
