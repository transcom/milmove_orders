package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	ordersOperations "github.com/transcom/milmove_orders/pkg/gen/ordersclient/operations"
)

const (
	// OrdersUUIDFlag is the orders uuid flag
	OrdersUUIDFlag string = "orders-uuid"
)

func initGetOrdersFlags(flag *pflag.FlagSet) {
	flag.String(OrdersUUIDFlag, "", "The UUID of the ")

	flag.SortFlags = false
}

func checkGetOrdersConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	ordersUUID := v.GetString(OrdersUUIDFlag)
	if ordersUUID == "" {
		return fmt.Errorf("An orders uuid must be provided")
	} else if !strfmt.IsUUID(ordersUUID) {
		return fmt.Errorf("Unable to parse uuid %q", ordersUUID)
	}

	return nil
}

func getOrders(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkGetOrdersConfig(v, args, logger)
	if err != nil {
		logger.Fatal(err)
	}

	ordersGateway, cacStore, errCreateClient := CreateClient(v)
	if errCreateClient != nil {
		return errCreateClient
	}

	// Defer closing the store until after the API call has completed
	if cacStore != nil {
		defer func() { _ = cacStore.Close() }()
	}

	var params ordersOperations.GetOrdersParams
	params.SetUUID(strfmt.UUID(v.GetString(OrdersUUIDFlag)))
	params.SetTimeout(time.Second * 30)
	resp, errGetOrders := ordersGateway.Operations.GetOrders(&params)
	if errGetOrders != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		logger.Fatal(errGetOrders.Error())
	}

	payload := resp.GetPayload()
	if payload != nil {
		payload, errJSONMarshall := json.Marshal(payload)
		if errJSONMarshall != nil {
			logger.Fatal(errJSONMarshall)
		}
		fmt.Println(string(payload))
	} else {
		logger.Fatal(resp.Error())
	}

	return nil
}
