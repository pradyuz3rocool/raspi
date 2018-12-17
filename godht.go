package main

import (
	"fmt"
	"time"

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
	payload   string
	action    string
	store     string
}

func main() {
	defer logger.FinalizeLogger()

	lg.Notify("***************************************************************************************************")
	lg.Notify("*** You can change verbosity of output, to modify logging level of module \"dht\"")
	lg.Notify("*** Uncomment/comment corresponding lines with call to ChangePackageLogLevel(...)")
	lg.Notify("***************************************************************************************************")
	// Uncomment/comment next line to suppress/increase verbosity of output
	// logger.ChangePackageLogLevel("dht", logger.InfoLevel)
	for {

		time.Sleep(time.Second * 10)
		sensorType := dht.DHT11

		var temp string

		temperature, humidity, retried, err :=
			dht.ReadDHTxxWithRetry(sensorType, 17, false, 10)
		if err != nil {
			lg.Fatal(err)
		}

		// print temperature and humidity
		lg.Infof("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)",
			sensorType, temperature, humidity, retried)

		temp = fmt.Sprintf("%f", temperature) // s == "123.456000"

		credentails := params{"topic", "tcp://192.168.1.74:1883", "password", "username", "host", false, 0, 1, temp, "pub", ":memory"}
		c1 := &credentails
		//	flag.Parse()

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
		fmt.Printf("\tmessage:   %s\n", c1.payload)
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

		//if *action == "pub" {
		client := MQTT.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("Sample Publisher Started")
		for i := 0; i < c1.num; i++ {
			fmt.Println("---- doing publish ----")
			token := client.Publish(c1.topic, byte(c1.qos), false, c1.payload)
			token.Wait()
		}

		client.Disconnect(250)
		fmt.Println("Sample Publisher Disconnected")
		//}
	}
}
