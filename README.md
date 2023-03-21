# breakdown

breaks down packages within repos

Owned by eng-infra

## Developing

- Update swagger.yml with your endpoints. See the [Swagger 2.0 spec](https://swagger.io/specification/v2/) for additional details on defining your swagger file.

- Run `make generate` to generate the supporting code

- Run `make build`, `make run`, or `make test` - This should fail with an error about having to implement the business logic.

- Implement aforementioned business logic so that code will build

## Using Yaml Aliases

If your swagger file uses yaml aliases, then you need to use python3 to substitute them in the file.
By default, the rule in the Makefile uses `wag-generate-mod`, replace this with `wag-yaml-aliases`.
Notice that this depends on the python library pyyaml, which can be installed with `pip3 install pyyaml`

## Deploying

```
ark start breakdown
```
