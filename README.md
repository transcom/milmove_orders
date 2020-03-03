# MilMove Electronic Orders

## License Information

Works created by U.S. Federal employees as part of their jobs typically are not eligible for copyright in the United
States. In places where the contributions of U.S. Federal employees are not eligible for copyright, this work is in
the public domain. In places where it is eligible for copyright, such as some foreign jurisdictions, the remainder of
this work is licensed under [the MIT License](https://opensource.org/licenses/MIT), the full text of which is included
in the [LICENSE.txt](./LICENSE.txt) file in this repository.

## Docker Compose

The primary development tool for building and testing is docker-compose. The advantage of developing with Docker
is the ability to develop and test inside the same environment as used by CircleCI and AWS. It also makes tool
distribution and versioning a lot simpler for distributed teams using different environments.

### Development

Development is handled entirely within docker using docker-compose. To enter the development environment run:

```sh
direnv allow
make dev
```

To reset your environment run

```sh
make dev_destroy dev
```

Every command in the `Makefile` is geared towards working inside the container. This means you can then run a command
like `make server_run` or `make db_dev_migrate` and it ought to work for you.

### Branch Builds

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

### Local Builds

Local testing of the application can be done using docker-compose. This method builds the docker images locally
and then orchestrates them to run together as a fully functional app. To run in this configuration use the command:

```sh
make docker_compose_local_up
```

And afterwards clean up with:

```sh
make docker_compose_local_down
```

## Access to AWS

Primary access to AWS is handled via the `aws-vault` command. This project uses `aws-vault` to store each user's
AWS credentials in a secure location like the macOS Keychain. All other access to AWS should happen inside of
docker. The best way to access those tools is by running docker-compose and using them inside the container:

```sh
make dev
```

Alternatively you can invoke docker directly to get a similar result. Replace `CLITOOL` in the following command with
a tool that requires access to AWS (like `aws` and `chamber`):

```sh
aws-vault exec "${AWS_PROFILE}" -- \
  docker run -it \
    -e AWS_ACCESS_KEY_ID \
    -e AWS_SECRET_ACCESS_KEY \
    -e AWS_SECURITY_TOKEN \
    -e AWS_SESSION_TOKEN \
    milmove/circleci-docker:milmove-orders CLITOOL
```
