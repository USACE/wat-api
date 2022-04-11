package wat

type EventConfiguration struct {
	OutputDestination string     `json:"output_destination"` //how do we manage ephemiral
	RealizationNumber int        `json:"realization_number"`
	EventNumber       int        `json:"event_number"`
	EventTimeWindow   TimeWindow `json:"time_window"`
	EventSeed         int64      `json:"event_seed"`       //natural variability
	RealizationSeed   int64      `json:"realization_seed"` //knowledge uncertainty
}

//should we have a deterministic event configuration and a stochastic event configuration?
