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
| `check-go-version` | checks the go version required for the project |
| `check-gopath` | checks the go path is correct for the project |
| `check-hosts-file` | Script helps ensure that /etc/hosts has all the correct entries in it |
| `prereqs` | validate if all prerequisite programs have been installed |

## AWS Scripts

These scripts are used for interacting with AWS or secrets in the AWS System Manager Parameter Store

| Script Name | Description |
| --- | --- |
| `aws` | Linked to aws-vault-wrapper. Runs the aws binary |
| `aws-vault-wrapper` | A wrapper to ensure AWS credentials are in the environment |
| `chamber` | Linked to aws-vault-wrapper. Runs chamber binary |

## Operations Scripts

These scripts are used to operate the system.

| Script Name | Description |
| --- | --- |

## Pre-commit Scripts

These scripts are used primarily to check our code before
committing.

| Script Name | Description |
| --- | --- |
| `pre-commit-go-mod` | modify `go.mod` and `go.sum` to match whats in the project |

## CircleCI Scripts

These scripts are primarily used for CircleCI workflows.

| Script Name | Description |
| --- | --- |
| `check-generated-code` | checks that the generated code has not changed |
| `circleci-announce-broken-branch` | announce that a branch is broken |
| `do-exclusively` | CircleCI's current recommendation for roughly serializing a subset of build commands for a given branch |

## Development Scripts

These scripts are primarily used for developing the application and
application testing

| Script Name | Description |
| --- | --- |

### Building

This subset of development scripts is used primarily for building the app.

| Script Name | Description |
| --- | --- |
| `gen-server` | generate swagger code from yaml files |

### Testing

This subset of development scripts is used for testing

| Script Name | Description |
| --- | --- |
| `run-server-test` | Run golang server tests |

### Secure Migrations

This subset of development scripts is used in developing secure
migrations.

| Script Name | Description |
| --- | --- |

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

### Mutual TLS

These scripts are primarily for working with Mutual TLS certificates

| Script Name | Description |
| --- | --- |

### Amazon Console Scripts

These scripts are used for quickly opening up tools in the AWS Console

| Script Name | Description |
| --- | --- |

### Vulnerability Scanning

These scripts are used to do vulnerability scanning on our code

| Script Name | Description |
| --- | --- |
