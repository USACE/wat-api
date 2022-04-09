package wat

type EventConfiguration struct {
	//InputData             string `json:"input_data_directoy"`should be defined by the dag
	OutputDestination        string     `json:"output_destination"` //how do we manage ephemiral
	RealizationNumber        int        `json:"realization_number"`
	EventNumber              int        `json:"event_number"`
	EventTimeWindow          TimeWindow `json:"time_window"`
	NaturalVariabilitySeed   int64      `json:"natural_variability_seed"`
	KnowledgeUncertaintySeed int64      `json:"knowledge_uncertainty_seed"`
}

//should we have a deterministic event configuration and a stochastic event configuration?
