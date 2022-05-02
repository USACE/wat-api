package wat

import (
	"fmt"
	"math/rand"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis"
	"github.com/usace/wat-api/config"
	"gopkg.in/yaml.v2"
)

//Job is defined by a manifest, provisions plugin resources, sends messages, and generates event payloads
type Job interface {
	//provisionresources
	ProvisionResources() error
	//sendmessage
	SendMessage(message string, sqs *sqs.SQS) error
	//does this thing need to "run" or "compute"
	GeneratePayloads(sqs *sqs.SQS, fs *filestore.FileStore, cache *redis.Client) error
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
	TotalRealizations            int          `json:"totalrealizations"`
	EventsPerRealization         int          `json:"eventsperrealization"`
	InitialRealizationSeed       int64        `json:"initialrealizationseed"`
	InitialEventSeed             int64        `json:"intitaleventseed"`
	Outputdestination            ResourceInfo `json:"outputdestination"`
	Inputsource                  ResourceInfo `json:"inputsource"`
	DeleteOutputAfterRealization bool         `json:"delete_after_realization"`
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
func (sj StochasticJob) GeneratePayloads(sqs *sqs.SQS, fs *filestore.FileStore, cache *redis.Client, config config.WatConfig) error {
	err := sj.ProvisionResources()
	eventrg := rand.New(rand.NewSource(sj.InitialEventSeed))             //Natural Variability
	realizationrg := rand.New(rand.NewSource(sj.InitialRealizationSeed)) //KnowledgeUncertianty
	if err != nil {
		return err
	}
	plugins := sj.SelectedPlugins
	pluginPayloadStubs := make([]ModelPayload, len(plugins))
	realizationRandomGeneratorByPlugin := make([]*rand.Rand, len(plugins))
	eventRandomGeneratorByPlugin := make([]*rand.Rand, len(plugins))
	for idx, p := range plugins {
		pluginPayloadStubs[idx] = MockModelPayload(sj.Inputsource, p)
		realizationSeeder := realizationrg.Int63()
		eventSeeder := eventrg.Int63()
		realizationRandomGeneratorByPlugin[idx] = rand.New(rand.NewSource(realizationSeeder))
		eventRandomGeneratorByPlugin[idx] = rand.New(rand.NewSource(eventSeeder))
	}
	for i := 0; i < sj.TotalRealizations; i++ { //knowledge uncertainty loop
		realizationIndexedSeeds := make([]IndexedSeed, len(plugins))
		for idx, _ := range plugins {
			realizationSeed := realizationRandomGeneratorByPlugin[idx].Int63()
			realizationIndexedSeeds[idx] = IndexedSeed{Index: i, Seed: realizationSeed}
		}
		for j := 0; j < sj.EventsPerRealization; j++ { //natural variability loop
			//ultimately need to send messages for each task in the event (defined by the dag)
			//event randoms will spawn in unpredictable ways if we dont pre spawn them.
			pluginEventIndexedSeeds := make([]IndexedSeed, len(plugins))
			for idx, _ := range plugins {
				pluginEventSeed := realizationRandomGeneratorByPlugin[idx].Int63()
				pluginEventIndexedSeeds[idx] = IndexedSeed{Index: j, Seed: pluginEventSeed}
			}
			go sj.ProcessDAG(config, j, pluginPayloadStubs, sqs, realizationIndexedSeeds, pluginEventIndexedSeeds)
		}
	}
	return nil
}

func (sj StochasticJob) ProcessDAG(config config.WatConfig, j int, pluginPayloadStubs []ModelPayload, sqs *sqs.SQS, realizationIndexedSeeds []IndexedSeed, eventIndexedSeedsByPlugin []IndexedSeed) {
	payloads := make([]ModelPayload, 0)
	for idx, _ := range sj.SelectedPlugins {
		event := eventIndexedSeedsByPlugin[idx]
		ec := EventConfiguration{
			OutputDestination: ResourceInfo{
				Scheme:    sj.Outputdestination.Scheme, //config.S3_ENDPOINT + "/" + config.S3_BUCKET,
				Authority: fmt.Sprintf("%v%v%v/%v%v", sj.Outputdestination.Authority, "realization_", realizationIndexedSeeds[idx].Index, "event_", event.Index),
			},
			Realization:     realizationIndexedSeeds[idx],
			Event:           event,
			EventTimeWindow: sj.TimeWindow,
		}
		pluginPayloadStubs[idx].EventConfiguration = ec
		payloads = append(payloads, pluginPayloadStubs[idx])
		payload := pluginPayloadStubs[idx]
		for idx, li := range payload.LinkedInputs {
			li.Scheme = ec.OutputDestination.Scheme
			li.Authority = ec.OutputDestination.Authority
			payload.LinkedInputs[idx] = li
		}
		bytes, err := yaml.Marshal(payload)
		if err != nil {
			panic(err)
		}
		//need to join this up with the model information to create a model manifest.
		err = sj.SendMessage(string(bytes), sqs)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}
