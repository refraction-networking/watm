package v1

import (
	"bytes"
	"log"
	"os"
	"syscall"
)

func readInboundConfig() (config []byte, err error) {
	// check if /conf/inbound.cfg exists
	file, err := os.Open("/conf/inbound.cfg")
	if err != nil {
		return nil, syscall.EACCES
	}
	return readConfig(file)
}

func readOutboundConfig() (config []byte, err error) {
	// check if /conf/outbound.cfg exists
	file, err := os.Open("/conf/outbound.cfg")
	if err != nil {
		return nil, syscall.EACCES
	}
	return readConfig(file)
}

func readConfig(file *os.File) (config []byte, err error) {
	// read the config file
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Println("readConfig: (*bytes.Buffer).ReadFrom:", err)
		return nil, syscall.EIO
	}

	config = buf.Bytes()

	// close the file
	if err = file.Close(); err != nil {
		err = syscall.EIO
	}
	return
}
