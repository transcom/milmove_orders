package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/pkg/cli"

	ordersOperations "github.com/transcom/milmove_orders/pkg/gen/ordersclient/operations"
	"github.com/transcom/milmove_orders/pkg/gen/ordersmessages"
)

var suffixes = []string{"JR", "SR", "II", "III", "IV", "V"}

type errInvalidCSVFile struct {
	Path string
}

func (e *errInvalidCSVFile) Error() string {
	return fmt.Sprintf("invalid CSV file path %q", e.Path)
}

const (
	// CSVFileFlag is the CSV file flag
	CSVFileFlag string = "csv-file"
)

func initPostRevisionsFlags(flag *pflag.FlagSet) {
	flag.String(CSVFileFlag, "", "The CSV File")
	flag.String(IssuerFlag, "navy", "The Issuer of the orders")

	flag.SortFlags = false
}

func checkPostRevisionsConfig(v *viper.Viper, args []string, logger *log.Logger) error {
	err := CheckRootConfig(v)
	if err != nil {
		logger.Fatal(err)
	}

	csvFile := v.GetString(CSVFileFlag)
	if len(csvFile) == 0 {
		return errors.New("missing csv file path, expected to be set")
	}
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		return fmt.Errorf("Expected %s to be a path in the filesystem: %w", csvFile, &errInvalidCSVFile{Path: csvFile})
	}

	// Currently the CSV parsing only support Navy orders
	issuer := v.GetString(IssuerFlag)
	validIssuers := []string{"navy"}
	if issuer == "" {
		return fmt.Errorf("An issuer must be provided")
	} else if !stringInSlice(issuer, validIssuers) {
		return fmt.Errorf("Invalid issuer %q, must be one of %q", issuer, validIssuers)
	}
	return nil
}

