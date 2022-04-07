package wat

//Job is defined by a manifest, provisions plugin resources, sends messages
type Job interface {
	//manifest //shouldnt a job just serialize to json or something?
	//provisionresources
	//sendmessage
	//does this thing need to "run" or "compute"
}

//DeterministicJob implements the Job interface for a Deterministic Compute
type DeterministicJob struct {
	//simulation name?
	//dag
	TimeWindow string `json:"timewindow"`
	//Outputdestination string                 `json:"outputdestination"`
	//Inputsource       string                 `json:"inputsource"`
}

//StochasticJob implements the job interface for a Stochastic Simulation
type StochasticJob struct {
	//dag
	TimeWindow                   string `json:"timewindow"`
	TotalRealizations            int    `json:"totalrealizations"`
	EventsPerRealization         int    `json:"eventsperrealization"`
	InitialRealizationSeed       int64  `json:"initialrealizationseed"`
	InitialEventSeed             int64  `json:"intitaleventseed"`
	Outputdestination            string `json:"outputdestination"`
	Inputsource                  string `json:"inputsource"`
	DeleteOutputAfterRealization bool   `json:"delete_after_realization"`
}
