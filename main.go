// +build !js

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/empijei/identigen/identities"
)

var minage = flag.Int("minage", 25, "The minimum age for random people generation. Must be positive and less than maxage.")
var maxage = flag.Int("maxage", 55, "The maximum age for random people generation. Must be positive, less or equal than 200 and more than minage.")
var number = flag.Int("number", 1, "The amount of random people to generate. Must be positive.")
var dtFmt = flag.String("dt_fmt", "eu", "The format of the dates. Supports: 'eu','us','ja'")
var format = flag.String("format", "human", "The comma separated list of formats for the output. Supports: 'json', 'csv', 'xml', 'human'.")
var fields = flag.String("fields", "all", "The comma separated case-sensitive list of fields to print. Use 'all' to print all of them. Supported fields are: "+strings.Join(identities.AllFields, ","))
var country = flag.String("country", "IT", "The two characters ISO 3166 code for the identity's nationality.")

func main() {
	flag.Parse()
	args := make(map[string]interface{})
	args["dt_fmt"] = *dtFmt
	args["minage"] = *minage
	args["maxage"] = *maxage
	args["number"] = *number
	args["format"] = *format
	args["fields"] = *fields
	args["country"] = *country

	err := identities.MainModule(args, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
}
