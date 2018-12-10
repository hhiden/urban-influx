/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/v01"

	"github.com/hhiden/urban-influx/pkg/config"
	client "github.com/influxdata/influxdb/client/v2"
)

type MessageDumper struct{}

type UOOutMessage struct {
	Location string
	Sensor   string
	Value    float64
}

var conf = config.GetConfig()

func (md *MessageDumper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	marshaller := v01.NewDefaultHTTPMarshaller()
	evt, err := marshaller.FromRequest(r)

	if err == nil {
		iface, ok := evt.Get("data")
		if ok {
			v := iface.(map[string]interface{})
			log.Println(v)
			send(v)
		}
	}
	w.WriteHeader(http.StatusOK)

	/*
		if err == nil {
			//sendReading(event)
			//log.Print(string(reqBytes))
			w.Write(reqBytes)
		} else {
			log.Printf("Error dumping the request: %+v :: %+v", err, r)
		}
	*/
}

func main() {
	http.ListenAndServe(":8080", &MessageDumper{})
}

func send(data map[string]interface{}) {
	location := data["location"].(string)
	sensor := data["sensor"].(string)
	value := data["value"].(float64)

	fmt.Println(location, sensor, value)

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: conf.InfluxURL,
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	tags := map[string]string{"location": location, "sensor": sensor}

	fields := map[string]interface{}{
		sensor: value,
	}

	pt, err := client.NewPoint(location, tags, fields, time.Now())
	if err == nil {
		bp, err2 := client.NewBatchPoints(
			client.BatchPointsConfig{
				Database:  conf.InfluxDB,
				Precision: "s",
			})
		if err2 == nil {
			bp.AddPoint(pt)
			c.Write(bp)
		} else {
			fmt.Println("Error writing point: ", err2.Error())
		}
	} else {
		fmt.Println(err.Error())
	}

}

func sendReading(event cloudevents.Event) {

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: conf.InfluxURL,
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	data := []byte("msg")
	var message UOOutMessage
	json.Unmarshal(data, &message)

	tags := map[string]string{"location": message.Location, "sensor": message.Sensor}

	fields := map[string]interface{}{
		message.Sensor: message.Value,
	}

	pt, err := client.NewPoint(message.Location, tags, fields, time.Now())
	if err == nil {
		bp, err2 := client.NewBatchPoints(
			client.BatchPointsConfig{
				Database:  conf.InfluxDB,
				Precision: "s",
			})
		if err2 == nil {
			bp.AddPoint(pt)
			c.Write(bp)
		} else {
			fmt.Println("Error writing point: ", err2.Error())
		}
	} else {
		fmt.Println(err.Error())
	}

}
