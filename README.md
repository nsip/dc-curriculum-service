# dc-curriculum-service
small graphql service that supplies syllabus information


#### install:
`go get github.com/nsip/dc-curriculum-service`

if you are using golang then build as normal then run with:
`./dc-curriculum-service` on the command-line

this will fire up a graphql endpoint on:

`http://localhost:1330/graphql`

if you are not building in golang, download the binary release for your platform from the Releases tab, and run from the command-line as above.

schema is self-documenting, so to explore non-programmatically download a graphql client and just point it at the endpoint.

The desktop graphiql is quick and easy, and can be downloaded from:

https://github.com/skevy/graphiql-app/releases

or if you are on Mac, you can use homebrew:

`brew cask install graphiql`

which will add the query explorer as an App on your system.

if you are a VS user the Apollo client has a plugin for your IDE:

https://marketplace.visualstudio.com/items?itemName=apollographql.vscode-apollo



