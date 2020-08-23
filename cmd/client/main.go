package main

import "C"
import (
	"fhe"
	"flag"
	"fmt"
	"github.com/ldsec/lattigo/bfv"
	"github.com/spf13/viper"
)

const (
	FheServer = "127.0.0.1:50051"
	Token     = "123"
)

var gKeys bool
var wFiles bool
var evalFiles bool

func init() {
	flag.BoolVar(&gKeys, "g", false, "generates new keys")
	flag.BoolVar(&wFiles, "w", false, "writes new files")
	flag.BoolVar(&evalFiles, "e", false, "evaluates files")
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetDefault("FheServer", FheServer)
	viper.SetDefault("Token", Token)
	err := viper.ReadInConfig()
	if err != nil {

	}
	fmt.Println("starting client")
	flag.Parse()
	params := bfv.DefaultParams[bfv.PN13QP218]
	params.T = 0x3ee0001
	c := fhe.NewClient(viper.GetString("FheServer"), viper.GetString("Token"))
	if gKeys {
		fmt.Println("generating keys")
		c.GenKeys(params)
	}
	if wFiles {
		fmt.Println("uploading files")
		c.UploadFile(params)
	}
	if evalFiles {
		fmt.Println("evaluating files Grpc")
		c.EvalReq(params)
	}
}
