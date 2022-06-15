package model

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/usace/wat-api/config"
	"gopkg.in/yaml.v3"
)

//Job is defined by a manifest, provisions plugin resources, sends messages, and generates event payloads
type Job interface {
	//provisionresources
	ProvisionResources() ([]ProvisionedResources, error)
	//sendmessage
	SendMessage(message string, sqs *sqs.SQS) error
	//does this thing need to "run" or "compute"
	GeneratePayloads(config config.WatConfig, fs filestore.FileStore, aws_batch *batch.Batch) error
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
	//SelectedPlugins              []Plugin `json:"plugins"` //ultimately this needs to be part of the dag somehow
	Dag                          DirectedAcyclicGraph `json:"directed_acyclic_graph"`
	TimeWindow                   `json:"timewindow"`
	TotalRealizations            int          `json:"totalrealizations"`
	EventsPerRealization         int          `json:"eventsperrealization"`
	InitialRealizationSeed       int64        `json:"initialrealizationseed"`
	InitialEventSeed             int64        `json:"intitaleventseed"`
	Outputdestination            ResourceInfo `json:"outputdestination"`
	Inputsource                  ResourceInfo `json:"inputsource"`
	DeleteOutputAfterRealization bool         `json:"delete_after_realization"`
}

func (sj StochasticJob) ProvisionResources(awsBatch *batch.Batch) ([]ProvisionedResources, error) {
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
func (sj StochasticJob) GeneratePayloads(config config.WatConfig, fs filestore.FileStore, awsBatch *batch.Batch) error {
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
		realizationIndexedSeeds := make([]IndexedSeed, len(nodes))
		for idx := range nodes {
			realizationSeed := realizationRandomGeneratorByPlugin[idx].Int63()
			realizationIndexedSeeds[idx] = IndexedSeed{Index: i, Seed: realizationSeed}
		}
		for j := 0; j < sj.EventsPerRealization; j++ { //natural variability loop
			//ultimately need to send messages for each task in the event (defined by the dag)
			//event randoms will spawn in unpredictable ways if we dont pre spawn them.
			pluginEventIndexedSeeds := make([]IndexedSeed, len(nodes))
			for idx := range nodes {
				pluginEventSeed := realizationRandomGeneratorByPlugin[idx].Int63()
				pluginEventIndexedSeeds[idx] = IndexedSeed{Index: j, Seed: pluginEventSeed}
			}
			go sj.ProcessDAG(config, i, j, realizationIndexedSeeds, pluginEventIndexedSeeds, fs, awsBatch, resources)
		}
	}
	fmt.Println("complete")
	return nil
}

func (sj StochasticJob) ProcessDAG(config config.WatConfig, realization int, event int, realizationIndexedSeeds []IndexedSeed, eventIndexedSeedsByPlugin []IndexedSeed, fs filestore.FileStore, awsBatch *batch.Batch, resources []ProvisionedResources) {
	outputDestinationPath := fmt.Sprintf("%v%v%v/%v%v", sj.Outputdestination.Fragment, "realization_", realization, "event_", event)
	for idx, n := range sj.Dag.Nodes {
		fmt.Println(n.ImageAndTag, outputDestinationPath)
		ec := EventConfiguration{
			OutputDestination: ResourceInfo{
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
		payload := Mock2DModelPayload(sj.Inputsource, ec.OutputDestination, outputDestinationPath, n.Plugin)
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
		//submit job to batch.
		//if n.Plugin.Name == "hydrograph_scaler" {
		/*s, err := utils.StartContainer(n.Plugin, path, config.EnvironmentVariables())
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Print(s)
		*/
		//}

	}
	fmt.Println("event", event, "realization", realization, "complete!")
}
