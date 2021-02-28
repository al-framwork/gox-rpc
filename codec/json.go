package codec

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

// create json Codec
func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}

func (jsonCodec *JsonCodec) Close() error {
	return jsonCodec.conn.Close()
}

// read header
func (jsonCodec *JsonCodec) ReadHeader(header *Header) error {
	return jsonCodec.dec.Decode(header)
}

// read body
func (jsonCodec *JsonCodec) ReadBody(body interface{}) error {
	return jsonCodec.dec.Decode(body)
}

func (jsonCodec *JsonCodec) Write(header *Header, body interface{}) (err error) {
	defer func() {
		_ = jsonCodec.buf.Flush()
		if err != nil {
			_ = jsonCodec.Close()
		}
	}()

	if err = jsonCodec.enc.Encode(header); err != nil {
		log.Println("rpc codec: json error encoding invoker:", err)
		return err
	}

	if err = jsonCodec.enc.Encode(body); err != nil {
		log.Println("rpc codec: json error encoding body:", err)
		return err
	}
	return nil
}
