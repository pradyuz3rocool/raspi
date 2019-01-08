package main

import (
	"fmt"

	dht "github.com/d2r2/go-dht"

	logger "github.com/d2r2/go-logger"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
	// logger.InfoLevel,
)

type params struct {
	topic     string
	broker    string
	password  string
	user      string
	id        string
	cleansess bool
	qos       int
	num       int

	action string
	store  string
}
type data struct {
	temp  float32 `json:"temperature"`
	humid float32 `json:"humidity"`
}

func main() {
	defer logger.FinalizeLogger()

	lg.Notify("***************************************************************************************************")
	lg.Notify("*** You can change verbosity of output, to modify logging level of module \"dht\"")
	lg.Notify("*** Uncomment/comment corresponding lines with call to ChangePackageLogLevel(...)")

	//dat := &d1
	credentails := params{"topic", "tcp://192.168.1.92:1883", "password", "username", "host", false, 0, 1, "pub", ":memory"}
	c1 := &credentails

	if c1.topic == "" {
		fmt.Println("Invalid setting for -topic, must not be empty")
		return
	}

	fmt.Printf("Sample Info:\n")
	fmt.Printf("\taction:    %s\n", c1.action)
	fmt.Printf("\tbroker:    %s\n", c1.broker)
	fmt.Printf("\tclientid:  %s\n", c1.id)
	fmt.Printf("\tuser:      %s\n", c1.user)
	fmt.Printf("\tpassword:  %s\n", c1.password)
	fmt.Printf("\ttopic:     %s\n", c1.topic)
	fmt.Printf("\tqos:       %d\n", c1.qos)
	fmt.Printf("\tcleansess: %v\n", c1.cleansess)
	fmt.Printf("\tnum:       %d\n", c1.num)
	fmt.Printf("\tstore:     %s\n", c1.store)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(c1.broker)
	opts.SetClientID(c1.id)
	opts.SetUsername(c1.user)
	opts.SetPassword(c1.password)
	opts.SetCleanSession(c1.cleansess)
	if c1.store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(c1.store))
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("*************************************Sample Publisher Started************************************************")

	sensorType := dht.DHT11

	//var data string
	//var i int
	// print temperature and humidity

	for {

		temperature, humidity, retried, err :=
			dht.ReadDHTxxWithRetry(sensorType, 17, false, 10)
		if err != nil {
			lg.Fatal(err)
		}

		d1 := data{temperature, humidity}

		lg.Infof("data formated ", d1)

		lg.Infof("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)",
			sensorType, temperature, humidity, retried)

		fmt.Println("************************************doing publish*******************************************************")

		token := client.Publish(c1.topic, byte(c1.qos), false, d1)
		token.Wait()
		lg.Infof("done...")
	}

	client.Disconnect(250)
	fmt.Println("Sample Publisher Disconnected")

}
