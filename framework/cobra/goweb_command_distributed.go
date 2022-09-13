package cobra

import (
	"log"
	"time"

	"github.com/awaketai/goweb/framework/contract"
	"github.com/robfig/cron/v3"
)

func (c *Command) AddDistributedCronCommand(serviceName, spec string, cmd *Command, holdTime time.Duration) {
	root := c.Root()
	if root.Cron == nil {
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}

	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type:        "distributed-cron",
		Cmd:         cmd,
		Spec:        spec,
		ServiceName: serviceName,
	})

	appService := root.GetContainer().MustMake(contract.AppKey).(contract.App)
	distributeService := root.GetContainer().MustMake(contract.DistributedKey).(contract.Distributed)
	appID := appService.AppID()

	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()

	root.Cron.AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		selectedAppID, err := distributeService.Select(serviceName, appID, holdTime)
		if err != nil {
			return
		}

		if selectedAppID != appID {
			return
		}
		// 如果自已被选择到了，执行这个任务
		err = cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})
}
