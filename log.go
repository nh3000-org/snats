// Copyright 2022-2023 The SNATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A Secure client using NATS messaging system (https://newhorizons3000.org).

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	//"os"
)

var idcount int
var MyLogLang string

type MessageStore struct {
	MSiduuid   string
	MSalias    string
	MShostname string
	MSipadrs   string
	MSmessage  string
	MSnodeuuid string
	MSdate     string
}

var MyLangs = map[string]string{
	"eng-mn-alias": "Intrusion Detection",
	"esp-mn-alias": "Detección de Intrusos",
	"eng-mn-lc":    "Log Connection",
	"esp-mn-lc":    "Conexión de Registro",
	"eng-fm-nhn":   "No Host Name",
	"esp-fm-nhn":   "Sin Nombre de Host",
	"eng-fm-hn":    "Host Name",
	"esp-fm-hn":    "Nombre de Host",
	"eng-fm-mi":    "Mac Ids",
	"esp-fm-mi":    "Identificadores de Mac",
	"eng-fm-ad":    "Address",
	"esp-fm-ad":    "Direccion",
	"eng-fm-ni":    "Node Id",
	"esp-fm-ni":    "Identificación del Nodo",
	"eng-fm-msg":   "Message Id",
	"esp-fm-msg":   "Identificación del mensaje",
	"eng-fm-on":    "On",
	"esp-fm-on":    "En",
	"eng-fm-fm":    "Format Message",
	"esp-fm-fm":    "Dar Formato al Mensaje",
	"eng-fm-con":   "Connection ",
	"esp-fm-con":   "Conexión ",
	"eng-fm-js":    "Jet Stream ",
	"esp-fm-js":    "Corriente en Chorro ",
}

func GetLangs(mylang, mystring string) string {

	value, err := MyLangs[mylang+"-"+mystring]
	//log.Println("GetLangs ", PreferedLanguage+"-"+mystring)
	if err == false {
		return "err"
	}
	return value
}
func FormatMessage(m string) []byte {
	EncMessage := MessageStore{}

	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "\n" + GetLangs(MyLogLang, "fm-nhn")
	} else {
		EncMessage.MShostname = "\n" + GetLangs(MyLogLang, "fm-hn") + " - " + name
	}
	ifas, err := net.Interfaces()
	if err == nil {

		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		EncMessage.MShostname += "\n" + GetLangs(MyLogLang, "fm-mi")
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\n" + GetLangs(MyLogLang, "fm-ad")
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}

	}

	EncMessage.MSalias = GetLangs(MyLogLang, "mn-alias")
	idcount++
	EncMessage.MSnodeuuid = "\n" + GetLangs(MyLogLang, "mn-ni") + " - " + strconv.Itoa(idcount)
	iduuid := uuid.New().String()
	EncMessage.MSiduuid = "\n" + GetLangs(MyLogLang, "mn-msg") + " - " + iduuid[0:8]
	EncMessage.MSdate = "\n" + GetLangs(MyLogLang, "mn-on") + " -" + time.Now().Format(time.UnixDate)
	//EncMessage.MSdate = "\nOn -"
	EncMessage.MSmessage = m
	//EncMessage += m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println(GetLangs(MyLogLang, "mn-fm"), jsonerr)
	}

	return []byte(jsonmsg)

}
func main() {
	//[67367] 2023/07/27 16:55:21.707056 [ERR] 87.236.176.182:54369 - cid:115 - TLS handshake error: EOF
	logLang := flag.String("loglang", "eng", "NATS Language to Use eng esp")
	logPattern := flag.String("logpattern", "[ERR]", "Log Pattern to Identify")
	CA := flag.String("ca", "./ca.pem", "Path to TLS CA Certificate Authority")
	ClientCert := flag.String("clientcert", "./clientcert.pem", "Path to TLS Client Cert")
	ClientKey := flag.String("clientkey", "./clientkey.pem", "Path to TLS Client Key")
	flag.Parse()
	fmt.Println("Usage:")
	fmt.Println("pipe | output of nats-server debug logging to this")
	fmt.Println("anything with [ERR] will be logged to the nats server")
	fmt.Println("plain text")
	fmt.Println("")
	fmt.Println("Run Options:")
	fmt.Println("-loglang: ", *logLang)
	MyLogLang = *logLang
	fmt.Println("-logpattern: ", *logPattern)
	fmt.Println("-ca: ", *CA)
	fmt.Println("-clientkey: ", *ClientCert)
	fmt.Println("-clientkey: ", *ClientKey)

	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				time.Sleep(time.Minute)
			}

		}

		if int64(len(buf)) != 0 {

			if strings.Contains(string(buf), *logPattern) {
				//log.Println("Bufs:", string(buf))
				nc, err := nats.Connect("nats://127.0.0.1:4222", nats.RootCAs(*CA), nats.ClientCert(*ClientCert, *ClientKey))
				if err != nil {
					log.Println(GetLangs(MyLogLang, "mn-con"), err.Error())
				}
				js, err := nc.JetStream()
				if err != nil {
					log.Println(GetLangs(MyLogLang, "mn-js"), err)
				}
				js.Publish("MESSAGES.log", []byte(FormatMessage(string(buf))))
			}
		}

		if err != nil && err != io.EOF {
			log.Println(err)
		}
	}

}
