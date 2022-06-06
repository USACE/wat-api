package wat

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/model"
	"github.com/usace/wat-api/utils"
)

func TestStochasticPayloadGeneration(t *testing.T) {
	tw := model.TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	sj := StochasticJob{

		TimeWindow:                   tw,
		TotalRealizations:            2,
		EventsPerRealization:         10,
		InitialRealizationSeed:       1234,
		InitialEventSeed:             1234,
		Outputdestination:            model.ResourceInfo{Authority: "testing"},
		Inputsource:                  model.ResourceInfo{Authority: "testSettings.InputDataDir"},
		DeleteOutputAfterRealization: false,
	}
	config := config.WatConfig{}
	//fs := filestore.FileStore{}
	err := sj.GeneratePayloads(nil, nil, nil, config, nil)
	if err != nil {
		t.Fail()
	}
}
func mockSimpleDag() model.DirectedAcyclicGraph {
	manifests := make([]model.ModelManifest, 1)
	t := "EC2"
	i := "m2.micro"
	var min int64 = 0
	var desired int64 = 2
	var max int64 = 128
	instance_types := make([]*string, 1)
	instance_types[0] = &i
	manifests[0] = model.ModelManifest{
		ModelComputeResources: model.ModelComputeResources{
			MinCpus:       &min,
			DesiredCpus:   &desired,
			MaxCpus:       &max,
			InstanceTypes: instance_types,
			Type:          &t,
			Managed:       true,
		},
		Plugin: model.Plugin{Name: "fragilitycurveplugin", ImageAndTag: "williamlehman/fragilitycurveplugin:v0.0.7"},
	}
	return model.DirectedAcyclicGraph{
		Nodes: manifests,
	}
}
func mockLoader() utils.ServicesLoader {

	cfg := config.WatConfig{
		APP_PORT:              "8080",
		SKIP_JWT:              false,
		AWS_ACCESS_KEY_ID:     "key",
		AWS_SECRET_ACCESS_KEY: "secret",
		AWS_DEFAULT_REGION:    "us-east-1",
		AWS_S3_REGION:         "us-east-1",
		S3_MOCK:               false,
		S3_BUCKET:             "fake",
		S3_ENDPOINT:           "data",
		S3_DISABLE_SSL:        false,
		S3_FORCE_PATH_STYLE:   false,
		REDIS_HOST:            "bla",
		REDIS_PORT:            "bla",
		REDIS_PASSWORD:        "bla",
		SQS_ENDPOINT:          "bla",
	}
	ldr, err := utils.InitLoaderWithConfig("", cfg)
	if err != nil {
		fmt.Print(err)
	}
	return ldr
}
func TestBatchComputeEnvironmentGeneration(t *testing.T) {
	loader := mockLoader()
	awsBatch, err := loader.InitBatch()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	dag := mockSimpleDag()
	fmt.Println("provisioning resources...")
	resources := make([]utils.ProvisionedResources, len(dag.Nodes))
	//create a compute environments
	for idx, n := range dag.Nodes {
		resources[idx] = utils.ProvisionedResources{
			Plugin: n.Plugin,
		}
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
		resources[idx].ComputeEnvironmentARN = output.ComputeEnvironmentArn
	}
}
func TestBatchComputeEnvironmentAndJobDefinitonAndJobQueue(t *testing.T) {
	loader := mockLoader()
	awsBatch, err := loader.InitBatch()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	dag := mockSimpleDag()
	fmt.Println("provisioning resources...")
	resources := make([]utils.ProvisionedResources, len(dag.Nodes))
	//create a compute environments
	for idx, n := range dag.Nodes {
		resources[idx] = utils.ProvisionedResources{
			Plugin: n.Plugin,
		}
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
		resources[idx].ComputeEnvironmentARN = output.ComputeEnvironmentArn
		computeEnvironments := make([]*batch.ComputeEnvironmentOrder, 1)
		var order int64 = 0
		computeEnvironments[0] = &batch.ComputeEnvironmentOrder{
			ComputeEnvironment: output.ComputeEnvironmentArn,
			Order:              &order, //lower gets priority?
		}
		//register the job

		inputRegister := &batch.RegisterJobDefinitionInput{
			ContainerProperties: &batch.ContainerProperties{
				Command: []*string{
					aws.String(".\\main -payload=" + "pathtopayload.yml"), //how do i pass the command line argument to the path dynamically?
				}, //overridden with container overrides in the submitjobinput
				Image: aws.String("busybox"),
				ResourceRequirements: []*batch.ResourceRequirement{
					{
						Type:  aws.String("MEMORY"),
						Value: aws.String("2048"),
					},
					{
						Type:  aws.String("VCPU"),
						Value: aws.String("1"),
					},
				},
			},
			JobDefinitionName: aws.String("execute go container with payload"),
			Tags: map[string]*string{
				"Agency": aws.String("USACE"),
				"User":   aws.String("JaneDoe"),
			},
			Type: aws.String("container"),
		}
		jobRegisterOutput, err := awsBatch.RegisterJobDefinition(inputRegister)
		if err != nil {
			fmt.Println(err)
			//panic(err)
		}
		resources[idx].JobARN = jobRegisterOutput.JobDefinitionArn
		//create a batch queue
		jobQueueName := fmt.Sprintf("%v_%v", n.ModelConfiguration.Name, n.Plugin.ImageAndTag)
		batchQueueOutput, err := awsBatch.CreateJobQueue(&batch.CreateJobQueueInput{
			ComputeEnvironmentOrder: computeEnvironments,
			JobQueueName:            &jobQueueName,
			Priority:                aws.Int64(1), //higher gets priority
			Tags:                    nil,
			SchedulingPolicyArn:     nil, //if not set FIFO
			State:                   nil, //&batch.JQStatusValid,"VALID"
		})
		if err != nil {
			fmt.Println(err)
		}
		resources[idx].QueueARN = batchQueueOutput.JobQueueArn
	}
}
