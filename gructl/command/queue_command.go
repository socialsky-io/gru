package command

import (
	"fmt"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/cli"
	etcdclient "github.com/coreos/etcd/client"
	"github.com/gosuri/uitable"
)

func NewQueueCommand() cli.Command {
	cmd := cli.Command{
		Name:   "queue",
		Usage:  "list minion task queue",
		Action: execQueueCommand,
	}

	return cmd
}

// Executes the "queue" command
func execQueueCommand(c *cli.Context) {
	if len(c.Args()) == 0 {
		displayError(errMissingMinion, 64)
	}

	minion := uuid.Parse(c.Args()[0])
	if minion == nil {
		displayError(errInvalidUUID, 64)
	}

	client := newEtcdMinionClientFromFlags(c)

	// Ignore errors about missing queue directory
	queue, err := client.MinionTaskQueue(minion)
	if err != nil {
		if eerr, ok := err.(etcdclient.Error); !ok || eerr.Code != etcdclient.ErrorCodeKeyNotFound {
			displayError(err, 1)
		}
	}

	if len(queue) == 0 {
		return
	}

	table := uitable.New()
	table.MaxColWidth = 40
	table.AddRow("TASK", "COMMAND", "TIME")
	for _, task := range queue {
		table.AddRow(task.TaskID, task.Command, time.Unix(task.TimeReceived, 0))
	}

	fmt.Println(table)
}