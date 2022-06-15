package model

type ProvisionedResources struct {
	Plugin
	ComputeEnvironmentARN *string
	JobARN                *string
	QueueARN              *string
}
