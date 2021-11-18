package counter

import (
	"fmt"
	"sync/atomic"

	contractsBackgroundtasks "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/backgroundtasks"
	contractsLogger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	servicesBackgroundtasks "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/backgroundtasks"
	"github.com/rs/zerolog/log"
)

// JOB
//------------------------------------------
type counterJob struct {
	counter int32
}

func newCounterJob() *counterJob {
	return &counterJob{}
}
func (j *counterJob) Run() {
	j.incrLocalCounter()
	log.Info().Str("count",
		fmt.Sprintf("%v", j.getLocalCounter())).
		Msg("Background Counter")
}
func (j *counterJob) incrLocalCounter() {
	atomic.AddInt32(&j.counter, 1)
}
func (j *counterJob) getLocalCounter() int32 {
	return atomic.LoadInt32(&j.counter)
}

// Job Provider
//------------------------------------------
type service struct {
	Logger contractsLogger.ILogger
}

func (s *service) GetOneTimeJobs() contractsBackgroundtasks.OneTimeJobs {
	return nil
}
func (s *service) GetScheduledJobs() contractsBackgroundtasks.ScheduledJobs {
	counterJob := newCounterJob()
	cronJob := servicesBackgroundtasks.NewScheduledJob(counterJob, "@every 0h0m5s")
	return servicesBackgroundtasks.NewScheduledJobs(cronJob)

}
