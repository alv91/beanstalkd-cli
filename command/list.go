package command

import (
	"fmt"
	"regexp"

	"github.com/urfave/cli"
)

// List displays all tubes or those fulfilling filter.
func (c *Command) List(cli *cli.Context) error {
	log := c.GetLogger(cli)

	// Build and connect to beanstalkd
	client, err := c.GetBeanstalkdClient(cli)
	if err != nil {
		log.WithError(err).Error("Could not connect to beanstalkd server")
		return err
	}

	// Here we get list with all tubes
	log.Debug("Listing tubes")
	tubes, err := client.ListTubes()

	// Filtering tubes list
	var re *regexp.Regexp
	for _, tube := range tubes {
		re = regexp.MustCompile(cli.String("filter"))
		result := re.FindAllString(tube, -1)
		if len(result) == 0 {
			continue
		}
		fmt.Println(tube)
	}

	if err != nil {
		log.WithError(err).WithField("filter", cli.String("filter")).Error("Failed listing tubes")
		return err
	}

	client.Quit()

	return nil
}
