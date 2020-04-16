module github.com/transcom/milmove_orders

go 1.14

require (
	github.com/aws/aws-sdk-go v1.30.7
	github.com/codegangsta/gin v0.0.0-20171026143024-cafe2ce98974
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/spec v0.19.7
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.8
	github.com/go-openapi/validate v0.19.8
	github.com/gobuffalo/envy v1.9.0
	github.com/gobuffalo/pop v4.13.1+incompatible
	github.com/gobuffalo/validate v2.0.4+incompatible
	github.com/gofrs/flock v0.7.1
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/imdario/mergo v0.3.9
	github.com/jessevdk/go-flags v1.4.0
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.3.0
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.7
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/objx v0.2.0
	github.com/stretchr/testify v1.5.1
	github.com/transcom/mymove v0.0.0-20200413222450-b2771f9c20a5
	go.uber.org/zap v1.14.1
	goji.io v2.0.2+incompatible
	golang.org/x/mod v0.2.0 // indirect
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae // indirect
	golang.org/x/tools v0.0.0-20200224181240-023911ca70b2 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543 // indirect
	gopkg.in/ini.v1 v1.52.0 // indirect
	pault.ag/go/pksigner v1.0.2
)

// transcom/sqlx v1.2.1 is just jmoiron's 1.2.0 with custom driver fixes
// This is a temporary solution till https://github.com/jmoiron/sqlx/pull/560
// is merged or a better solution is completed as mentioned in
// https://github.com/jmoiron/sqlx/pull/520
replace github.com/jmoiron/sqlx v1.2.0 => github.com/transcom/sqlx v1.2.1

// Update to ignore compiler warnings on macOS catalina
// https://github.com/keybase/go-keychain/pull/55
// https://github.com/99designs/aws-vault/pull/427
replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
