package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/usace/wat-api/config"
	"github.com/usace/wat-api/model"
)

type BatchClient struct {
	Client *batch.Batch
}
type ProvisionedResources struct {
	model.Plugin
	ComputeEnvironmentARN *string
	QueueARN              *string
	JobARN                *string
}

func (awsBatch BatchClient) CreateBatchComputeEnvironment(manifest model.ModelManifest) (*string, error) {
	fmt.Println("creating compute environment for", manifest.ImageAndTag)
	managed := "MANAGED"
	if !manifest.Managed {
		managed = "UNMANAGED"
	}
	computeEnvironment := &batch.CreateComputeEnvironmentInput{
		ComputeEnvironmentName: &manifest.ImageAndTag,
		ComputeResources: &batch.ComputeResource{
			DesiredvCpus:  manifest.DesiredCpus,
			Ec2KeyPair:    &manifest.ModelConfiguration.Name, //not sure we need it
			InstanceRole:  nil,                               //this probably needs to be preset
			InstanceTypes: manifest.InstanceTypes,
			MaxvCpus:      manifest.MaxCpus,
			MinvCpus:      manifest.MinCpus,
			SecurityGroupIds: []*string{
				nil, //needs to be passed in somehow.
			},
			Subnets: []*string{
				nil, //not sure i need this
			},
			Tags: map[string]*string{
				"nil": nil,
			},
			Type: manifest.Type,
		},
		ServiceRole: nil, //this is needed
		State:       aws.String("ENABLED"),
		Type:        &managed,
	}
	output, err := awsBatch.Client.CreateComputeEnvironment(computeEnvironment)
	if err != nil {
		return nil, err
	}
	return output.ComputeEnvironmentArn, nil
}
func (awsBatch BatchClient) SubmitWatTaskAsBatchJob(idx int, path string, dependsOn []*batch.JobDependency, resources []ProvisionedResources, config config.WatConfig) {
	//send a job to batch
	proptags := true
	batchOutput, err := awsBatch.Client.SubmitJob(&batch.SubmitJobInput{
		DependsOn: dependsOn,
		ContainerOverrides: &batch.ContainerOverrides{
			Command: []*string{
				aws.String(".\\main -payload=" + path),
			},
			Environment: config.BatchEnvironmentVariables(),
		},
		JobDefinition:              resources[idx].JobARN, //need to verify this.
		JobName:                    &path,
		JobQueue:                   resources[idx].QueueARN,
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
