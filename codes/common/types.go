package common

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

const (
	cmdPP        = 0x0
	cmdHandshake = 0x1
)

type ClientInfo struct {
	ClientID         string
	PublicIP         string
	PublicPort       uint16
	PublicProtocol   string
	InternalIP       string
	InternalPort     uint16
	InternalProtocol string
}

type ProxyProtocol struct {
	ClientID         string
	PublicIP         string
	PublicPort       uint16
	PublicProtocol   string
	InternalIP       string
	InternalPort     uint16
	InternalProtocol string
}

// 1byte version
// 1byte cmd
// 2bytes length
// length body

func (pp *ProxyProtocol) Encode() ([]byte, error) {
	hdr := make([]byte, 4)
	hdr[0] = 0x0
	hdr[1] = cmdPP

	body, err := json.Marshal(pp)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint16(hdr[2:4], uint16(len(body)))
	return append(hdr, body...), nil
}

func (pp *ProxyProtocol) Decode(reader io.Reader) error {
	hdr := make([]byte, 4)
	_, err := io.ReadFull(reader, hdr)

	if err != nil {
		return err
	}

	cmd := hdr[1]

	if cmd != cmdPP {
		return fmt.Errorf("invalid pp cmd")
	}

	bodyLen := binary.BigEndian.Uint16(hdr[2:4])

	body := make([]byte, bodyLen)

	_, err = io.ReadFull(reader, body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, pp)
	if err != nil {
		return err
	}

	return nil
}

type HandshakeReq struct {
	ClientID string
}

func (hsr *HandshakeReq) Encode() ([]byte, error) {
	hdr := make([]byte, 4)
	hdr[0] = 0x0
	hdr[1] = cmdHandshake

	body, err := json.Marshal(hsr)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint16(hdr[2:4], uint16(len(body)))
	return append(hdr, body...), nil
}

func (hsr *HandshakeReq) Decode(reader io.Reader) error {
	hdr := make([]byte, 4)
	_, err := io.ReadFull(reader, hdr)

	if err != nil {
		return err
	}

	cmd := hdr[1]

	if cmd != cmdHandshake {
		return fmt.Errorf("invalid handshake cmd")
	}

	bodyLen := binary.BigEndian.Uint16(hdr[2:4])

	body := make([]byte, bodyLen)

	_, err = io.ReadFull(reader, body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, hsr)
	if err != nil {
		return err
	}

	return nil
}