func postRevisions(cmd *cobra.Command, args []string) error {
	v := viper.New()

	//Create the logger
	//Remove the prefix and any datetime data
	logger := log.New(os.Stdout, "", log.LstdFlags)

	errParseFlags := ParseFlags(cmd, v, args)
	if errParseFlags != nil {
		return errParseFlags
	}

	// Check the config before talking to the CAC
	err := checkPostRevisionsConfig(v, args, logger)
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

	csvPath := v.GetString(CSVFileFlag)
	fileReader, err := os.Open(filepath.Clean(csvPath))
	if err != nil {
		logger.Fatal(err)
	}
	csvReader := csv.NewReader(fileReader)

	// First line contains the column headers; make a hash table that keys on the header with the column index as the value
	headers, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			log.Fatal("Empty file; no headers found")
		} else {
			log.Fatal(err)
		}
	}
	fields := make(map[string]int)
	for i := 0; i < len(headers); i++ {
		fields[headers[i]] = i
	}

	// every subsequent line can now be picked apart using this information
	for record, recordErr := csvReader.Read(); recordErr == nil; record, recordErr = csvReader.Read() {
		var rev ordersmessages.Revision
		rev.Member = new(ordersmessages.Member)
		rev.Member.Affiliation = ordersmessages.AffiliationNavy
		// The sailor's name, in the format LASTNAME,FIRSTNAME (optional MI) (optional suffix)
		fullname := record[fields["Service Member Name"]]
		names := strings.SplitN(fullname, ",", 2)
		rev.Member.FamilyName = names[0]
		names = strings.Fields(names[1])
		rev.Member.GivenName = names[0]
		if len(names) > 1 {
			if stringInSlice(names[len(names)-1], suffixes) {
				rev.Member.Suffix = &names[len(names)-1]
				if len(names) > 2 {
					middleName := strings.Join(names[1:len(names)-1], " ")
					rev.Member.MiddleName = &middleName
				}
			} else {
				middleName := strings.Join(names[1:], " ")
				rev.Member.MiddleName = &middleName
			}
		}

		daysStarting31Dec1899, _ := strconv.Atoi(record[fields["Order Create/Modification Date"]])
		dateIssued := time.Date(1899, time.December, 30+daysStarting31Dec1899, 0, 0, 0, 0, time.Local)
		fmtDateIssued := strfmt.DateTime(dateIssued)
		rev.DateIssued = &fmtDateIssued

		orderModNbr, orderModNbrErr := strconv.Atoi(record[fields["Order Modification Number"]])
		if orderModNbrErr != nil {
			orderModNbr = 0
		}
		obligModNbr, obligModNbrErr := strconv.Atoi(record[fields["Obligation Modification Number"]])
		if obligModNbrErr != nil {
			obligModNbr = 0
		}
		seqNum := int64(orderModNbr + obligModNbr)
		rev.SeqNum = &seqNum

		if record[fields["Obligation Status Code"]] == "D" {
			rev.Status = ordersmessages.StatusCanceled
		} else {
			rev.Status = ordersmessages.StatusAuthorized
		}
		rev.Member.Title = &record[fields["Rank Classification  Description"]]
		categorizedRank := paygradeToRank[record[fields["Paygrade"]]]
		rev.Member.Rank = categorizedRank.paygrade

		purpose := record[fields["CIC Purpose Information Code (OBLGTN)"]]
		if categorizedRank.officer {
			rev.OrdersType = officerOrdersTypes[purpose]
		} else {
			rev.OrdersType = enlistedOrdersTypes[purpose]
		}

		rev.LosingUnit = new(ordersmessages.Unit)
		if name := strings.TrimSpace(record[fields["Detach UIC Home Port"]]); len(name) > 0 {
			rev.LosingUnit.Name = &name
		}
		if uic := strings.TrimSpace(record[fields["Detach UIC"]]); len(uic) > 0 {
			fmtUIC := fmt.Sprintf("N%05s", uic)
			rev.LosingUnit.Uic = &fmtUIC
		}
		if city := strings.TrimSpace(record[fields["Detach UIC City Name"]]); len(city) > 0 {
			rev.LosingUnit.City = &city
		}
		if state := strings.TrimSpace(record[fields["Detach State Code"]]); len(state) > 0 {
			rev.LosingUnit.Locality = &state
		}
		if country := strings.TrimSpace(record[fields["Detach Country Code"]]); len(country) > 0 {
			rev.LosingUnit.Country = &country
		}

		daysStarting31Dec1899, daysError := strconv.Atoi(record[fields["Ultimate Estimated Arrival Date"]])
		if daysError == nil {
			estArrivalDate := time.Date(1899, time.December, 30+daysStarting31Dec1899, 0, 0, 0, 0, time.Local)
			rev.ReportNoLaterThan = new(strfmt.Date)
			*rev.ReportNoLaterThan = strfmt.Date(estArrivalDate)
		}

		rev.GainingUnit = new(ordersmessages.Unit)
		if name := strings.TrimSpace(record[fields["Ultimate UIC Home Port"]]); len(name) > 0 {
			rev.GainingUnit.Name = &name
		}
		if uic := strings.TrimSpace(record[fields["Ultimate UIC"]]); len(uic) > 0 {
			fmtUIC := fmt.Sprintf("N%05s", uic)
			rev.GainingUnit.Uic = &fmtUIC
		}
		if city := strings.TrimSpace(record[fields["Ultimate UIC City Name"]]); len(city) > 0 {
			rev.GainingUnit.City = &city
		}
		if state := strings.TrimSpace(record[fields["Ultimate State Code"]]); len(state) > 0 {
			rev.GainingUnit.Locality = &state
		}
		if country := strings.TrimSpace(record[fields["Ultimate Country Code"]]); len(country) > 0 {
			rev.GainingUnit.Country = &country
		}

		if record[fields["Entitlement Indicator"]] == "Y" {
			rev.NoCostMove = false
		} else {
			rev.NoCostMove = true
		}

		rev.HasDependents = new(bool)
		*rev.HasDependents = record[fields["Count of Dependents Participating in Move (STATIC)"]] != "0"

		if tdyEnRoute, tdyError := strconv.Atoi(record[fields["Count of Intermediate Stops (STATIC)"]]); tdyError == nil {
			rev.TdyEnRoute = tdyEnRoute > 0
		}

		rev.PcsAccounting = new(ordersmessages.Accounting)
		rev.PcsAccounting.Tac = &record[fields["TAC"]]

		if v.GetBool(cli.VerboseFlag) {
			bodyBuf := &bytes.Buffer{}
			encoder := json.NewEncoder(bodyBuf)
			encoder.SetIndent("", "  ")
			encoderErr := encoder.Encode(rev)
			if encoderErr != nil {
				log.Fatal(err)
			}

			fmt.Print(bodyBuf.String())
		}

		var params ordersOperations.PostRevisionParams
		params.SetMemberID(record[fields["Ssn (obligation)"]])
		params.SetOrdersNum(record[fields["Primary SDN"]])
		params.SetIssuer(string(ordersmessages.IssuerNavy))
		params.SetRevision(&rev)
		params.SetTimeout(time.Second * 30)
		resp, errGetOrders := ordersGateway.Operations.PostRevision(&params)
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
	}

	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return nil
}
