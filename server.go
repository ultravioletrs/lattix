package fhe

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/ldsec/lattigo/bfv"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type In struct {
	Params bfv.Parameters
	Res    bfv.Ciphertext
}

type Server struct {
	UnimplementedFhesrvServer
	filesDir       string
	filesExtension string
	token          string
}

func NewServer(filesDir, filesExtension, token string) *Server {
	return &Server{
		UnimplementedFhesrvServer: UnimplementedFhesrvServer{},
		filesDir:                  filesDir,
		filesExtension:            filesExtension,
		token:                     token,
	}
}

func (s *Server) EvalFiles(ctx context.Context, in *EvalRequest) (*EvalReply, error) {
	log.Printf("Received from: %d to: %d", in.Fromtimestamp, in.Totimestamp)
	md, _ := metadata.FromIncomingContext(ctx)
	token := md.Get("token")[0]
	if token != s.token {
		fmt.Println(fmt.Sprintf("Got wrong token %s, expected %s", token, s.token))
		return nil, errors.New("wrong token")

	}
	ib := bytes.NewBuffer(in.Request)
	ob := bytes.Buffer{}
	decoder := gob.NewDecoder(ib)
	encoder := gob.NewEncoder(&ob)
	inr := In{}

	time.Sleep(time.Second)
	err := decoder.Decode(&inr)
	if err != nil {
		fmt.Println("Fatal error while handling decode", err.Error())
		os.Exit(1)
	}
	s.eval(&inr.Params, &inr.Res, in.Fromtimestamp, in.Totimestamp)
	encoder.Encode(&inr.Res)
	return &EvalReply{
		Message:  fmt.Sprintf("Result for period from: %d to: %d", in.Fromtimestamp, in.Totimestamp),
		Response: ob.Bytes(),
	}, nil
}

func (s *Server) UploadFile(ctx context.Context, in *UploadRequest) (*UploadReply, error) {
	log.Printf("Received: upload request")
	err := ioutil.WriteFile(fmt.Sprintf("%s%d.%s", s.filesDir, time.Now().UnixNano(), s.filesExtension), in.File, 0644)
	if err != nil {
		return &UploadReply{
			Message: "NOT OK",
		}, err
	}
	return &UploadReply{
		Message: "OK",
	}, nil
}

func (s *Server) eval(params *bfv.Parameters, res *bfv.Ciphertext, fromTimestamp, toTimestamp int64) *bfv.Ciphertext {
	var files []int64
	evaluator := bfv.NewEvaluator(params)
	filepath.Walk(s.filesDir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		ext := path.Ext(p)
		base := path.Base(p)
		if ext != "."+s.filesExtension {
			return nil
		}
		i, err := strconv.ParseInt(strings.Split(base, ".")[0], 10, 64)
		if err != nil {
			return nil
		}
		if i >= fromTimestamp && i <= toTimestamp {
			files = append(files, i)
		}
		return nil
	})
	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("%s%d.%s", s.filesDir, file, s.filesExtension))
		if err != nil {
			fmt.Println("Fatal error ", err.Error())
			os.Exit(1)
		}
		decoder := gob.NewDecoder(f)
		var driver bfv.Ciphertext
		err = decoder.Decode(&driver)
		if err != nil {
			fmt.Println("Fatal error ", err.Error())
			os.Exit(1)
		}
		evaluator.Add(res, driver, res)
	}
	return res
}
