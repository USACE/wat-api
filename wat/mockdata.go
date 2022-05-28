package wat

func MockModelPayload(inputSource ResourceInfo, plugin Plugin) ModelPayload {
	mconfig := ModelConfiguration{}
	inputs := make([]ComputedOutput, 0)
	switch plugin.Name {
	case "fragilitycurveplugin":
		mconfig.Name = "levee_failures"
		mconfig.Alternative = "st. louis river"
		inputs = append(inputs, ComputedOutput{
			Name:      "Project File",
			Parameter: "Project Specification",
			Format:    ".json",
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "fc.json",
			},
		})
	case "hydrograph_scaler":
		mconfig.Name = "hydrographs"
		inputs = append(inputs, ComputedOutput{
			Name:      "Project File",
			Parameter: "Project Specification",
			Format:    ".json",
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "hsm.json",
			},
		})
		outputs := make([]Output, 3)
		outputs[0] = Output{
			Name:      "hsm1.csv",
			Parameter: "flow",
			Format:    "csv",
		}
		outputs[1] = Output{
			Name:      "hsm2.csv",
			Parameter: "flow",
			Format:    "csv",
		}
		outputs[2] = Output{
			Name:      "hsm3.csv",
			Parameter: "flow",
			Format:    "csv",
		}
		payload := ModelPayload{
			ModelConfiguration: mconfig,
			ModelLinks: ModelLinks{
				LinkedInputs:     inputs,
				NecessaryOutputs: outputs,
			},
		}
		return payload
	case "hydrograph_stats":
		mconfig.Name = "hydrograph_stats"
		inputs = make([]ComputedOutput, 2)
		inputs[0] = ComputedOutput{
			Name:      "Project File",
			Parameter: "Project Specification",
			Format:    ".yml",
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "config_aws.yml",
			},
		}
		inputs[1] = ComputedOutput{
			Name:      "hsm.csv",
			Parameter: "flow",
			Format:    "csv",
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "hsm.csv",
			},
		}
		outputs := make([]Output, 1)
		outputs[0] = Output{
			Name:      "results-wat.json",
			Parameter: "scalar",
			Format:    "json",
		}
		payload := ModelPayload{
			ModelConfiguration: mconfig,
			ModelLinks: ModelLinks{
				LinkedInputs:     inputs,
				NecessaryOutputs: outputs,
			},
		}
		return payload
	}
	payload := ModelPayload{
		ModelConfiguration: mconfig,
		ModelLinks: ModelLinks{
			LinkedInputs: inputs,
		},
	}
	return payload
}
