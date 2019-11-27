
package data

import  (
  "bytes"
  "encoding/gob"
)

// RPCdata : format for RPC
type RPCdata struct {
  Name string
  Args []interface{}
  Err string
}

// Encode : sent of the network
func Encode(data RPCdata) ([]byte, error) {
  var buf bytes.Buffer
  encoder := gob.NewEncoder(&buf)
  if err := encoder.Encode(data); err != nil {
    return nil, err 
  }

  return buf.Bytes(), nil
}

// Decode : decode the binary data into the go struct
func Decode(b []byte) (RPCdata, error) {
  buf := bytes.NewBuffer(b)
  decoder := gob.NewDecoder(buf)
  var data RPCdata 
  if err := decoder.Decode(&data); err != nil {
    return RPCdata{}, err 
  }

  return data, nil
}
