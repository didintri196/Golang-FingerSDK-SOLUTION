package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	xj "github.com/basgys/goxml2json"
)

type GetUserInfoResponse struct {
	Row *Row `json:"row"`
}

type Absen struct {
	Row *Row `json:"GetAttLogResponse"`
}
type Row struct {
	Data []Data `json:"row"`
}
type Data struct {
	Status   string `json:"status"`
	WorkCode string `jsno:"workcode"`
	PIN      string `json:"pin"`
	DateTime string `json:"datetime"`
	Verified string `json:"verified"`
}

type Data1 struct {
	Status   string `json:"status"`
	WorkCode string `jsno:"workcode"`
	PIN      string `json:"pin"`
	DateTime string `json:"datetime"`
	Verified string `json:"verified"`
}

type MachineModels struct{}

func main() {
	// var dtabsen *Data
	data, jum := get_absen()
	fmt.Println("JUMLAH : " + strconv.Itoa(jum))
	for _, dtabsen := range data.Row.Data {
		fmt.Println("PIN => " + dtabsen.PIN + "\n")
		fmt.Println("Waktu => " + dtabsen.DateTime + "\n")
		fmt.Println("Verified => " + dtabsen.Verified + "\n")
	}
	// fmt.Println(data)
}

func get_absen() (data *Absen, count int) {
	key := "0"
	ipmesin := "192.168.1.201:80"
	soap_request := "<GetAttLog><ArgComKey xsi:type=\"xsd:integer\">" + key + "</ArgComKey><Arg><PIN xsi:type=\"xsd:integer\">All</PIN></Arg></GetAttLog>"
	jum_len := strconv.Itoa(len(soap_request))
	conn, err := net.Dial("tcp", ipmesin)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "POST /iWsService HTTP/1.0\r\n")
	fmt.Fprintf(conn, "Content-Type: text/xml\r\n")
	fmt.Fprintf(conn, "Content-Length: "+jum_len+"\r\n\r\n")
	fmt.Fprintf(conn, soap_request+"\r\n")
	var buf bytes.Buffer
	io.Copy(&buf, conn)
	parsing := "<GetAttLogResponse>"
	pecah := strings.Split(buf.String(), parsing)
	xml := strings.NewReader(parsing + pecah[1])
	jso, err := xj.Convert(xml)
	if err != nil {
		panic("That's embarrassing...")
	}
	ab := &Absen{}

	if len(jso.String()) >= 129 && len(jso.String()) <= 130 {
		err = json.Unmarshal([]byte(jso.String()), ab)
		count = 1
	} else {
		err = json.Unmarshal([]byte(jso.String()), ab)
		count = len(ab.Row.Data)
	}

	return ab, count
}
