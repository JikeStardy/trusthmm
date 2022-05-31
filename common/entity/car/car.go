package car

import (
	"encoding/json"
	"github.com/google/uuid"
	"trust/common"
)

type Car struct {
	CarId string `json:"car_id"`
	Info  string `json:"car_info"`
}

func (c *Car) ToJsonBytes() []byte {
	ret, err := json.Marshal(c)
	if err != nil {
		return common.EMPTY_JSON_BYTE
	}
	return ret
}

func NewCar() Car {
	return Car{CarId: uuid.NewString()}
}

func FromJSON(data string) (Car, error) {
	carAsBytes := []byte(data)
	car := Car{}
	err := json.Unmarshal(carAsBytes, &car)
	return car, err
}
