package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/neverlee/detriot/lrpc/log"
)

type RequestHead struct {
}

type RequestMessage struct {
	Head RequestHead     `json:"head"`
	Body json.RawMessage `json:"body"`
}

type ResponseHead struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseMessage struct {
	Head ResponseHead `json:"head"`
	Body any          `json:"body"`
}

type Handler struct {
	Name         string
	method       reflect.Value
	mainPointer  reflect.Value
	requestType  reflect.Type
	responseType reflect.Type
}

type Server struct {
	httpSrv  http.Server
	handlers map[string]*Handler
}

func NewServer(addr string) *Server {
	s := Server{
		httpSrv: http.Server{
			Addr: addr,
		},
		handlers: make(map[string]*Handler),
	}
	return &s
}

func (srv *Server) Register(vh interface{}) {
	rt := reflect.TypeOf(vh)
	rv := reflect.ValueOf(vh)

	const handlePrefix = "Handle"
	vhTypeName := rt.Elem().Name()
	keyPrefix := "/qrpc/" + strings.ToLower(vhTypeName) + "/"
	for i := 0; i < rt.NumMethod(); i++ {
		rtmethod := rt.Method(i)
		rttmethod := rtmethod.Type

		if !strings.HasPrefix(rtmethod.Name, handlePrefix) || rttmethod.NumIn() != 3 || rttmethod.NumOut() != 1 {
			continue
		}

		hkey := keyPrefix + strings.ToLower(rtmethod.Name[len(handlePrefix):])

		handler := Handler{
			Name:         fmt.Sprintf("%s.%s", vhTypeName, rtmethod.Name),
			method:       rtmethod.Func,
			mainPointer:  rv,
			requestType:  rttmethod.In(1).Elem(),
			responseType: rttmethod.In(2).Elem(),
		}

		srv.handlers[hkey] = &handler
		log.Infof("register handle %s to %s\n", handler.Name, hkey)
	}
}

func (srv *Server) Handle(rspw http.ResponseWriter, req *http.Request) (rerr error) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			rerr = err
		}
	}()

	handler, ok := srv.handlers[strings.ToLower(req.URL.Path)]
	if !ok {
		return errors.New("no such handler")
	}

	dec := json.NewDecoder(req.Body)
	reqv := reflect.New(handler.requestType)
	reqvi := reqv.Interface()
	if err := dec.Decode(reqvi); err != nil {
		return err
	}

	rspv := reflect.New(handler.responseType)
	rspvi := rspv.Interface()

	log.Info("request handle call: ", req.URL.Path, handler.Name)
	rvresults := handler.method.Call([]reflect.Value{handler.mainPointer, reqv, rspv})

	if !rvresults[0].IsNil() {
		if cerr := rvresults[0].Interface().(error); cerr != nil {
			rspw.Header().Set("qrpc_code", "-1") // call error
			rspw.Header().Set("qrpc_message", cerr.Error())
			return nil
		}
	} else {
		enc := json.NewEncoder(rspw)
		if err := enc.Encode(rspvi); err != nil {
			return err
		}
	}

	return nil
}

func (srv *Server) ServeHTTP(rspw http.ResponseWriter, req *http.Request) {
	err := srv.Handle(rspw, req)
	if err != nil {
		rspw.Header().Set("qrpc_code", "-2") // call error
		rspw.Header().Set("qrpc_message", err.Error())
		log.Warn("request process err", req.URL.Path, err)
	}
}

func (srv *Server) Run() error {
	srv.httpSrv.Handler = srv
	return srv.httpSrv.ListenAndServe()
}
