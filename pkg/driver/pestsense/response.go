package pestsense

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/smart-core-os/sc-api/go/traits"
)

type Response struct {
	Id       string `json:"id"`
	Occupied bool   `json:"occupied"`
}

func handleResponse(body []byte, devices map[string]*PestSensor) {
	fmt.Printf("Received message: %s\n", body)
	response := Response{}

	err := json.Unmarshal(body, &response)

	if err != nil {
		fmt.Println("Error unmarshalling")
		return
	}

	fmt.Println("ID: " + response.Id)
	fmt.Println("Occupied: " + strconv.FormatBool(response.Occupied))

	device, exists := devices[response.Id]

	if exists {
		if response.Occupied {
			fmt.Println("Setting occupied for device " + response.Id)
			device.Occupancy.Set(&traits.Occupancy{State: traits.Occupancy_OCCUPIED})
		} else {
			fmt.Println("Setting unoccupied for device " + response.Id)
			device.Occupancy.Set(&traits.Occupancy{State: traits.Occupancy_UNOCCUPIED})
		}
	}
}
