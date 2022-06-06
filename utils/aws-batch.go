package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

func batchSession() (*batch.Batch, error) {
	var batchClient *batch.Batch
	creds := credentials.NewStaticCredentials(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"",
	)
	cfg := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")).WithCredentials(creds)
	s, err := session.NewSession(cfg)
	if err != nil {
		return batchClient, nil
	}
	batchClient = batch.New(s)
	return batchClient, nil
}

func ProvisionWatResources(definitionFile string, outputDir string) error {

	batchClient, err := batchSession()
	if err != nil {
		return err
	}
	fmt.Println("batchClient", batchClient)

	instructions, err := ioutil.ReadFile(definitionFile)
	if err != nil {
		return err
	}

	var batchPayload AWSBatchPayload
	err = json.Unmarshal(instructions, &batchPayload)
	if err != nil {
		return err
	}

	creationTime := time.Now().Local().Format("2006-01-02_15_04_05")

	// Compute Environment
	computeEnvOutput := filepath.Join(outputDir, fmt.Sprintf("compute-environment-%s.json", creationTime))
	computeEnvData, err := batchPayload.NewComputeEnvironment(batchClient)
	if err != nil {
		return err
	}

	file, err := json.MarshalIndent(computeEnvData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(computeEnvOutput, file, 0644)

	// Job Definition
	jobDefinitionOutput := filepath.Join(outputDir, fmt.Sprintf("job-definition-%s.json", creationTime))
	jobDefinitionData, err := batchPayload.NewJobDefinition(batchClient)
	if err != nil {
		return err
	}

	file, err = json.MarshalIndent(jobDefinitionData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jobDefinitionOutput, file, 0644)

	//  Job Queue
	jobQueueOutput := filepath.Join(outputDir, fmt.Sprintf("job-queue-%s.json", creationTime))
	jobQueueData, err := batchPayload.NewQueue(batchClient, computeEnvData.ComputeEnvironmentArn)
	file, err = json.MarshalIndent(jobQueueData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jobQueueOutput, file, 0644)
	return err
}

func DeleteWatResources(defFile string, outputDir string) error {
	batchClient, err := batchSession()
	if err != nil {
		return err
	}

	instructions, err := ioutil.ReadFile(defFile)
	if err != nil {
		return err
	}

	var batchPayload AWSBatchPayload
	err = json.Unmarshal(instructions, &batchPayload)
	if err != nil {
		return err
	}

	creationTime := time.Now().Local().Format("2006-01-02_15_04_05")

	//-----------------Job Queue-----------------//
	jobQueueOutput := filepath.Join(outputDir, fmt.Sprintf("delete-jobQueue-%s.json", creationTime))
	jobQueueData, err := batchPayload.DeleteQueue(batchClient)
	if err != nil {
		return err
	}

	file, err := json.MarshalIndent(jobQueueData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jobQueueOutput, file, 0644)

	batchPayload.DeleteQueue(batchClient)

	//-----------------Job Definition-----------------//
	jobDefinitionOutput := filepath.Join(outputDir, fmt.Sprintf("delete-jobDef-%s.json", creationTime))
	jobDefinitionData, err := batchPayload.DeleteJobDefinition(batchClient)
	if err != nil {
		return err
	}

	file, err = json.MarshalIndent(jobDefinitionData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jobDefinitionOutput, file, 0644)

	//-----------------Compute Environment-----------------//
	computeEnvironmentOutput := filepath.Join(outputDir, fmt.Sprintf("delete-computeEnv-%s.json", creationTime))
	computeEnvironmentData, err := batchPayload.DeleteJobDefinition(batchClient)
	if err != nil {
		return err
	}

	file, err = json.MarshalIndent(computeEnvironmentData, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(computeEnvironmentOutput, file, 0644)
	return err
}

type AWSBatchPayload struct {
	ComputeEnvironmentFile string `json:"computeEnvironmentFile"`
	JobDefinitionFile      string `json:"jobDefinitionFile"`
	JobQueueFile           string `json:"jobQueueFile"`
}

// Compute Environment Handlers
func (bp AWSBatchPayload) ComputeEnvironmentInputFromJson() (batch.CreateComputeEnvironmentInput, error) {
	var computeEnvironment batch.CreateComputeEnvironmentInput
	instructions, err := ioutil.ReadFile(bp.ComputeEnvironmentFile)
	if err != nil {
		return computeEnvironment, err
	}

	err = json.Unmarshal(instructions, &computeEnvironment)
	if err != nil {
		return computeEnvironment, err
	}
	return computeEnvironment, nil
}

func (bp AWSBatchPayload) ComputeEnvironmentOutputFromJson() (batch.CreateComputeEnvironmentOutput, error) {
	var computeEnvironment batch.CreateComputeEnvironmentOutput
	instructions, err := ioutil.ReadFile(bp.ComputeEnvironmentFile)
	if err != nil {
		return computeEnvironment, err
	}

	err = json.Unmarshal(instructions, &computeEnvironment)
	if err != nil {
		return computeEnvironment, err
	}
	return computeEnvironment, nil
}

func (bp AWSBatchPayload) NewComputeEnvironment(bc *batch.Batch) (output *batch.CreateComputeEnvironmentOutput, err error) {
	computeEnvironment, err := bp.ComputeEnvironmentInputFromJson()
	if err != nil {
		return output, err
	}

	output, err = bc.CreateComputeEnvironment(&computeEnvironment)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (bp AWSBatchPayload) DeleteComputeEnvironment(bc *batch.Batch) (output *batch.DeleteComputeEnvironmentOutput, err error) {
	computeEnvironment, err := bp.ComputeEnvironmentOutputFromJson()
	if err != nil {
		log.Fatal(err)
	}

	updateComputeEnvironmentData := batch.UpdateComputeEnvironmentInput{ComputeEnvironment: computeEnvironment.ComputeEnvironmentName,
		State: aws.String("DISABLED")}

	_, err = bc.UpdateComputeEnvironment(&updateComputeEnvironmentData)
	if err != nil {
		return output, err
	}

	// Wait for AWS to update resources
	time.Sleep(90 * time.Second)

	deleteComputeEnvironmentData := batch.DeleteComputeEnvironmentInput{ComputeEnvironment: computeEnvironment.ComputeEnvironmentName}

	output, err = bc.DeleteComputeEnvironment(&deleteComputeEnvironmentData)
	if err != nil {
		return output, err
	}

	return output, err
}

// Job Definition Handlers
func (bp AWSBatchPayload) JobDefinitionInputFromJson() (batch.RegisterJobDefinitionInput, error) {
	var jobDefinition batch.RegisterJobDefinitionInput
	instructions, err := ioutil.ReadFile(bp.JobDefinitionFile)
	if err != nil {
		return jobDefinition, err
	}

	err = json.Unmarshal(instructions, &jobDefinition)
	if err != nil {
		return jobDefinition, err
	}
	return jobDefinition, nil
}

func (bp AWSBatchPayload) JobDefinitionOutputFromJson() (batch.RegisterJobDefinitionOutput, error) {
	var jobDefinition batch.RegisterJobDefinitionOutput
	instructions, err := ioutil.ReadFile(bp.JobDefinitionFile)
	if err != nil {
		return jobDefinition, err
	}

	err = json.Unmarshal(instructions, &jobDefinition)
	if err != nil {
		return jobDefinition, err
	}
	return jobDefinition, nil
}

func (bp AWSBatchPayload) NewJobDefinition(bc *batch.Batch) (output *batch.RegisterJobDefinitionOutput, err error) {
	jobDefinition, err := bp.JobDefinitionInputFromJson()
	if err != nil {
		return output, err
	}

	output, err = bc.RegisterJobDefinition(&jobDefinition)
	if err != nil {
		return output, err
	}

	// write to output file
	return output, err
}

func (bp AWSBatchPayload) DeleteJobDefinition(bc *batch.Batch) (output *batch.DeregisterJobDefinitionOutput, err error) {
	jobDefinitionData, err := bp.JobDefinitionOutputFromJson()
	if err != nil {
		return output, err
	}

	jobDefinitionDataInput := batch.DeregisterJobDefinitionInput{JobDefinition: jobDefinitionData.JobDefinitionArn}

	_, err = bc.DeregisterJobDefinition(&jobDefinitionDataInput)

	if err != nil {
		return output, err
	}
	return output, err
}

// Job Queue Handlers
func (bp AWSBatchPayload) QueueInputFromJson() (batch.CreateJobQueueInput, error) {

	var jobQueue batch.CreateJobQueueInput
	instructions, err := ioutil.ReadFile(bp.JobQueueFile)
	if err != nil {
		return jobQueue, err
	}

	err = json.Unmarshal(instructions, &jobQueue)
	if err != nil {
		return jobQueue, err
	}
	return jobQueue, nil
}

func (bp AWSBatchPayload) QueueOutputFromJson() (output batch.CreateJobQueueOutput, err error) {

	instructions, err := ioutil.ReadFile(bp.JobQueueFile)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal(instructions, &output)
	if err != nil {
		return output, err
	}
	return output, nil
}

func (bp AWSBatchPayload) NewQueue(bc *batch.Batch, computeEnvironment *string) (output *batch.CreateJobQueueOutput, err error) {
	jobQueue, err := bp.QueueInputFromJson()

	if err != nil {
		return output, err
	}

	// TODO: Think through the jobQueue.ComputeEnvironmentOrder list
	if *computeEnvironment != "" {
		jobQueue.ComputeEnvironmentOrder[0].ComputeEnvironment = computeEnvironment
	}

	output, err = bc.CreateJobQueue(&jobQueue)
	if err != nil {
		return output, err
	}

	return output, err
}

func (bp AWSBatchPayload) DeleteQueue(bc *batch.Batch) (output *batch.DeleteJobQueueOutput, err error) {
	jobQueue, err := bp.QueueOutputFromJson()
	if err != nil {
		return output, err
	}

	updateQueueData := batch.UpdateJobQueueInput{JobQueue: jobQueue.JobQueueName,
		State: aws.String("DISABLED")}

	updatedJobQueueData, err := bc.UpdateJobQueue(&updateQueueData)
	if err != nil {
		fmt.Println("Error....", err)
	}

	// Wait for AWS to update resources
	time.Sleep(30 * time.Second)

	jobQueueData := batch.DeleteJobQueueInput{JobQueue: updatedJobQueueData.JobQueueName}
	_, err = bc.DeleteJobQueue(&jobQueueData)
	if err != nil {
		return output, err
	}

	return output, err
}
