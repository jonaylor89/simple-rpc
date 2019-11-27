
package client

import (
  "errors"
  "net"
  "reflect"

  "github.com/jonaylor89/simple-rpc/data"
  "github.com/jonaylor89/simple-rpc/transport"
)

// Client : client struct
type Client struct {
  conn net.Conn
}

// NewClient : creates a new client
func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

// CallRPC : call remote function
func (c *Client) CallRPC(rpcName string, fPtr interface{}) {
  container := reflect.ValueOf(fPtr).Elem()
  f := func(req []reflect.Value) []reflect.Value {
    cReqTransport := transport.NewTransport(c.conn) 
    errorHandler := func(err error) []reflect.Value {
      outArgs := make([]reflect.Value, container.Type().NumOut()) 
      for i := 0; i < len(outArgs)-1; i++ {
        outArgs[i] = reflect.Zero(container.Type().Out(i)) 
      }

      outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()

      return outArgs
    }

    // Process input parameters
    inArgs := make([]interface{}, 0, len(req))
    for _, arg := range req {
      inArgs = append(inArgs, arg.Interface()) 
    }


    reqRPC := data.RPCdata{Name: rpcName, Args: inArgs}

    b, err := data.Encode(reqRPC)
    if err != nil {
      panic(err) 
    }

    err = cReqTransport.Send(b)
    if err != nil {
      return errorHandler(err) 
    }

    rsp, err := cReqTransport.Read()
    if err != nil {
      return errorHandler(err) 
    }

    rspDecode, _ := data.Decode(rsp)
    if rspDecode.Err != "" {
      return errorHandler(errors.New(rspDecode.Err)) 
    }

    if len(rspDecode.Args) == 0 {
      rspDecode.Args = make([]interface{}, container.Type().NumOut())
    }

    numOut := container.Type().NumOut()
    outArgs := make([]reflect.Value, numOut)
    for i := 0; i < numOut; i++ {
      if i != numOut-1 {
        if rspDecode.Args[i] == nil {
          outArgs[i] = reflect.Zero(container.Type().Out(i)) 
        } else {
          outArgs[i] = reflect.ValueOf(rspDecode.Args[i]) 
        }
      } else {
        outArgs[i] = reflect.Zero(container.Type().Out(i)) 
      }
    }

    return outArgs
    
  }

  container.Set(reflect.MakeFunc(container.Type(), f))
}
