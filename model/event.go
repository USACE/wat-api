package model

type EventConfiguration struct {
	OutputDestination ResourceInfo `json:"output_destination" yaml:"output_destination"` //how do we manage ephemiral
	Realization       IndexedSeed  `json:"realization" yaml:"realization"`               //knowledge uncertainty
	Event             IndexedSeed  `json:"event" yaml:"event"`                           //natural variability
	EventTimeWindow   TimeWindow   `json:"time_window" yaml:"time_window"`
}

func (ec EventConfiguration) ToInput() LinkedDataDescription {
	return LinkedDataDescription{
		DataDescription: DataDescription{
			Name:      "Event Configuration",
			Parameter: "Event Specification",
			Format:    ".json",
		},
		ResourceInfo: ResourceInfo{
			Scheme:    "https",
			Authority: "needs to be replaced",
			Fragment:  "EventConfiguration.json",
		},
	}
}

type IndexedSeed struct {
	Index int   `json:"index" yaml:"index"`
	Seed  int64 `json:"seed" yaml:"seed"`
}

//should we have a deterministic event configuration and a stochastic event configuration?
