package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Memphis struct{}

func newMem() Memphis {
	mem := Memphis{}
	return mem
}

func (m Memphis) Connect(mem *Memphis, host string, username string, connectionToken string, options ...Option) error {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	opts := GetDefaultOptions()
	for _, opt := range options {
		if opt != nil {
			if err := opt(&opts); err != nil {
				return err
			}
		}
	}
	hostCon := normalizeHost(host)
	conn, err := net.Dial("tcp", hostCon)
	if err != nil {
		return err
	}
	opts.client = conn
	tcpConn := tcpCreds{
		Username:     username,
		Broker_creds: connectionToken,
	}

	tcpToBytes := new(bytes.Buffer)
	json.NewEncoder(tcpToBytes).Encode(tcpConn)
	_, err = opts.client.Write(tcpToBytes.Bytes())
	if err != nil {
		return err
	}
	go readTCP(wg, conn)

	tcpChanMsg := make(chan string)

	// go func() {
	// 	data := make([]byte, 4096)
	// 	for range time.Tick(time.Second * 10) {
	// 		_, err := conn.Read(data)
	// 		if err != nil {
	// 			if err != io.EOF {
	// 				fmt.Println("read error:", err)
	// 			}
	// 			break
	// 		}
	// 		// fmt.Println(string(data))
	// 		tcpChanMsg <- string(data)
	// 	}
	// }()

	msg := <-tcpChanMsg
	fmt.Println(reflect.TypeOf(msg))
	fmt.Println(msg)

	// data := make([]byte, 4096)
	// for range time.Tick(time.Second * 10) {
	// 	_, err := conn.Read(data)
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			fmt.Println("read error:", err)
	// 		}
	// 		break
	// 	}
	// 	fmt.Println(string(data))
	// }

	wg.Wait()
	return nil
}

func normalizeHost(host string) string {
	if strings.HasPrefix(host, "http://") {
		res := strings.ReplaceAll(host, "http://", "")
		return res
	} else if strings.HasPrefix(host, "https://") {
		res := strings.ReplaceAll(host, "https://", "")
		return res
	} else {
		return host
	}
}

type Option func(*Options) error

type Options struct {
	host                string
	managementPort      int
	tcpPort             int
	dataPort            int
	username            string
	connectionToken     string
	reconnect           bool
	maxReconnect        int
	reconnectIntervalMs int
	timeoutMs           int
	isConnectionActive  bool
	connectionId        string
	accessToken         string
	client              net.Conn
	brokerConnection    string
	brokerManager       string
	brokerStats         string
	pingTimeout         int
	accessTokenTimeout  int
	reconnectAttempts   int
}

func GetDefaultOptions() Options {
	return Options{
		managementPort:      5555,
		tcpPort:             6666,
		dataPort:            7766,
		reconnect:           true,
		maxReconnect:        3,
		reconnectIntervalMs: 200,
		timeoutMs:           15000,
		isConnectionActive:  false,
	}
}

type tcpCreds struct {
	Username     string
	Broker_creds string
}

func readTCP(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	data := make([]byte, 4096)
	for range time.Tick(time.Second * 10) {
		_, err := conn.Read(data)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		fmt.Println(string(data))
	}
}
