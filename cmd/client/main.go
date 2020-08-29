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
)

const (
	FheServer = "127.0.0.1:50051"
	Token     = "123"
)

var gKeys bool
var wData bool
var wFiles bool
var evalFiles bool

func init() {
	flag.BoolVar(&gKeys, "g", false, "generates new keys")
	flag.BoolVar(&wData, "w", false, "writes new data")
	flag.BoolVar(&wFiles, "f", false, "writes new data read from csv files")
	flag.BoolVar(&evalFiles, "e", false, "evaluates files")
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
		fmt.Println("Evaluating files Grpc")
		c.EvalReq(params)
	}
	fmt.Println("Client finished")
}
