package main

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func main() {

	log.Info("Create new cron")
	c := cron.New()

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if filepath.Ext(file.Name()) == ".properties" {

			var properties, _ = properties.LoadFile(file.Name(), properties.UTF8)

			get, ok := properties.Get("schedule")

			if ok {
				c.AddFunc(get, func() {

					client := &http.Client{}

					api, o := properties.Get("healthApi")

					if !o {
						fmt.Printf("Invalide health check url")
						return
					}

					req, err := http.NewRequest("GET", api, nil)
					if err != nil {
						fmt.Print(err.Error())
						return
					}
					req.Header.Add("Accept", "application/json")
					req.Header.Add("Content-Type", "application/json")
					resp, err := client.Do(req)
					if err != nil {
						fmt.Print(err.Error())
						return
					}

					defer resp.Body.Close()
					bodyBytes, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Print(err.Error())
						return
					}
					response := string(bodyBytes)
					fmt.Printf("%+v\n", response)

				})
			}
		}
	}

	c.Start()

	select {}

}

func printCronEntries(cronEntries []cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}
