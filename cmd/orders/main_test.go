package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

type webServerSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger logger
}

func TestWebServerSuite(t *testing.T) {

	flag := pflag.CommandLine
	initServeFlags(flag)
	errParse := flag.Parse([]string{})
	if errParse != nil {
		log.Fatalf("Failed to parse flags due to %v", errParse)
	}

	v := viper.New()
	errBindPFlags := v.BindPFlags(flag)
	if errBindPFlags != nil {
		log.Fatalf("Failed to bind flags due to %v", errBindPFlags)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, errLoggingConfig := logging.Config(v.GetString(cli.DbEnvFlag), v.GetBool(cli.VerboseFlag))
	if errLoggingConfig != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", errLoggingConfig)
	}

	fields := make([]zap.Field, 0)
	if len(gitBranch) > 0 {
		fields = append(fields, zap.String("git_branch", gitBranch))
	}
	if len(gitCommit) > 0 {
		fields = append(fields, zap.String("git_commit", gitCommit))
	}
	logger = logger.With(fields...)
	zap.ReplaceGlobals(logger)

	ss := &webServerSuite{
		viper:  v,
		logger: logger,
	}

	suite.Run(t, ss)
}

// TestCheckServeConfigOrders is the acceptance test for the milmove webserver
// This will run all checks against the local environment and fail if something isn't configured
func (suite *webServerSuite) TestCheckServeConfigOrders() {
	if testEnv := os.Getenv("TEST_ACC_ENV"); len(testEnv) > 0 {
		filenameOrders := fmt.Sprintf("%s/config/env/%s.orders.env", os.Getenv("TEST_ACC_CWD"), testEnv)
		suite.logger.Info(fmt.Sprintf("Loading environment variables from file %s", filenameOrders))
		suite.applyContext(suite.patchContext(suite.loadContext(filenameOrders)))
	}

	suite.Nil(checkServeConfig(suite.viper, suite.logger))
}

// TestCheckServeConfigMigrate is the acceptance test for the milmove migration command
// This will run all checks against the local environment and fail if something isn't configured
func (suite *webServerSuite) TestCheckServeConfigMigrate() {
	if testEnv := os.Getenv("TEST_ACC_ENV"); len(testEnv) > 0 {
		filenameOrders := fmt.Sprintf("%s/config/env/%s.migrations.env", os.Getenv("TEST_ACC_CWD"), testEnv)
		suite.logger.Info(fmt.Sprintf("Loading environment variables from file %s", filenameOrders))
		suite.applyContext(suite.patchContext(suite.loadContext(filenameOrders)))
	}

	suite.Nil(checkMigrateConfig(suite.viper, suite.logger))
}

func (suite *webServerSuite) loadContext(variablesFile string) map[string]string {
	ctx := map[string]string{}
	if len(variablesFile) > 0 {
		if _, variablesFileStatErr := os.Stat(variablesFile); os.IsNotExist(variablesFileStatErr) {
			suite.logger.Fatal(fmt.Sprintf("File %q does not exist", variablesFile))
		}
		// Read contents of variables file into vars
		vars, err := ioutil.ReadFile(filepath.Clean(variablesFile))
		if err != nil {
			suite.logger.Fatal(fmt.Sprintf("error reading variables from file %s", variablesFile))
		}

		// Adds variables from file into context
		for _, x := range strings.Split(string(vars), "\n") {
			// If a line is empty or starts with #, then skip.
			if len(x) > 0 && x[0] != '#' {
				// Split each line on the first equals sign into []string{name, value}
				pair := strings.SplitAfterN(x, "=", 2)
				ctx[pair[0][0:len(pair[0])-1]] = pair[1]
			}
		}
	}
	return ctx
}

// patchContext updates specific variables based on value
func (suite *webServerSuite) patchContext(ctx map[string]string) map[string]string {
	for k, v := range ctx {
		if strings.HasPrefix(v, "/bin/") {
			ctx[k] = filepath.Join(os.Getenv("TEST_ACC_CWD"), v[1:])
		}
		// Overwrite the migration path to something on the local system
		if k == "MIGRATION_PATH" {
			ctx[k] = "file:///home/circleci/transcom/milmove_orders/migrations/orders/schema;file:///home/circleci/transcom/milmove_orders/migrations/orders/secure"
		}
	}
	return ctx
}

func (suite *webServerSuite) applyContext(ctx map[string]string) {
	for k, v := range ctx {
		suite.logger.Info("overriding " + k)
		suite.viper.Set(strings.Replace(strings.ToLower(k), "_", "-", -1), v)
	}
}
