package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/tokopedia/inbox/src/conf"
	"github.com/tokopedia/inbox/src/message"

	gcfg "gopkg.in/gcfg.v1"
	grace "gopkg.in/paytm/grace.v1"
)

type Config struct {
	Database  DatabaseStruct
	Database2 Database2Struct
}

type DatabaseStruct struct {
	SlaveConn string
}

type Database2Struct struct {
	MasterConn string
}

func main() {
	fileLocation := "files/etc/inbox.ini"
	log.Println(fmt.Sprintf("Using configuration : %s", fileLocation))

	// read config
	var config Config
	ok := readConfig(&config, fileLocation)
	if !ok {
		log.Println("Config file error")
	}
	fmt.Println("Database Conn =", config.Database.SlaveConn)

	// store db in global
	conf.InitDB(
		&config.Database.SlaveConn,
		&config.Database2.MasterConn,
	)

	// API: message/v1/inbox
	mux := httprouter.New()
	n := negroni.New()
	n.UseHandler(mux)

	mux.GET("/message/v1/inbox/:user_id", message.ReadInbox)
	// mux.POST("/message/v1/delete", message.Delete)

	grace.Serve(":3000", n)
}

func readConfig(c *Config, filePath string) bool {
	if err := gcfg.ReadFileInto(c, filePath); err != nil {
		log.Printf("%s\n", err)
		return false
	}
	return true
}
