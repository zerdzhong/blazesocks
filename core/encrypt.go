package core

import (
	"io"
	"net"
)

const (
	BufSize = 1024
)

// Cipher 编码器
type Cipher struct {
	encodePassword *Password
	decodePassword *Password
}

func (cipher *Cipher) encode(data []byte) {
	for i, v := range data {
		data[i] = cipher.encodePassword[v]
	}
}

func (cipher *Cipher) decode(data []byte) {
	for i, v := range data {
		data[i] = cipher.decodePassword[v]
	}
}

//NewCipher 初始化方法
func NewCipher(encodePassword *Password) *Cipher {
	decodePassword := &Password{}

	for i, v := range encodePassword {
		encodePassword[i] = v
		decodePassword[v] = byte(i)
	}

	return &Cipher{
		encodePassword: encodePassword,
		decodePassword: decodePassword,
	}
}

// SecureSocket 传输加密的 Socket
type SecureSocket struct {
	Cipher *Cipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

func (secureSocket *SecureSocket) DecodeRead(conn *net.TCPConn, data []byte) (n int, err error){
	n, err = conn.Read(data)
	if err != nil {
		return
	}

	secureSocket.Cipher.decode(data)

	return
}

func (secureSocket *SecureSocket) EncodeWrite(conn *net.TCPConn, data []byte) (n int, err error) {
	secureSocket.Cipher.encode(data)
	return conn.Write(data)
}

func (secureSocket *SecureSocket) EncodeCopy(dst *net.TCPConn, src *net.TCPConn) error {
	buf := make([]byte, BufSize)
	for {
		readCount, readErr := secureSocket.DecodeRead(src, buf)
		if nil != readErr {
			if readErr != io.EOF {
				return readErr
			} 
			break
		}

		if readCount < 0 {
			continue
		}
		
		writeCount, writeErr := dst.Write(buf[0:readCount])

		if writeErr != nil {
			return writeErr
		}

		if readCount != writeCount {
			return io.ErrShortWrite
		}
	}

	return nil
}
