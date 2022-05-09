package wat

import (
	"fmt"
	"math/rand"
	"time"

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
	GeneratePayloads(sqs *sqs.SQS, fs filestore.FileStore, cache *redis.Client) error
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

func (sj StochasticJob) ProvisionResources(queue *sqs.SQS) error {
	fmt.Println("provisioning resources...")
	/*	for _, p := range sj.SelectedPlugins {
			piq := sqs.CreateQueueInput{
				QueueName: &p.Name,
			}
			output, err := queue.CreateQueue(&piq)
			if err != nil {
				return err
			}
			fmt.Println(output.QueueUrl, "created")
		}
	*/
	messages := "messages"
	miq := sqs.CreateQueueInput{
		QueueName: &messages,
	}
	output, err := queue.CreateQueue(&miq)
	if err != nil {
		return err
	}
	fmt.Println(output.QueueUrl, "created")
	events := "events"
	eiq := sqs.CreateQueueInput{
		QueueName: &events,
	}
	output, err = queue.CreateQueue(&eiq)
	if err != nil {
		return err
	}
	fmt.Println(output.QueueUrl, "created")
	return nil
}
func (sj StochasticJob) SendMessage(message string, queue *sqs.SQS, queueName string) error {
	fmt.Println("sending message: " + message)
	input := sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}
	queueURL, err := queue.GetQueueUrl(&input) //fmt.Sprintf("%v/queue/messages", queue.Endpoint)
	if err != nil {
		return err
	}
	fmt.Println("sending message to:", queueURL.QueueUrl)
	output, err := queue.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  aws.String(message),
		QueueUrl:     queueURL.QueueUrl,
	})
	fmt.Println("message sent")
	if err != nil {
		return err
	}
	fmt.Println(output.String())
	return nil
}
func (sj StochasticJob) GeneratePayloads(sqs *sqs.SQS, fs filestore.FileStore, cache *redis.Client, config config.WatConfig) error {
	err := sj.ProvisionResources(sqs)
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
		pluginPayloadStubs[idx] = MockModelPayload(sj.Inputsource, p) //TODO: remove once DAG is developed.
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
			go sj.ProcessDAG(config, j, pluginPayloadStubs, sqs, realizationIndexedSeeds, pluginEventIndexedSeeds, fs, cache)
		}
	}
	return nil
}

func (sj StochasticJob) ProcessDAG(config config.WatConfig, j int, pluginPayloadStubs []ModelPayload, sqs *sqs.SQS, realizationIndexedSeeds []IndexedSeed, eventIndexedSeedsByPlugin []IndexedSeed, fs filestore.FileStore, cache *redis.Client) {
	payloads := make([]ModelPayload, 0)
	key := ""
	for idx, _ := range sj.SelectedPlugins {
		if key != "" {
			for {
				value := cache.Get(key)
				fmt.Println(value)
				if value.Val() == "in progress" {
					time.Sleep(time.Second * 2)
				} else {
					break
				}
			}

		}
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
		//put payload in s3
		path := payload.EventConfiguration.OutputDestination.Authority + "/" + payload.Name + "_payload.yml"
		fmt.Println("putting object in fs:", path)
		_, err = fs.PutObject(path, bytes)
		if err != nil {
			fmt.Println("failure to push payload to filestore:", err)
			panic(err)
		}
		//set status in redis
		key = payload.PluginImageAndTag + "_" + payload.Name + "_R" + fmt.Sprint(payload.Realization.Index) + "_E" + fmt.Sprint(payload.Event.Index)
		cache.Set(key, "in progress", 0)
		//send message to sqs
		err = sj.SendMessage(string(bytes), sqs, "messages") //p.Name
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	for {
		value := cache.Get(key)
		fmt.Println(value)
		if value.Val() == "in progress" {
			time.Sleep(time.Second * 2)
		} else {
			fmt.Println("Realization", realizationIndexedSeeds[0].Index, "Event", eventIndexedSeedsByPlugin[0].Index, "Complete!")
			break
		}
	}

}
