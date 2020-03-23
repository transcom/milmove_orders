# Scripts

This directory holds the scripts that are not compiled go code. For
compiled go code please look in the `bin/` directory of the project.

If you want to see if scripts are not listed in this file you can run
`find-scripts-missing-in-readme`.

## Dev Environment

These scripts are primarily used for managing the developers
environment.

| Script Name | Description |
| --- | --- |
| `check-aws-vault-version` | checks the aws-vault version required for the project |
| `check-hosts-file` | Script helps ensure that /etc/hosts has all the correct entries in it |
| `prereqs` | validate if all prerequisite programs have been installed |

## Operations Scripts

These scripts are used to operate the system.

| Script Name | Description |
| --- | --- |
| `deploy-orders` | Deploy the app |
| `deploy-orders-migrations` | Deploy the orders migrations |
| `health-tls-check` | Run health and TLS version checks. |

## Pre-commit Scripts

These scripts are used primarily to check our code before
committing.

| Script Name | Description |
| --- | --- |
| `pre-commit-go-mod` | modify `go.mod` and `go.sum` to match whats in the project |
| `pre-commit-swagger-validate` | Run swagger validate command inside docker container |

## CircleCI Scripts

These scripts are primarily used for CircleCI workflows.

| Script Name | Description |
| --- | --- |
| `check-deployed-commit` | checks that the deployed commit and given commit match. |
| `check-generated-code` | checks that the generated code has not changed |
| `circleci-announce-broken-branch` | announce that a branch is broken |
| `compare-deployed-commit` | checks that the given commit is ahead of the currently deployed commit |
| `ecr-describe-image-scan-findings` | Checks an uploaded image scan results |
| `ecs-deploy-service-container` | Updates the named service with the given name, image, and environment. |
| `ecs-run-orders-migrations-container` | Creates and runs a migration task using the given container definition. |
| `rds-snapshot-orders-db` | Creates a snapshot of the orders database for the given environment. |

## Development Scripts

These scripts are primarily used for developing the application and
application testing

| Script Name | Description |
| --- | --- |
| `update-docker-compose` | Update branch name before running docker-compose |

### Building

This subset of development scripts is used primarily for building the app.

| Script Name | Description |
| --- | --- |
| `gen-server` | generate swagger code from yaml files |

### Testing

This subset of development scripts is used for testing

| Script Name | Description |
| --- | --- |
| `run-e2e-test` | Runs integration tests with orders-api-client |
| `run-server-test` | Run golang server tests |

### Secure Migrations

This subset of development scripts is used in developing secure
migrations.

| Script Name | Description |
| --- | --- |
| `download-secure-migration` |  A script to download secure migrations from all environments |
| `generate-secure-migration` |  A script to help manage the creation of secure migrations |
| `upload-secure-migration` | A script to upload secure migrations to all environments |

### Database Scripts

These scripts are primarily used for working with the database

| Script Name | Description |
| --- | --- |
| `psql-dev` | Convenience script to drop into development postgres DB |
| `psql-deployed-migrations` | Convenience script to drop into deployed migrations postgres DB |
| `psql-schema` | Convenience script to dump the schema from the postgres DB |
| `psql-test` | Convenience script to drop into testing postgres DB |
| `psql-wrapper` | A wrapper around `psql` that sets correct values |
| `update-migrations-manifest` | Update manifest for migrations |
| `wait-for-db` |  waits for an available database connection, or until a timeout is reached |
| `wait-for-db-docker` |  waits for an available database connection, or until a timeout is reached using docker |

### CAC Scripts

These scripts are primarily used for working with a CAC and the Orders API

| Script Name | Description |
| --- | --- |
| `cac-prereqs` | Check the prereqs for CAC |

### Mutual TLS

These scripts are primarily for working with Mutual TLS certificates

| Script Name | Description |
| --- | --- |
| `mutual-tls-extract-fingerprint` | Get SHA 256 fingerprint of the public certificate from a cert file |
| `mutual-tls-extract-subject` | Get a sha256 hash of the certificate from a cert file|

### Amazon Console Scripts

These scripts are used for quickly opening up tools in the AWS Console

| Script Name | Description |
| --- | --- |

### Vulnerability Scanning

These scripts are used to do vulnerability scanning on our code

| Script Name | Description |
| --- | --- |
| `anti-virus` | Scan the source code for viruses |
