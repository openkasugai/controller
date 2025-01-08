package main

import (
	"context"
	"flag"
	"fmt"
	"calcapp/api/operator"
	"net"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./api/operator/operator.proto

func init() {
	config := zap.NewDevelopmentConfig()
	l, _ := config.Build()
	zap.ReplaceGlobals(l)
}

var usage = `Usage: %s
(server mode) -server -port <port> -operator <operator>
(client mode) -host <host> -port <port> <1.0 ...values>

`

const (
	ConfigOperatorPlus           = "plus"
	ConfigOperatorMinus          = "minus"
	ConfigOperatorMultiply       = "multiply"
	ConfigOperatorDivide         = "divide"
	ConfigOperatorAverageResults = "average_results"
	ConfigOperatorReceiver       = "receiver"
)

func main() {
	log := zap.L()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}

	operators := []string{
		ConfigOperatorPlus,
		ConfigOperatorMinus,
		ConfigOperatorMultiply,
		ConfigOperatorDivide,
		ConfigOperatorAverageResults,
		ConfigOperatorReceiver,
	}

	isServerMonde := false
	flag.BoolVar(&isServerMonde, "server", false, "run server mode")
	port := 8080
	flag.IntVar(&port, "port", port, "port number")
	host := ""
	flag.StringVar(&host, "host", host, "host address")
	operator := ""
	flag.StringVar(&operator, "operator", operator, "operator <"+strings.Join(operators, ",")+">")
	values := []float64{}

	flag.Parse()

	if isServerMonde {
		if !slices.Contains(operators, operator) {
			log.Error("operator not match")
			flag.Usage()
			os.Exit(1)
		}
		runServerMode(port, operator)
	} else {
		if flag.NArg() == 0 {
			log.Error("zero values")
			flag.Usage()
			os.Exit(1)
		}
		for _, v := range flag.Args() {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.Error("values:", zap.Error(err))
				flag.Usage()
				os.Exit(1)
			}
			values = append(values, f)
		}

		if host == "" {
			cfg, found, err := loadConfig()
			if err != nil {
				log.Error("load config file", zap.Bool("found", found), zap.Error(err))
			}
			log.Info("load config", zap.Bool("found", found), zap.Any("config", cfg))
			if found {
				host = cfg.Next.Host
				port = cfg.Next.Port
			}
		}

		if host == "" {
			log.Error("error", zap.String("host", host), zap.Int("port", port))
			flag.Usage()
			os.Exit(1)
		}
		runCientMode(host, port, values)
	}
}

func registerNotifyConfigChange(ctx context.Context) <-chan struct{} {
	notify := make(chan struct{}, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		defer signal.Stop(c)
		defer close(notify)

		log := zap.L()
		log.Info("start notify loader")
		defer log.Info("end notify loader")

		for {
			select {
			case _, ok := <-c:
				if !ok {
					return
				}
				notify <- struct{}{}
			case <-ctx.Done():
				return
			}
		}
	}()

	return notify
}

