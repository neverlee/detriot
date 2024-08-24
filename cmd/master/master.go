package master

import (
	"os"

	"github.com/neverlee/detriot/lrpc/log"
	"github.com/neverlee/detriot/lrpc/server"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Bind string `yaml:"bind"`
}

func LoadFileConfig(filepath string, conf interface{}) error {
	fdata, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(fdata, conf)
}

type Master struct {
}

func NewMaster() *Master {
	m := Master{}

	return &m
}

type TestRequest struct {
	Hello string `json:"hello"`
}

type TestResponse struct {
	Message string
}

func (ms *Master) HandleTest(req *TestRequest, rsp *TestResponse) error {
	rsp.Message = "rsp: " + req.Hello
	return nil
}

func Run(configPath string) error {
	log.Info("configPath:", configPath)

	srv := server.NewServer(":8000")
	ms := NewMaster()
	srv.Register(ms)
	err := srv.Run()
	return err
}

// 6478
