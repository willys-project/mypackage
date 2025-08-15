package functions

import "github.com/willys-project/mypackage/model"

func ReduceBoc(boc []model.BOC) map[string][]string {
	data := make(map[string][]string)

	for _, currVal := range boc {
		if _, exists := data[currVal.Jabatan]; !exists {
			data[currVal.Jabatan] = []string{currVal.Nama}
		} else {
			data[currVal.Jabatan] = append(data[currVal.Jabatan], currVal.Nama)
		}
	}

	return data
}
