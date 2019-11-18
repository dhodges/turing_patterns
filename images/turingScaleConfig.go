package images

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// TuringScaleConfig parameters that define a turingScaleGrid
type TuringScaleConfig struct {
	Width  int
	Height int
	Scales []turingScale
}

// DefaultConfig each turingScaleGrid will default to using these params
var DefaultConfig TuringScaleConfig = TuringScaleConfig{
	Scales: []turingScale{
		turingScale{20, 40, 0.04, 1, 2},
		turingScale{10, 20, 0.03, 1, 2},
		turingScale{5, 10, 0.02, 1, 2},
		turingScale{1, 2, 0.01, 1, 2},
	},
}

// ReadConfigFromJSONFile Unmarshal config from a JSON file
func ReadConfigFromJSONFile(filename string) (TuringScaleConfig, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	config := TuringScaleConfig{}
	if err = json.Unmarshal([]byte(file), &config); err != nil {
		log.Fatal(err)
	}

	return config, err
}

// WriteConfigToJSONFile Marshal config to a JSON file
func WriteConfigToJSONFile(cfg TuringScaleConfig, filename string) error {
	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(filename, configJSON, 0644); err != nil {
		log.Fatal(err)
	}

	return err
}
