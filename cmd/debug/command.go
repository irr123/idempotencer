package debug

import (
	"fmt"

	"github.com/c1rno/idempotencer/pkg/config"
	"github.com/c1rno/idempotencer/pkg/logging"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   `debug`,
	Short: `Don't use it`,
	Run: func(cmd *cobra.Command, args []string) {
		spew.Config.DisablePointerAddresses = true
		spew.Config.DisableCapacities = true
		spew.Config.SortKeys = true
		spew.Config.SpewKeys = true

		logger := logging.NewLogger(0)
		conf, err := config.NewConfig(logger)
		logger.Debug(fmt.Sprintf("Env config: %s", spew.Sdump(conf)), map[string]interface{}{
			"err": err,
		})
		logger = logging.NewLogger(conf.LogLevel)
	},
}