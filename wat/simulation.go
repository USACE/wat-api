package wat

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/USACE/filestore"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
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

func (sj StochasticJob) ProvisionResources(queue *sqs.SQS, awsBatch *batch.Batch) error {
	fmt.Println("provisioning resources...")
	//create a compute environments
	for _, n := range sj.Dag.Nodes {
		fmt.Println("creating compute environment for", n.ImageAndTag)
		managed := "MANAGED"
		if !n.Managed {
			managed = "UNMANAGED"
		}
		computeEnvironment := &batch.CreateComputeEnvironmentInput{
			ComputeEnvironmentName: &n.ImageAndTag,
			ComputeResources: &batch.ComputeResource{
				DesiredvCpus:  n.DesiredCpus,
				Ec2KeyPair:    &n.ModelConfiguration.Name, //not sure we need it
				InstanceRole:  nil,                        //this probably needs to be preset
				InstanceTypes: n.InstanceTypes,
				MaxvCpus:      n.MaxCpus,
				MinvCpus:      n.MinCpus,
				SecurityGroupIds: []*string{
					nil, //needs to be passed in somehow.
				},
				Subnets: []*string{
					nil, //not sure i need this
				},
				Tags: map[string]*string{
					"nil": nil,
				},
				Type: n.Type,
			},
			ServiceRole: nil, //this is needed
			State:       aws.String("ENABLED"),
			Type:        &managed,
		}
		output, err := awsBatch.CreateComputeEnvironment(computeEnvironment)
		if err != nil {
			fmt.Println(err)
		}
		computeEnvironments := make([]*batch.ComputeEnvironmentOrder, 1)
		var order int64 = 0
		computeEnvironments[0] = &batch.ComputeEnvironmentOrder{
			ComputeEnvironment: output.ComputeEnvironmentArn,
			Order:              &order, //lower gets priority?
		}
		//should i be making the batch queues here?
		jobQueueName := fmt.Sprintf("%v_%v", n.ModelConfiguration.Name, n.Plugin.ImageAndTag)
		batchQueueOutput, err := awsBatch.CreateJobQueue(&batch.CreateJobQueueInput{
			ComputeEnvironmentOrder: computeEnvironments,
			JobQueueName:            &jobQueueName,
			Priority:                nil, //higher gets priority
			Tags:                    nil,
			SchedulingPolicyArn:     nil, //if not set FIFO
			State:                   nil, //&batch.JQStatusValid,"VALID"
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(batchQueueOutput)
	}
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
func (sj StochasticJob) GeneratePayloads(sqs *sqs.SQS, fs filestore.FileStore, cache *redis.Client, config config.WatConfig, awsBatch *batch.Batch) error {
	err := sj.ProvisionResources(sqs, awsBatch)
	eventrg := rand.New(rand.NewSource(sj.InitialEventSeed))             //Natural Variability
	realizationrg := rand.New(rand.NewSource(sj.InitialRealizationSeed)) //KnowledgeUncertianty
	if err != nil {
		return err
	}
	nodes := sj.Dag.Nodes
	pluginPayloadStubs := make([]ModelPayload, len(nodes))
	realizationRandomGeneratorByPlugin := make([]*rand.Rand, len(nodes))
	eventRandomGeneratorByPlugin := make([]*rand.Rand, len(nodes))
	for idx, n := range nodes {
		pluginPayloadStubs[idx] = MockModelPayload(sj.Inputsource, n.Plugin) //TODO: remove once DAG is developed to create a payload from a linked manifest
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
			go sj.ProcessDAG(config, j, pluginPayloadStubs, sqs, realizationIndexedSeeds, pluginEventIndexedSeeds, fs, cache, awsBatch)
		}
	}
	return nil
}

func (sj StochasticJob) ProcessDAG(config config.WatConfig, j int, pluginPayloadStubs []ModelPayload, sqs *sqs.SQS, realizationIndexedSeeds []IndexedSeed, eventIndexedSeedsByPlugin []IndexedSeed, fs filestore.FileStore, cache *redis.Client, awsBatch *batch.Batch) {
	key := ""
	dependsOn := make([]*batch.JobDependency, 1)

	for idx := range sj.Dag.Nodes {
		if key != "" {
			//dependency in batch
			dependsOn[0] = &batch.JobDependency{
				JobId: &key,
			}
			//dependency through redis.
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
		//send a job to batch
		proptags := true
		batchOutput, err := awsBatch.SubmitJob(&batch.SubmitJobInput{
			DependsOn:                  dependsOn,
			JobDefinition:              &payload.PluginImageAndTag, //need to verify this.
			JobName:                    &key,
			JobQueue:                   &key,      //i have no queue
			Parameters:                 nil,       //parameters?
			PropagateTags:              &proptags, //i think.
			RetryStrategy:              nil,
			SchedulingPriorityOverride: nil,
			ShareIdentifier:            nil,
			Tags:                       nil,
			Timeout:                    nil,
		})
		fmt.Println("batchoutput", batchOutput)
		if err != nil {
			fmt.Println("batcherror", err)
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
