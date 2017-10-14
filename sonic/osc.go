// SonicPi UDP channel baster. Connects to a SonicPi and sends out chunks of music code.
package sonic

import (
	"log"
	"net"
)

type Conn struct {
	Conn *net.UDPConn
}

// Dial connects to the SonicPI service
func Dial(address string) (s *Conn, e error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: conn}, nil
}

// Blast some notes to SonicPi
func (c *Conn) Blast(cmd string) error {
	log.Println("Sonic blast...")

	b := []byte{}
	b = append(b, encodeStr("/run-code")...)
	b = append(b, encodeStr(",ss")...)
	b = append(b, encodeStr("go")...)
	b = append(b, encodeStr(cmd)...)
	b = append(b, encodeStr("")...)

	_, err := c.Conn.Write(b)

	return err
}

func encodeStr(v string) []byte {
	if v == "" {
		return []byte{0, 0}
	}
	b := []byte(v)
	for i := 0; i < 4-len(v)%4; i++ {
		b = append(b, 0)
	}
	return b
}
