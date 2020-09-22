// +build ignore

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// TODO(bassosimone): most of the code in this file should be moved to
// github.com/ooni/probe-engine once we're always using the new API. When
// this happens, this code will be part of the probeservices package.

// The following errors are returned by this API.
var (
	ErrEmptyReportID = errors.New("empty report ID")
)

// Config contains the configuration for fetching a measurement.
type Config struct {
	// ReportID is the mandatory report ID.
	ReportID string

	// Full indicates whether we also want the full measurement body.
	Full bool

	// Input is the optional input.
	Input string

	// Debugf is a function called to emit debug messages.
	Debugf func(format string, v ...interface{})

	// Client is the optional HTTP client.
	Client *http.Client

	// BaseURL is the optional base URL.
	BaseURL string
}

// MeasurementMeta contains measurement metadata.
type MeasurementMeta struct {
	// Fields returned by the API server whenever we are
	// calling /api/v1/measurement_meta.
	Anomaly              bool      `json:"anomaly"`
	CategoryCode         string    `json:"category_code"`
	Confirmed            bool      `json:"confirmed"`
	Failure              bool      `json:"failure"`
	Input                *string   `json:"input"`
	MeasurementStartTime time.Time `json:"measurement_start_time"`
	ProbeASN             int64     `json:"probe_asn"`
	ProbeCC              string    `json:"probe_cc"`
	ReportID             string    `json:"report_id"`
	Scores               string    `json:"scores"`
	TestName             string    `json:"test_name"`
	TestStartTime        time.Time `json:"test_start_time"`

	// This field is only included if the user has specified
	// the config.Full option, otherwise it's empty.
	RawMeasurement string `json:"raw_measurement"`

	// This field contains the body that we received from the
	// API server and it's here to help debugging.
	RawBody []byte `json:"-"`
}

// GetMeasurementMeta gets measurement metadata.
func GetMeasurementMeta(ctx context.Context, config Config) (MeasurementMeta, error) {
	if config.ReportID == "" {
		return MeasurementMeta{}, ErrEmptyReportID
	}
	if config.Debugf == nil {
		config.Debugf = log.Printf
	}
	if config.Client == nil {
		config.Client = http.DefaultClient
	}
	if config.BaseURL == "" {
		config.BaseURL = "https://ams-pg.ooni.org"
	}
	URL, err := url.Parse(config.BaseURL)
	if err != nil {
		return MeasurementMeta{}, err
	}
	URL.Path = "/api/v1/measurement_meta"
	query := url.Values{}
	query.Add("report_id", config.ReportID)
	if config.Input != "" {
		query.Add("input", config.Input)
	}
	if config.Full {
		query.Add("full", "true")
	}
	URL.RawQuery = query.Encode()
	config.Debugf("> GET %s", URL.String())
	resp, err := config.Client.Get(URL.String())
	if err != nil {
		return MeasurementMeta{}, err
	}
	config.Debugf("< %d", resp.StatusCode)
	defer resp.Body.Close()
	// TODO(bassosimone): this would be nice to have in most
	// github.com/ooni/probe-engine/probeservices APIs.
	reader := io.LimitReader(resp.Body, 1<<25)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return MeasurementMeta{}, err
	}
	var mmeta MeasurementMeta
	err = json.Unmarshal(body, &mmeta)
	mmeta.RawBody = body // helps debugging
	if err != nil {
		return mmeta, err
	}
	return mmeta, nil
}

func fatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	reportid := flag.String("report-id", "", "Report ID of the measurement")
	input := flag.String("input", "", "Input of the measurement")
	full := flag.Bool("full", false, "Also include the measurement body")
	flag.Parse()
	mmeta, err := GetMeasurementMeta(context.Background(), Config{
		ReportID: *reportid,
		Input:    *input,
		Full:     *full,
	})
	fatalOnError(err)
	data, err := json.Marshal(mmeta)
	fatalOnError(err)
	fmt.Printf("%s\n", data)
}
