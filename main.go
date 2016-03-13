package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"

	"github.com/codegangsta/cli"
)

func run(c *cli.Context) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.String("mqtt-server"))
	opts.SetUsername(c.String("mqtt-username"))
	opts.SetPassword(c.String("mqtt-password"))
	opts.SetClientID("loratestapp")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	token := client.Subscribe("application/+/node/+/rx", 0, func(c *mqtt.Client, msg mqtt.Message) {
		log.Printf("topic: %s, payload: %s", msg.Topic(), msg.Payload())
	})
	token.Wait()
	if token.Error() != nil {
		log.Fatal(token.Error())
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Println(<-sigChan)
}

func main() {
	app := cli.NewApp()
	app.Name = "loratestapp"
	app.Usage = "test application to test the loraserver"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "mqtt-server",
			Usage:  "MQTT server",
			Value:  "tcp://localhost:1883",
			EnvVar: "MQTT_SERVER",
		},
		cli.StringFlag{
			Name:   "mqtt-username",
			Usage:  "MQTT username",
			EnvVar: "MQTT_USERNAME",
		},
		cli.StringFlag{
			Name:   "mqtt-password",
			Usage:  "MQTT password",
			EnvVar: "MQTT_PASSWORD",
		},
	}
	app.Run(os.Args)
}
