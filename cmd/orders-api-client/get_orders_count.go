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
	// IssuerFlag is the orders uuid flag
	IssuerFlag string = "issuer"
	// RelativeTimeFlag is the relative time search flag
	RelativeTimeFlag string = "relative-time"
	// StartTimestampFlag is the start time search flag
	StartTimestampFlag string = "start-time"
	// EndTimestampFlag is the end time search flag
	EndTimestampFlag string = "end-time"
)

func initGetOrdersCountFlags(flag *pflag.FlagSet) {
	flag.String(IssuerFlag, "navy", "The Issuer of the orders")
	flag.Duration(RelativeTimeFlag, time.Hour*24, "The relative time to search backwards from when the command is invoked in UTC, set to '-1m' to disable")
	flag.String(StartTimestampFlag, "", "The Start time to search from, overrides relative-time search")
	flag.String(EndTimestampFlag, "", "The End time to search to, overrides relative-time search")

	flag.SortFlags = false
}

func checkGetOrdersCountConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	issuer := v.GetString(IssuerFlag)
	validIssuers := []string{"army", "navy", "air-force", "marine-corps", "coast-guard"}
	if issuer == "" {
		return fmt.Errorf("An issuer must be provided")
	} else if !stringInSlice(issuer, validIssuers) {
		return fmt.Errorf("Invalid issuer %q, must be one of %q", issuer, validIssuers)
	}

	startTimestamp := v.GetString(StartTimestampFlag)
	var startTime *time.Time
	if startTimestamp != "" {
		st, errStartTime := time.Parse(time.RFC3339, startTimestamp)
		startTime = &st
		if errStartTime != nil {
			return fmt.Errorf("Invalid value for %q flag: %w", StartTimestampFlag, errStartTime)
		}
	}
	endTimestamp := v.GetString(EndTimestampFlag)
	var endTime *time.Time
	if endTimestamp != "" {
		et, errEndTime := time.Parse(time.RFC3339, endTimestamp)
		endTime = &et
		if errEndTime != nil {
			return fmt.Errorf("Invalid value for %q flag: %w", EndTimestampFlag, errEndTime)
		}
	}
	if (startTime != nil && endTime != nil) && endTime.Before(*startTime) {
		return fmt.Errorf("Timestamp for %q flag  (%q)must be before %q flag  (%q)", StartTimestampFlag, *startTime, EndTimestampFlag, *endTime)
	}

	return nil
}

func getOrdersCount(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkGetOrdersCountConfig(v, args, logger)
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

	var params ordersOperations.GetOrdersCountByIssuerParams
	params.SetIssuer(v.GetString(IssuerFlag))

	// Get time search params
	startTimestamp := v.GetString(StartTimestampFlag)
	endTimestamp := v.GetString(EndTimestampFlag)

	if startTimestamp != "" {
		startTime, _ := time.Parse(time.RFC3339, startTimestamp)
		params.SetStartDateTime((*strfmt.DateTime)(&startTime))
	}
	if endTimestamp != "" {
		endTime, _ := time.Parse(time.RFC3339, endTimestamp)
		params.SetEndDateTime((*strfmt.DateTime)(&endTime))
	}
	if startTimestamp == "" && endTimestamp == "" {
		now := time.Now().UTC()
		relativeTime := v.GetDuration(RelativeTimeFlag)
		if relativeTime > time.Hour*0 {
			startTime := now.Add(relativeTime * -1)
			params.SetStartDateTime((*strfmt.DateTime)(&startTime))
			params.SetEndDateTime((*strfmt.DateTime)(&now))
		}
	}

	params.SetTimeout(time.Second * 30)
	resp, errGetOrdersCount := ordersGateway.Operations.GetOrdersCountByIssuer(&params)
	if errGetOrdersCount != nil {
		// If the response cannot be parsed as JSON you may see an error like
		// is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface
		// Likely this is because the API doesn't return JSON response for BadRequest OR
		// The response type is not being set to text
		logger.Fatal(errGetOrdersCount.Error())
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
