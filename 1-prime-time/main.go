package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net"

	"github.com/charmbracelet/log"
)

const PORT = 6942

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

type Request struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

func (r *Request) isValid(validity bool) bool {
	if !validity || r.Method != "isPrime" {
		return false
	}
	if math.IsNaN(r.Number) || math.IsInf(r.Number, 0) {
		return false
	}
	return true
}

func (r *Request) resp() Response {
	if r.Number != math.Floor(r.Number) {
		return Response{
			Method: "isPrime",
			Prime:  false,
		}
	}
	if r.Number < 2 {
		return Response{
			Method: "isPrime",
			Prime:  false,
		}
	}
	number := big.NewInt(int64(r.Number))
	return Response{
		Method: "isPrime",
		Prime:  number.ProbablyPrime(20),
	}
}

func (r *Request) sendResponse(conn *net.Conn, log *log.Logger, validity bool) (closeConnection bool) {
	enc := json.NewEncoder(*conn)

	if r.isValid(validity) {
		resp := r.resp()
		log.Info("Sending response", "resp", resp)
		if err := enc.Encode(resp); err != nil {
			log.Fatal(err)
		}
		return false
	} else {
		log.Warn("Request is invalid, closing connection")
		(*conn).Write([]byte("malformed\n"))
		return true
	}
}
func main() {
	log := getNewLogger("main")
	log.Info("PrimeTime!")

	l, err := net.Listen("tcp4", fmt.Sprintf(":%d", PORT))
	if err != nil {
		log.Fatal(err)
	}
	addr := l.Addr().String()
	log.Infof("Listening on %s", addr)

	for i := 1; ; i++ {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		raddr := conn.RemoteAddr().String()
		client := fmt.Sprintf("client-%d : %s", i, raddr)
		log.Infof("Received connection from %s", client)
		go reqHandler(conn, getNewLogger(client))
	}
}

func reqHandler(conn net.Conn, log *log.Logger) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		log.Info("Raw request received", "raw", line)

		var rawMap map[string]any
		if err := json.Unmarshal([]byte(line), &rawMap); err != nil {
			log.Error("Failed to unmarshal JSON", "error", err)
			log.Warn("Malformed request, closing connection")
			conn.Write([]byte("malformed\n"))
			return
		}

		_, hasMethod := rawMap["method"]
		_, hasNumber := rawMap["number"]
		if !hasMethod || !hasNumber {
			log.Warn("Missing required field(s)", "hasMethod", hasMethod, "hasNumber", hasNumber)
			log.Warn("Malformed request, closing connection")
			conn.Write([]byte("malformed\n"))
			return
		}

		req := Request{}
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Error("Failed to unmarshal into Request struct", "error", err)
			log.Warn("Malformed request, closing connection")
			conn.Write([]byte("malformed\n"))
			return
		}

		log.Info("Parsed request", "request", req)
		if closeConnection := req.sendResponse(&conn, log, true); closeConnection {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error("Scanner error", "error", err)
	} else {
		log.Info("Client closed connection cleanly")
	}
}
