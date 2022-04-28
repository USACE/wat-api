package wat

import (
	"fmt"
	"math/rand"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"gopkg.in/yaml.v2"
)

//Job is defined by a manifest, provisions plugin resources, sends messages, and generates event payloads
type Job interface {
	//provisionresources
	ProvisionResources() error
	//sendmessage
	SendMessage(message string, sqs *sqs.SQS) error
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
	SelectedPlugins              []Plugin `json:"plugins"` //ultimately this needs to be part of the dag somehow
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
func (sj StochasticJob) SendMessage(message string, queue *sqs.SQS) error {
	fmt.Println("sending message: " + message)
	queueURL := fmt.Sprintf("%v/queue/messages", queue.Endpoint)
	fmt.Println("sending message to:", queueURL)
	output, err := queue.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  aws.String(message),
		QueueUrl:     &queueURL,
	})
	fmt.Println("message sent")
	if err != nil {
		return err
	}
	fmt.Println(output.String())
	return nil
}
func (sj StochasticJob) GeneratePayloads(sqs *sqs.SQS) ([]ModelPayload, error) {
	err := sj.ProvisionResources()
	payloads := make([]ModelPayload, 0)
	fcp := Plugin{"fragilitycurveplugin", "williamlehman/fragilitycurveplugin:v0.0.2"}
	payload := MockModelPayload(sj.Inputsource, fcp)
	if err != nil {
		return payloads, err
	}
	eventrg := rand.New(rand.NewSource(sj.InitialEventSeed))             //Natural Variability
	realizationrg := rand.New(rand.NewSource(sj.InitialRealizationSeed)) //KnowledgeUncertianty
	for i := 0; i < sj.TotalRealizations; i++ {                          //knowledge uncertainty loop
		realizationSeed := realizationrg.Int63()
		realization := IndexedSeed{Index: i, Seed: realizationSeed}
		for j := 0; j < sj.EventsPerRealization; j++ { //natural variability loop
			//ultimately need to send messages for each task in the event (defined by the dag)
			eventSeed := eventrg.Int63()
			event := IndexedSeed{Index: j, Seed: eventSeed}
			ec := EventConfiguration{
				OutputDestination: fmt.Sprintf("%v%v%v/%v%v", sj.Outputdestination, "realization_", realization.Index, "event_", event.Index),
				Realization:       realization,
				Event:             event,
				EventTimeWindow:   sj.TimeWindow,
			}
			payload.EventConfiguration = ec
			payloads = append(payloads, payload)
			bytes, err := yaml.Marshal(payload)
			if err != nil {
				return payloads, err
			}
			//need to join this up with the model information to create a model manifest.
			err = sj.SendMessage(string(bytes), sqs)
			if err != nil {
				fmt.Println(err)
				return payloads, err
			}
		}
	}
	return payloads, nil
}
