# Scope

This repository contains a query language for use with Go Buffalo.

# Motivation

See Blog post.

# Integration
This is a go package.  Integrate this repository by running 
`go get github.com/alphaflow/scope` within an executable go project.

# Usage and Examples

Example parameters can be found in `USAGE.md`.

# Local Development

You will need to have a postgres database running in order to run tests.  Set `TEST_DATABASE_URL` in your environment and then run:

```
 buffalo pop create -e test;
 buffalo pop migrate -e test;
```

To create the scope test database with buffalo pop.

When you have completed your changes, please read `DEPLOYMENT.md` to merge your changes.