func registerNotifiedConfigLoader(ctx context.Context, notify <-chan struct{}, svr *OperatorServer) {
	go func() {
		log := zap.L()
		log.Info("start config loader")
		defer log.Info("end config loader")

		for {
			select {
			case _, ok := <-notify:
				log.Info("loading config")
				if !ok {
					log.Info("notify closed")
					return
				}
				cfg, found, err := loadConfig()
				if err != nil {
					log.Error("load config file", zap.Bool("found", found), zap.Error(err))
					continue
				}
				log.Info("load config", zap.Bool("found", found), zap.Any("config", cfg))
				if found {
					svr.SetNextAddress(fmt.Sprintf("%s:%d", cfg.Next.Host, cfg.Next.Port))
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

type Config struct {
	Next *struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"next"`
}

func runServerMode(port int, operate string) {
	log := zap.L()

	var opeServer *OperatorServer = nil
	switch operate {
	case ConfigOperatorPlus:
		opeServer = &OperatorServer{operator: &PlusOperator{}}
	case ConfigOperatorMinus:
		opeServer = &OperatorServer{operator: &MinusOperator{}}
	case ConfigOperatorMultiply:
		opeServer = &OperatorServer{operator: &MultiplyOperator{}}
	case ConfigOperatorDivide:
		opeServer = &OperatorServer{operator: &DivideOperator{}}
	case ConfigOperatorAverageResults:
		opeServer = &OperatorServer{operator: &AverageResultsOperator{}}
	case ConfigOperatorReceiver:
		opeServer = &OperatorServer{operator: &ReceiverOperator{}}
	}

	ctx, cancel := context.WithCancel(context.Background())

	notify := registerNotifyConfigChange(ctx)

	registerNotifiedConfigLoader(ctx, notify, opeServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("server listener faild", zap.Error(err))
	}

	server := grpc.NewServer()

	operator.RegisterOperatorServiceServer(server, opeServer)

	go func() {
		defer cancel()

		log.Info("start gRPC server", zap.Int("port", port))
		err := server.Serve(lis)
		if err != nil {
			log.Error("start gRPC faild", zap.Error(err))
		}
		log.Info("end gRPC server", zap.Int("port", port))
	}()

	<-ctx.Done()

	log.Info("stopping gRPC server...", zap.Int("port", port))
	if context.Cause(ctx) == context.Canceled {
		server.GracefulStop()
	}
}

func loadConfig() (*Config, bool, error) {
	const ConfigFilePath = "/config/config.yaml"

	buf, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		} else {
			return nil, true, err
		}
	}

	config := &Config{}
	err = yaml.Unmarshal(buf, config)

	if err != nil {
		return nil, true, err
	}
	return config, true, nil
}

func runCientMode(host string, port int, values []float64) {
	log := zap.L()

	ctx := context.Background()

	req := &operator.OperateRequest{
		Inputs: values,
	}

	address := fmt.Sprintf("%s:%d", host, port)

	log.Info("send message", zap.String("address", address), zap.Any("request", req))
	nextRes, err := sendMessage(ctx, address, req)
	if err != nil {
		log.Error("send message faild", zap.Error(err))
	}
	log.Info("recive message", zap.String("address", address), zap.Any("response", nextRes))
}

func sendMessage(ctx context.Context, address string, req *operator.OperateRequest) (*operator.OperateResponse, error) {
	log := zap.L()

	log.Info("connecting", zap.String("address", address))
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("did not connect", zap.String("address", address), zap.Error(err))
	}
	defer conn.Close()

	c := operator.NewOperatorServiceClient(conn)
	res, err := c.Operate(ctx, req)
	if err != nil {
		log.Error("did not rpc", zap.String("address", address), zap.Error(err))
	}
	return res, err
}

type Operator interface {
	OperatorType() string
	Operate(*operator.OperateRequest) float64
	HasNextSendFeature() bool
}
type OperatorServer struct {
	operator.UnimplementedOperatorServiceServer
	operator    Operator
	nextAddress string
	m           sync.RWMutex
}

func (svr *OperatorServer) Operate(ctx context.Context, req *operator.OperateRequest) (*operator.OperateResponse, error) {
	log := zap.L()

	log.Info("start call Do", zap.Any("req", req))

	if svr.operator.HasNextSendFeature() {
		nextReq := &operator.OperateRequest{}
		proto.Merge(nextReq, req)

		nextReq.Results = append(nextReq.Results, &operator.OperateResult{
			Operator: svr.operator.OperatorType(),
			Value:    svr.operator.Operate(req),
		})

		svr.m.RLock()
		defer svr.m.RUnlock()

		address := svr.GetNextAddress()
		nextRes, err := sendMessage(ctx, address, nextReq)
		if err != nil {
			log.Error("send message faild", zap.Error(err))
		}
		log.Info("send message", zap.Any("res", nextRes))

		log.Info("end call Do", zap.Any("req", req))
	}
	return &operator.OperateResponse{
		Status: operator.OperateResponseStatus_OK,
	}, nil
}

func (svr *OperatorServer) GetNextAddress() string {
	svr.m.RLock()
	defer svr.m.RUnlock()
	return svr.nextAddress
}
func (svr *OperatorServer) SetNextAddress(nextAddress string) {
	svr.m.Lock()
	defer svr.m.Unlock()
	svr.nextAddress = nextAddress
}

type PlusOperator struct {
}

func (*PlusOperator) HasNextSendFeature() bool { return true }
func (*PlusOperator) OperatorType() string {
	return "plus"
}
func (*PlusOperator) Operate(req *operator.OperateRequest) float64 {
	var res float64
	if 1 < len(req.Inputs) {
		res = req.Inputs[0]
		for _, v := range req.Inputs[1:] {
			res = res + v
		}
	}
	return res
}

type MinusOperator struct {
}

func (*MinusOperator) HasNextSendFeature() bool { return true }

func (*MinusOperator) OperatorType() string {
	return "minus"
}
func (*MinusOperator) Operate(req *operator.OperateRequest) float64 {
	var res float64
	if 1 < len(req.Inputs) {
		res = req.Inputs[0]
		for _, v := range req.Inputs[1:] {
			res = res - v
		}
	}
	return res
}

type MultiplyOperator struct {
}

func (*MultiplyOperator) HasNextSendFeature() bool { return true }

func (*MultiplyOperator) OperatorType() string {
	return "multiply"
}
func (*MultiplyOperator) Operate(req *operator.OperateRequest) float64 {
	var res float64
	if 1 < len(req.Inputs) {
		res = req.Inputs[0]
		for _, v := range req.Inputs[1:] {
			res = res * v
		}
	}
	return res
}

type DivideOperator struct {
}

func (*DivideOperator) HasNextSendFeature() bool { return true }

func (*DivideOperator) OperatorType() string {
	return "divide"
}
func (*DivideOperator) Operate(req *operator.OperateRequest) float64 {
	var res float64
	if 1 < len(req.Inputs) {
		res = req.Inputs[0]
		for _, v := range req.Inputs[1:] {
			res = res / v
		}
	}
	return res
}

type AverageResultsOperator struct {
}

func (*AverageResultsOperator) HasNextSendFeature() bool { return true }

func (*AverageResultsOperator) OperatorType() string {
	return "average"
}
func (*AverageResultsOperator) Operate(req *operator.OperateRequest) float64 {
	s := 0.0
	for _, v := range req.Results {
		s += v.Value
	}
	return s / float64(len(req.Results))
}

type ReceiverOperator struct {
}

func (*ReceiverOperator) HasNextSendFeature() bool { return false }

func (*ReceiverOperator) OperatorType() string {
	return "receiver"
}
func (*ReceiverOperator) Operate(req *operator.OperateRequest) float64 {
	return 0.0
}
