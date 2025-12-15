package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ludovicplt/muscubackend/config"
	"github.com/ludovicplt/muscubackend/graph"
	"github.com/ludovicplt/muscubackend/graph/model"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
		return err
	}
	return nil
}

func gqlInit(db *gorm.DB) *handler.Server {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func main() {
	r := gin.Default()

	if err := LoadEnv(); err != nil {
		log.Println("Continuing without .env file")
	}

	db := config.ConnectDB(config.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
	})

	// Migrate the schema
	db.Migrator().DropTable(&model.User{})
	db.AutoMigrate(&model.User{})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := gqlInit(db)

	r.POST("/query", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL playground", "/query").ServeHTTP(c.Writer, c.Request)
	})

	r.Run(":" + port)
}
