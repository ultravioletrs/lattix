package main

import "C"
import (
	"fhe"
	"flag"
	"fmt"
	"github.com/ldsec/lattigo/bfv"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"strings"
)

const (
	FheServer = "127.0.0.1:50051"
	Token     = "123"
)

var gKeys bool
var wData bool
var wFiles bool
var evalFiles bool
var fromTimestamp int64
var toTimestamp int64
var returnElements string

func init() {
	flag.BoolVar(&gKeys, "g", false, "generates new keys")
	flag.BoolVar(&wData, "w", false, "writes new data")
	flag.BoolVar(&wFiles, "c", false, "writes new data read from csv files")
	flag.BoolVar(&evalFiles, "e", false, "evaluates files")
	flag.Int64Var(&fromTimestamp, "f", 0, "set from timestamp")
	flag.Int64Var(&toTimestamp, "t", 0, "set to timestamp")
	flag.StringVar(&returnElements, "r", "", "set comma separated elements to be returned")
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("client")
	viper.SetDefault("fhe_server", FheServer)
	viper.SetDefault("token", Token)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file: %s", err)
	}
	fmt.Println("Starting client")
	flag.Parse()
	params := bfv.DefaultParams[bfv.PN13QP218]
	params.T = 0x3ee0001
	c := fhe.NewClient(viper.GetString("fhe_server"), viper.GetString("token"))
	if gKeys {
		fmt.Println("Generating keys")
		c.GenKeys(params)
	}
	if wData {
		vals := flag.Args()
		var data []uint64
		for _, val := range vals {
			u, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				log.Fatalf("Invalid parameter: %s", err)
			}
			data = append(data, u)
		}
		fmt.Println("Writing data ", data)
		err := c.Write(params, data...)
		if err != nil {
			log.Fatalf("Error while writing data: %s", err)
		}
	}
	if wFiles {
		files := flag.Args()
		for _, file := range files {
			fmt.Println("Writing data from file ", file)
			if err := c.WriteFromFile(params, file); err != nil {
				log.Fatalf("Error while writing from file: %s", err)
			}
		}
	}
	if evalFiles {
		fmt.Println("Evaluating files on server")
		res, err := c.EvalReq(params, fromTimestamp, toTimestamp)
		if err != nil {
			log.Fatalf("An error occured while evaluate request: %s", err)
		}
		positions := strings.Split(returnElements, ",")
		for _, position := range positions {
			pos, err := strconv.ParseInt(position, 10, 64)
			if err != nil {
				log.Fatalf("Invalid parameter: %s", err)
			}
			fmt.Printf("Result at position %d: %d\n", pos, res[pos])
		}
		if len(positions) == 0 {
			fmt.Println(res)
		}

	}
}
