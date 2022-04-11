package wat

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

//Job is defined by a manifest, provisions plugin resources, sends messages, and generates event payloads
type Job interface {
	//provisionresources
	ProvisionResources() error
	//sendmessage
	SendMessage(message string) error
	//does this thing need to "run" or "compute"
	GeneratePayloads() error
}

//DeterministicJob implements the Job interface for a Deterministic Compute
type DeterministicJob struct {
	//simulation name?
	//dag
	TimeWindow `json:"timewindow"`
	//Outputdestination string                 `json:"outputdestination"`
	//Inputsource       string                 `json:"inputsource"`
}

//StochasticJob implements the job interface for a Stochastic Simulation
type StochasticJob struct {
	//dag
	TimeWindow                   `json:"timewindow"`
	TotalRealizations            int    `json:"totalrealizations"`
	EventsPerRealization         int    `json:"eventsperrealization"`
	InitialRealizationSeed       int64  `json:"initialrealizationseed"`
	InitialEventSeed             int64  `json:"intitaleventseed"`
	Outputdestination            string `json:"outputdestination"`
	Inputsource                  string `json:"inputsource"`
	DeleteOutputAfterRealization bool   `json:"delete_after_realization"`
}

func (sj StochasticJob) ProvisionResources() error {
	fmt.Println("provisioning resources...")
	return nil
}
func (sj StochasticJob) SendMessage(message string) error {
	fmt.Println("sending message: " + message)
	return nil
}
func (sj StochasticJob) GeneratePayloads() error {
	err := sj.ProvisionResources()
	if err != nil {
		return err
	}
	eventrg := rand.New(rand.NewSource(sj.InitialEventSeed))             //Natural Variability
	realizationrg := rand.New(rand.NewSource(sj.InitialRealizationSeed)) //KnowledgeUncertianty
	for i := 0; i < sj.TotalRealizations; i++ {                          //knowledge uncertainty loop
		realizationSeed := realizationrg.Int63()
		for j := 0; j < sj.EventsPerRealization; j++ { //natural variability loop
			//ultimately need to send messages for each task in the event (defined by the dag)
			eventSeed := eventrg.Int63()
			ec := EventConfiguration{
				OutputDestination: sj.Outputdestination,
				RealizationNumber: i,
				RealizationSeed:   realizationSeed,
				EventNumber:       j,
				EventSeed:         eventSeed,
				EventTimeWindow:   sj.TimeWindow,
			}
			bytes, err := json.Marshal(ec)
			if err != nil {
				return err
			}
			//need to join this up with the model information to create a model manifest.
			sj.SendMessage(string(bytes))
		}
	}
	return nil
}
