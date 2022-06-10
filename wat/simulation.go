package wat

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/model"
	"github.com/usace/wat-api/utils"
	"gopkg.in/yaml.v3"
)

//Job is defined by a manifest, provisions plugin resources, sends messages, and generates event payloads
type Job interface {
	//provisionresources
	ProvisionResources() ([]utils.ProvisionedResources, error)
	//sendmessage
	SendMessage(message string, sqs *sqs.SQS) error
	//does this thing need to "run" or "compute"
	GeneratePayloads(sqs *sqs.SQS, fs filestore.FileStore, cache *redis.Client) error
}

//DeterministicJob implements the Job interface for a Deterministic Compute
type DeterministicJob struct {
	//simulation name?
	//dag
	model.TimeWindow `json:"timewindow"`
	//Outputdestination string                 `json:"outputdestination"`
	//Inputsource       string                 `json:"inputsource"`
}

//StochasticJob implements the job interface for a Stochastic Simulation
type StochasticJob struct {
	//dag
	//SelectedPlugins              []Plugin `json:"plugins"` //ultimately this needs to be part of the dag somehow
	Dag                          model.DirectedAcyclicGraph `json:"directed_acyclic_graph"`
	model.TimeWindow             `json:"timewindow"`
	TotalRealizations            int                `json:"totalrealizations"`
	EventsPerRealization         int                `json:"eventsperrealization"`
	InitialRealizationSeed       int64              `json:"initialrealizationseed"`
	InitialEventSeed             int64              `json:"intitaleventseed"`
	Outputdestination            model.ResourceInfo `json:"outputdestination"`
	Inputsource                  model.ResourceInfo `json:"inputsource"`
	DeleteOutputAfterRealization bool               `json:"delete_after_realization"`
}

func (sj StochasticJob) ProvisionResources(awsBatch *batch.Batch) ([]utils.ProvisionedResources, error) {
	fmt.Println("provisioning resources...")
	return nil, nil
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
func (sj StochasticJob) GeneratePayloads(sqs *sqs.SQS, fs filestore.FileStore, cache *redis.Client, config config.WatConfig, awsBatch *batch.Batch) error {
	//provision resources
	resources, err := sj.ProvisionResources(awsBatch)
	//create random seed generators.
	eventrg := rand.New(rand.NewSource(sj.InitialEventSeed))             //Natural Variability
	realizationrg := rand.New(rand.NewSource(sj.InitialRealizationSeed)) //KnowledgeUncertianty
	if err != nil {
		return err
	}
	nodes := sj.Dag.Nodes
	realizationRandomGeneratorByPlugin := make([]*rand.Rand, len(nodes))
	eventRandomGeneratorByPlugin := make([]*rand.Rand, len(nodes))
	for idx := range nodes {
		realizationSeeder := realizationrg.Int63()
		eventSeeder := eventrg.Int63()
		realizationRandomGeneratorByPlugin[idx] = rand.New(rand.NewSource(realizationSeeder))
		eventRandomGeneratorByPlugin[idx] = rand.New(rand.NewSource(eventSeeder))
	}
	for i := 0; i < sj.TotalRealizations; i++ { //knowledge uncertainty loop
		realizationIndexedSeeds := make([]model.IndexedSeed, len(nodes))
		for idx := range nodes {
			realizationSeed := realizationRandomGeneratorByPlugin[idx].Int63()
			realizationIndexedSeeds[idx] = model.IndexedSeed{Index: i, Seed: realizationSeed}
		}
		for j := 0; j < sj.EventsPerRealization; j++ { //natural variability loop
			//ultimately need to send messages for each task in the event (defined by the dag)
			//event randoms will spawn in unpredictable ways if we dont pre spawn them.
			pluginEventIndexedSeeds := make([]model.IndexedSeed, len(nodes))
			for idx := range nodes {
				pluginEventSeed := realizationRandomGeneratorByPlugin[idx].Int63()
				pluginEventIndexedSeeds[idx] = model.IndexedSeed{Index: j, Seed: pluginEventSeed}
			}
			go sj.ProcessDAG(config, i, j, sqs, realizationIndexedSeeds, pluginEventIndexedSeeds, fs, cache, awsBatch, resources)
		}
	}
	fmt.Println("complete")
	return nil
}

func (sj StochasticJob) ProcessDAG(config config.WatConfig, realization int, event int, sqs *sqs.SQS, realizationIndexedSeeds []model.IndexedSeed, eventIndexedSeedsByPlugin []model.IndexedSeed, fs filestore.FileStore, cache *redis.Client, awsBatch *batch.Batch, resources []utils.ProvisionedResources) {
	outputDestinationPath := fmt.Sprintf("%v%v%v/%v%v", sj.Outputdestination.Fragment, "realization_", realization, "event_", event)
	for idx, n := range sj.Dag.Nodes {
		fmt.Println(n.ImageAndTag, outputDestinationPath)
		ec := model.EventConfiguration{
			OutputDestination: model.ResourceInfo{
				Scheme:    sj.Outputdestination.Scheme,
				Authority: sj.Outputdestination.Authority,
				Fragment:  outputDestinationPath,
			},
			Realization:     realizationIndexedSeeds[idx],
			Event:           eventIndexedSeedsByPlugin[idx],
			EventTimeWindow: sj.TimeWindow,
		}
		//write event configuration to s3.
		ecbytes, err := json.Marshal(ec)
		if err != nil {
			panic(err)
		}
		path := outputDestinationPath + "/" + n.Plugin.Name + "_Event Configuration.json"
		fmt.Println("putting object in fs:", path)
		_, err = fs.PutObject(path, ecbytes)
		if err != nil {
			fmt.Println("failure to push event configuration to filestore:", err)
			panic(err)
		}
		payload := model.Mock2DModelPayload(sj.Inputsource, ec.OutputDestination, outputDestinationPath, n.Plugin)
		bytes, err := yaml.Marshal(payload)
		if err != nil {
			panic(err)
		}
		//put payload in s3
		path = outputDestinationPath + "/" + n.Plugin.Name + "_payload.yml"
		fmt.Println("putting object in fs:", path)
		_, err = fs.PutObject(path, bytes)
		if err != nil {
			fmt.Println("failure to push payload to filestore:", err)
			panic(err)
		}
		//set status in redis
		//key = payload.Alternative + "_" + payload.Name + "_R" + fmt.Sprint(payload.EventConfiguration().Realization.Index) + "_E" + fmt.Sprint(payload.EventConfiguration().Event.Index)
		//cache.Set(key, "in progress", 0)
		//send message to sqs
		/*mess := model.PayloadMessage{
			Plugin:      n.Plugin,
			PayloadPath: path,
		}
		byt, err := yaml.Marshal(mess)
		if err != nil {
			panic(err)
		}
		err = sj.SendMessage(string(byt), sqs, "messages")
		*/

		//submit job to batch.
		//if n.Plugin.Name == "hydrograph_scaler" {
		s, err := utils.StartContainer(n.Plugin, path, config.EnvironmentVariables())
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Print(s)
		//}

	}
	fmt.Println("event", event, "realization", realization, "complete!")
}
