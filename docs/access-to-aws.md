# Access to AWS

Primary access to AWS is handled via the `aws-vault` command. This project uses `aws-vault` to store each user's
AWS credentials in a secure location like the macOS Keychain. All other access to AWS should happen inside of
docker. The best way to access those tools is by running docker-compose and using them inside the container:

```sh
make dev_up
make dev
```

Credentials are only needed when calling `make dev` and not when starting the development server with `make dev_up`.
If your AWS session token expires just drop out of the development environment and log back in. You will be prompted
for your MFA token and then dropped back into your development environment with proper credentials.

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
