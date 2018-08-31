package command

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func (c *Command) Put(cli *cli.Context) {
	log := c.GetLogger(cli)

	// Check if the job data is not empty
	if len(cli.String("data")) < 1 {
		log.Error("No data passed, use the --data parameter to insert data into a tube.")
		return
	}

	data := []byte(cli.String("data"))

	// Read data from stdin
	var err error
	if cli.String("data") == "-" {
		// Read data from stdin until EOF
		data, err = ioutil.ReadAll(os.Stdin)

		// Exit with an error if we couldnt read from stdin.
		if err != nil {
			log.WithError(err).Error("Could not read from stdin")
		}
	}

	// Build and connect to beanstalkd
	client, err := c.GetBeanstalkdClient(cli)
	if err != nil {
		log.WithError(err).Error("Could not connect to beanstalkd server")
		return
	}

	log.Debug("Selecting correct tube...")

	// Select a tube to insert
	client.Use(cli.String("tube"))
	if err != nil {
		log.WithError(err).WithField("tube", cli.String("tube")).Error("Failed to select tube")
		return
	}

	log.Debug("Inserting job...")

	// Insert the new job
	id, err := client.Put(
		uint32(cli.Int("priority")),
		cli.Duration("delay"),
		cli.Duration("duration"),
		data,
	)

	if err != nil {
		log.WithError(err).WithField("tube", cli.String("tube")).Error("Failed to insert job in queue")
		return
	}

	log.WithFields(logrus.Fields{
		"tube": cli.String("tube"),
		"id":   id,
	}).Info("Succesfully inserted job")

	client.Quit()
}
