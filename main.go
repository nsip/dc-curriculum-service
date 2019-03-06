//
// experimental package to host nsw curriculum data in
// small embedded db and serve using graphql
//

package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	graphql "github.com/playlyfe/go-graphql"
)

var syllabus_executor *graphql.Executor
var search_executor *graphql.Executor

func main() {

	// ensure database is available
	initialiseDB()
	defer closeDB()

	// construct gql resolvers & schema
	resolvers := buildResolvers()
	schema := buildSchema()

	// initialise the qgl executors
	var executorErr error
	syllabus_executor, executorErr = graphql.NewExecutor(schema, "syllabus", "", resolvers)
	if executorErr != nil {
		log.Fatal("cannot create syllabus executor: ", executorErr)
	}
	search_executor, executorErr = graphql.NewExecutor(schema, "websearch", "", resolvers)
	if executorErr != nil {
		log.Fatal("cannot create web-search executor: ", executorErr)
	}

	// start the gql web server
	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})) // allow cors requests during testing

	// the graphql handlers
	e.POST("/graphql", gqlHandlerSyllabus)
	e.POST("/search/graphql", gqlHandlerSearch)

	// run the server
	e.Logger.Fatal(e.Start(":1330"))

}

//
// wrapper type to capture graphql input
//
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

//
// the graphql handler routine for sylabus requests
//
func gqlHandlerSyllabus(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}

	result, err := syllabus_executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	// log.Printf("result:\n\n%#v\n\n", result)

	return c.JSON(http.StatusOK, result)

}

//
// the graphql handler for web-search requests
//
func gqlHandlerSearch(c echo.Context) error {

	grq := new(GQLRequest)
	if err := c.Bind(grq); err != nil {
		return err
	}

	query := grq.Query
	variables := grq.Variables
	gqlContext := map[string]interface{}{}

	result, err := search_executor.Execute(gqlContext, query, variables, "")
	if err != nil {
		panic(err)
	}

	// log.Printf("result:\n\n%#v\n\n", result)

	return c.JSON(http.StatusOK, result)

}
