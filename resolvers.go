// resolvers.go

//
// resolver functions for
// returning curriculum data through
// graphql
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	graphql "github.com/playlyfe/go-graphql"
)

//
// creates the set of resolvers available to the executor
//
func buildResolvers() map[string]interface{} {

	resolvers := map[string]interface{}{}

	//
	// resolvers for sylabus queries
	//
	resolvers["syllabus/overview"] = func(params *graphql.ResolveParams) (interface{}, error) {
		lookupKey := deriveKey(params)
		return getJSONMap(lookupKey)
	}

	resolvers["syllabus/content"] = func(params *graphql.ResolveParams) (interface{}, error) {
		lookupKey := deriveKey(params)
		return getJSONMap(lookupKey)
	}

	//
	// resolver for web-seach queries
	//
	resolvers["websearch/searchRequest"] = func(params *graphql.ResolveParams) (interface{}, error) {

		searchTerms := deriveSearchTerms(params)
		// log.Println("search terms: ", searchTerms)

		const endpoint = "https://api.cognitive.microsoft.com/bing/v7.0/search"
		token := "72754c1352f04b4f96f9a1b46d854257"

		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			log.Println(err)
		}

		param := req.URL.Query()
		param.Add("q", searchTerms)
		req.URL.RawQuery = param.Encode()

		req.Header.Add("Ocp-Apim-Subscription-Key", token)

		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		ans := new(BingAnswer)
		err = json.Unmarshal(body, &ans)
		if err != nil {
			return nil, err
		}

		return ans, nil
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

//
// parses the parameters to construct a set of search terms to pass to the
// web search engine
//
func deriveSearchTerms(params *graphql.ResolveParams) string {

	// check if an input block was supplied
	input_map, ok := params.Args["terms"].(map[string]interface{})
	if !ok {
		return ""
	}

	// now grab the variables passed
	terms := make([]string, 0)
	kws, ok := input_map["keywords"].([]interface{})
	if ok {
		for _, keyword := range kws {
			terms = append(terms, keyword.(string))
		}
	}
	learning_area, ok := input_map["learning_area"].(string)
	if ok {
		terms = append(terms, learning_area)
	}
	subject, ok := input_map["subject"].(string)
	if ok {
		terms = append(terms, subject)
	}
	stage, ok := input_map["stage"].(string)
	if ok {
		terms = append(terms, "stage "+stage)
	}
	course, ok := input_map["course_name"].(string)
	if ok {
		terms = append(terms, course)
	}
	content_area, ok := input_map["content_area"].(string)
	if ok {
		terms = append(terms, content_area)
	}

	// log.Println("terms: ", terms)

	//make sure the search terms are within size tolerance for bing search
	allTerms := strings.Join(terms, " ")
	if len(allTerms) > 1024 {
		safeTerms := allTerms[0:1023]
		return safeTerms
	} else {
		return allTerms
	}

}
