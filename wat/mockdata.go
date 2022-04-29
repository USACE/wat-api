package wat

func MockModelPayload(inputSource string, plugin Plugin) ModelPayload {
	mconfig := ModelConfiguration{}
	switch plugin.Name {
	case "fragilitycurveplugin":
		paths := make([]string, 1)
		paths[0] = inputSource + "fc.json"
		mconfig.Name = "levee_failures"
		mconfig.ModelConfigurationPaths = paths
	case "hydrograph_scaler":
		paths := make([]string, 1)
		paths[0] = inputSource + "hsm.json"
		mconfig.Name = "hydrographs"
		mconfig.ModelConfigurationPaths = paths
		outputs := make([]PossibleOutput, 3)
		outputs[0] = PossibleOutput{
			Name:      "hsm1.csv",
			Parameter: "flow",
			Format:    "csv",
		}
		outputs[1] = PossibleOutput{
			Name:      "hsm2.csv",
			Parameter: "flow",
			Format:    "csv",
		}
		outputs[2] = PossibleOutput{
			Name:      "hsm3.csv",
			Parameter: "flow",
			Format:    "csv",
		}
		payload := ModelPayload{
			TargetPlugin:       plugin.Name,
			PluginImageAndTag:  plugin.ImageAndTag,
			ModelConfiguration: mconfig,
			ModelLinks: ModelLinks{
				NecessaryOutputs: outputs,
			},
		}
		return payload
	case "hydrograph_stats":
		paths := make([]string, 1)
		paths[0] = inputSource + "config_aws.yml"
		mconfig.Name = "hydrograph_stats"
		mconfig.ModelConfigurationPaths = paths
		inputs := make([]ComputedOutput, 1)
		inputs[0] = ComputedOutput{
			Name:      "hsm.csv",
			Parameter: "flow",
			Format:    "csv",
			ResourceInfo: ResourceInfo{
				Scheme:    "s3?",
				Authority: "/data/realization_0/event_1",
				Fragment:  "hsm.csv",
			},
		}
		outputs := make([]PossibleOutput, 1)
		outputs[0] = PossibleOutput{
			Name:      "results-wat.json",
			Parameter: "scalar",
			Format:    "json",
		}
		payload := ModelPayload{
			TargetPlugin:       plugin.Name,
			PluginImageAndTag:  plugin.ImageAndTag,
			ModelConfiguration: mconfig,
			ModelLinks: ModelLinks{
				LinkedInputs:     inputs,
				NecessaryOutputs: outputs,
			},
		}
		return payload
	}
	payload := ModelPayload{
		TargetPlugin:       plugin.Name,
		PluginImageAndTag:  plugin.ImageAndTag,
		ModelConfiguration: mconfig,
	}
	return payload
}
