package model

import "time"

func MockModelPayload(inputSource ResourceInfo, outputDestination ResourceInfo, eventParts string, plugin Plugin) ModelPayload {
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
				Fragment:  inputSource.Fragment + "hsm.json",
			},
		})
		inputs = append(inputs, LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "Event Configuration",
				Parameter: "Event Configuration",
				Format:    ".json",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    outputDestination.Scheme,
				Authority: outputDestination.Authority,
				Fragment:  outputDestination.Fragment + "/hydrograph_scaler_Event Configuration.json",
			},
		})
		outputs := make([]LinkedDataDescription, 1)
		outputs[0] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:      "muncie-r1-e2-White-RS-5696.24.csv",
				Parameter: "flow",
				Format:    "csv",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    outputDestination.Scheme,
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/muncie-White-RS-5696.24.csv",
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
	case "ras-mutator":
		mconfig.Name = "Muncie-Mutator"
		inputs = make([]LinkedDataDescription, 2)
		inputs[0] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "self",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.p04.tmp.hdf", //this does not change
			},
		}
		inputs[1] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "/Event Conditions/Unsteady/Boundary Conditions/Flow Hydrographs/River: White  Reach: Muncie  RS: 15696.24",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/muncie-White-RS-5696.24.csv",
			},
		}
		outputs := make([]LinkedDataDescription, 1)
		outputs[0] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "self",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/Muncie.p04.tmp.hdf",
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
	case "ras-unsteady":
		mconfig.Name = "Muncie"
		inputs = make([]LinkedDataDescription, 5)
		inputs[0] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "Muncie.p04.tmp.hdf",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,        //this actually needs to change to output source authority.
				Fragment:  eventParts + "/Muncie.p04.tmp.hdf", //provided by the mutator - changes each event
			},
		}
		inputs[1] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "Muncie.b04",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.b04",
			},
		}
		inputs[2] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "Muncie.prj",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.prj",
			},
		}
		inputs[3] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "Muncie.x04",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.x04",
			},
		}
		inputs[4] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "Muncie.c04",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: inputSource.Authority,
				Fragment:  inputSource.Fragment + "models/Muncie/Muncie.c04",
			},
		}
		outputs := make([]LinkedDataDescription, 3)
		outputs[0] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "Muncie.p04.hdf",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/Muncie.p04.hdf",
			},
		}
		outputs[1] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "Muncie.log",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/Muncie.log",
			},
		}
		outputs[2] = LinkedDataDescription{
			DataDescription: DataDescription{
				Name:   "Muncie.dss",
				Format: "object",
			},
			ResourceInfo: ResourceInfo{
				Scheme:    "s3",
				Authority: outputDestination.Authority,
				Fragment:  eventParts + "/Muncie.dss",
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
