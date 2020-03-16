# Developing with docker-compose

The primary development tool for building and testing is docker-compose. The advantage of developing with Docker
is the ability to develop and test inside the same environment as used by CircleCI and AWS. It also makes tool
distribution and versioning a lot simpler for distributed teams using different environments.

## Development

Development is handled entirely within docker using docker-compose. To start the development environment run:

```sh
direnv allow
make dev_up
```

Now the docker-compose file has been brought up, running both a Postgres container and the app server (which is
equivalent to `make server_run`).

To enter the development environment with a Bash prompt run:

```sh
make dev
```

To reset your environment run

```sh
make dev_reset
```

Every command in the `Makefile` is geared towards working inside the container. This means you can then run a command
like `make server_run` or `make db_dev_migrate` and it ought to work for you.

To destroy your development environment run:

```sh
make dev_down
```

**Note on AWS Credentials:** credentials are passed in when you run `make dev`. When credentials expire you simply
exit the docker container with `Ctrl-C` and then `make dev` again. You will be prompted for your MFA token and then
logged into your development environment with updated AWS credentials. It isn't necessary to do this for the running
server since there is no case where the development server needs to talk to AWS. That means you do not need to do
`make dev_reset` when the creds expire.

## Branch Builds

Requirements:

- AWS Account with access to AWS ECR repos `orders` and `orders-migrations`

Branch testing of the application can be done using docker-compose. This method uses docker images built in CircleCI
and then orchestrates them to run together as a fully functional app. To run in this configuration use the command:

```sh
make docker_compose_local_up
```

And afterwards clean up with:

```sh
make docker_compose_local_down
```

## Local Builds

Local testing of the application can be done using docker-compose. This method builds the docker images locally
and then orchestrates them to run together as a fully functional app. To run in this configuration use the command:

```sh
make docker_compose_local_up
```

And afterwards clean up with:

```sh
make docker_compose_local_down
```
