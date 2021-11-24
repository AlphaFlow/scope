![Logo](https://www.alphaflow.com/wp-content/themes/Alphaflow/res/alphaflow-logo.png) 
# Scope
The primary purpose of the scope package is to provide scopes that can be easily attached to Pop/Gorm (See gorm package) queries.  Think `tx.Scope([scopeFunc]).All(&books)`.

This repository also contains a limited "query language" for use with Go Buffalo.  These scopes take query parameters and produce scope functions to generically manipulate `pop.Models`, see `scope.For[Filter|Sort|Paginate]FromParams(...)`. The `GetAggregationsFromParams` function serves a similar purpose of taking a set of query params and providing generic output, however because the result structure is different the aggregations execute a query to return the results for you.

# Motivation
See Blog post.

# Integration
This is a go package.  Integrate this repository by running 
`go get github.com/alphaflow/scope` within an executable go project.

# Usage and Examples

A description of all available parameters can be found in [USAGE.md](/USAGE.md).

An example buffalo resource file can be found in [/examples/resource.go](/examples/resource.go).

# Local Development

You will need to have a postgres database running in order to run tests.  Set `TEST_DATABASE_URL` in your environment and then run:

```
 buffalo pop create -e test;
 buffalo pop migrate -e test;
```

To create the scope test database with buffalo pop.

When you have completed your changes, please read [DEPLOYMENT.md](/DEPLOYMENT.md) to merge your changes.
