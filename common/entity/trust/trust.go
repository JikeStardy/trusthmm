package trust

import (
	"encoding/json"
	"trust/common"
)

type VehicleTrust struct {
	Id string `json:"id"`
	// TODO 论文里未声明 VT 是整数还是小数，基于 k 属于 0~1，那么先假设 VT 为小数
	VehicleTrust float64 `json:"vt"`
}

func (vt *VehicleTrust) ToJsonBytes() []byte {
	vtAsBytes, err := json.Marshal(vt)
	if err != nil {
		return common.EMPTY_JSON_BYTE
	}
	return vtAsBytes
}

func FromJSON(data string) (VehicleTrust, error) {
	vtAsBytes := []byte(data)
	vt := VehicleTrust{}
	err := json.Unmarshal(vtAsBytes, &vt)
	return vt, err
}
