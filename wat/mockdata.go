package wat

import "time"

func MockModelPayload(inputSource ResourceInfo, plugin Plugin) ModelPayload {
	mconfig := ModelConfiguration{}
	inputs := make([]LinkedDataDescription, 0)
	switch plugin.Name {
	case "fragilitycurveplugin":
		mconfig.Name = "levee_failures"
		mconfig.Alternative = "st. louis river"
		inputs = append(inputs, LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "Project File",
				Parameter: "Project Specification",
				Format:    ".json",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "fc.json",
			},
		})
	case "hydrograph_scaler":
		mconfig.Name = "hydrographs"
		inputs = append(inputs, LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "Project File",
				Parameter: "Project Specification",
				Format:    ".json",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "hsm.json",
			},
		})
		outputs := make([]LinkedDataDescription, 3)
		outputs[0] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "hsm1.csv",
				Parameter: "flow",
				Format:    "csv",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "hsm1.csv",
			},
		}
		outputs[1] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "hsm2.csv",
				Parameter: "flow",
				Format:    "csv",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "hsm2.csv",
			},
		}
		outputs[2] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "hsm3.csv",
				Parameter: "flow",
				Format:    "csv",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "hsm3.csv",
			},
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
		inputs = make([]LinkedDataDescription, 2)
		inputs[0] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "Project File",
				Parameter: "Project Specification",
				Format:    ".yml",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "config_aws.yml",
			},
		}
		inputs[1] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "hsm.csv",
				Parameter: "flow",
				Format:    "csv",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "hsm.csv",
			},
		}
		outputs := make([]LinkedDataDescription, 1)
		outputs[0] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "results-wat.json",
				Parameter: "scalar",
				Format:    "json",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    inputSource.Scheme,
				Authority: inputSource.Authority,
				Fragment:  "results-wat.json",
			},
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
func MockEventConfiguration() EventConfiguration {
	tw := TimeWindow{StartTime: time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local), EndTime: time.Date(2020, time.December, 31, 1, 1, 1, 1, time.Local)}
	event := IndexedSeed{Index: 1, Seed: 5678}
	realization := IndexedSeed{Index: 1, Seed: 1234}
	eventConfiguration := EventConfiguration{
		OutputDestination: ResourceInfo{
			Scheme:    "http",
			Authority: "/minio/runs/realization_1/event_1",
		},
		Realization:     realization,
		Event:           event,
		EventTimeWindow: tw,
	}
	return eventConfiguration
}
