# REST API

A REST API template in Golang

## Prerequisites

- Install [Air](https://github.com/air-verse/air/tree/master) for live reload for application.
- Install [golang-migrate](https://github.com/golang-migrate/migrate/blob/v4.18.1/cmd/migrate/README.md) for running database schema migration.

## Available Scripts

Go to the root of the `nextgo-starter` and use the following commands in terminal:

| Name                                                     | Description                                                                             |
| -------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `npx nx run api:serve`                                   | Runs application via `go run` command                                                   |
| `npx nx run api:dev`                                     | Runs application with hot reloading                                                     |
| `npx nx run api:test`                                    | Runs unit test                                                                          |
| `npx nx run api:lint`                                    | Formats and lints application using the `go fmt`                                        |
| `npx nx run api:build`                                   | Builds application                                                                      |
| `npx nx run api:tidy`                                    | Runs go mod tidy to ensures that the go.mod file matches the source code in a Go module |
| `npx nx run api:create-migration "--args=filename=test"` | Generates a database migration `up` and `down` files with filename 'test'               |
| `npx nx run api:migrate-up`                              | Runs database migration                                                                 |
| `npx nx run api:migrate-down`                            | Reverses database migration                                                             |

All tasks are defined in `project.json`.
