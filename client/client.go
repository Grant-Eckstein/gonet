package gonet

import "crypto/tls"

func Send(remote string, data []byte, responseLength int) ([]byte, error) {
	// InsecureSkipVerify allows for self-signed certs
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Create connection with server
	conn, err := tls.Dial("tcp", remote, conf)
	if err != nil {
		return nil, err
	}

	// Handle error on closure
	defer func(conn *tls.Conn) {
		err = conn.Close()
	}(conn)

	// Send data
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	// Receive response of specified length
	buf := make([]byte, responseLength)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	// Resize buffer and return
	buf = buf[:n]
	return buf, nil
}
