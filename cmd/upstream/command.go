package upstream

import (
	"fmt"

	. "github.com/c1rno/idempotencer/cmd/shared"
	"github.com/c1rno/idempotencer/pkg/dto"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/queue"
	_ "github.com/pebbe/zmq4"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   `upstream`,
	Short: `Simple kafka consumer, produces events into 0MQ broker`,
	Run: func(cmd *cobra.Command, args []string) {
		setup, err := InitialSetup()
		helpers.Panicer(err)
		defer setup.Waiter()

		setup.Wg.Add(2)
		client := queue.NewClient(setup.Conf.QueueProducer, setup.Log)
		go func() {
			<-setup.Ctx.Done()
			helpers.Panicer(client.Disconnect())
			setup.Wg.Done()
		}()
		go func() {
			var err errors.Error
			helpers.Panicer(client.Connect())
			id := helpers.UniqIdentity()
			i := 0
			for setup.Ctx.Err() == nil {
				i += 1
				err = client.Push(dto.NewRawMsg(fmt.Sprintf("upstream-%s: %d", id, i)))
				if err == nil {
					_, err = client.Pull()
					for err != nil && setup.Ctx.Err() == nil {
						_, err = client.Pull()
					}
				}
			}
			setup.Wg.Done()
		}()
	},
}
