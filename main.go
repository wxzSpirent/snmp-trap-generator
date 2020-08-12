package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/soniah/gosnmp"
	"gopkg.in/yaml.v2"
)

var host = flag.String("host", "127.0.0.1", "Host name of SNMP server")
var port = flag.Int("port", 9900, "Port")
var file = flag.String("file", "trap-input.txt", "File listing traps to be generated")
var community = flag.String("community", "public", "community string")

type trap struct {
	Oid       string `yaml:"oid"`
	Name      string `yaml:"name"`
	Vendor    string `yaml:"vendor"`
	Severity  string `yaml:"severity"`
	EventName string `yaml:"event_name"`
}

func main() {
	flag.Parse()

	// load trap file
	yamlFile, err := os.Open(*file)
	if err != nil {
		fmt.Println("error opening yaml file: ", err)
	}

	defer yamlFile.Close()
	data, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		fmt.Println("error reading yaml file: ", err)
	}

	var traps []trap
	err = yaml.Unmarshal(data, &traps)
	if err != nil {
		fmt.Println("error decoding ", err)
	}

	// setup gosnmp
	gosnmp.Default.Target = *host
	gosnmp.Default.Port = (uint16)(*port)
	gosnmp.Default.Version = gosnmp.Version2c
	gosnmp.Default.Community = *community
	gosnmp.Default.Logger = log.New(os.Stdout, "", 0)

	err = gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer gosnmp.Default.Conn.Close()

	for _, t := range traps {
		pdu := gosnmp.SnmpPDU{
			Name:  ".1.3.6.1.6.3.1.1.4.1.0",
			Type:  gosnmp.ObjectIdentifier,
			Value: t.Oid,
		}
		trap := gosnmp.SnmpTrap{
			Variables: []gosnmp.SnmpPDU{pdu},
		}

		_, err = gosnmp.Default.SendTrap(trap)
		if err != nil {
			log.Fatalf("SendTrap() err: %v", err)
		}
	}
}
