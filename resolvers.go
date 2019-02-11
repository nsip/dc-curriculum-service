// resolvers.go

//
// resolver functions for
// returning curriculum data through
// graphql
//

package main

import (
	"fmt"

	graphql "github.com/playlyfe/go-graphql"
)

//
// creates the set of resolvers available to the executor
//
func buildResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	resolvers["syllabus/overview"] = func(params *graphql.ResolveParams) (interface{}, error) {
		lookupKey := deriveKey(params)
		return getJSONMap(lookupKey)
	}

	resolvers["syllabus/content"] = func(params *graphql.ResolveParams) (interface{}, error) {
		lookupKey := deriveKey(params)
		return getJSONMap(lookupKey)
	}

	return resolvers
}

//
// derives the json document lookup key from the values
// passed in the gql query parameters
//
func deriveKey(params *graphql.ResolveParams) string {

	args := params.Args
	docType := params.Field.Name.Value

	key := fmt.Sprintf("%s-%s-%s-stage%s-%s", args["state"], args["learning_area"],
		args["subject"], args["stage"], docType)

	// log.Println("lookup key is:\t", key)

	return key
}
