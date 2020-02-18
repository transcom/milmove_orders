module github.com/transcom/milmove_orders

go 1.13

require (
	github.com/aws/aws-sdk-go v1.29.4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/loads v0.19.4
	github.com/go-openapi/runtime v0.19.10
	github.com/go-openapi/spec v0.19.5
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.7
	github.com/go-openapi/validate v0.19.6
	github.com/gobuffalo/pop v4.13.1+incompatible
	github.com/gobuffalo/validate v2.0.4+incompatible
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gorilla/csrf v1.6.2
	github.com/jessevdk/go-flags v1.4.0
	github.com/lib/pq v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/rickar/cal v1.0.3
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
	github.com/transcom/mymove v0.0.0-20200217234508-33d518e6bb13
	go.uber.org/zap v1.13.0
	goji.io v2.0.2+incompatible
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
)

// Update to ignore compiler warnings on macOS catalina
// https://github.com/keybase/go-keychain/pull/55
// https://github.com/99designs/aws-vault/pull/427
replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
