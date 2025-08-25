package functions

import (
	// Sesuaikan path di bawah dengan lokasi struct model.BOC Anda
	"errors"

	"github.com/willys-project/mypackage/model"
)

// Key yang didukung
const (
	KeyCommissioner              = "commissioner"
	KeyIndependentCommissioner   = "independent_commissioner"
	KeyPresidentCommissioner     = "president_commissioner"
	KeyVicePresidentCommissioner = "vice_president_commissioner"
)

// mapping key → jabatan
var jabatanMap = map[string]string{
	KeyCommissioner:              "KOMISARIS",
	KeyIndependentCommissioner:   "KOMISARIS INDEPENDEN",
	KeyPresidentCommissioner:     "PRESIDEN KOMISARIS",
	KeyVicePresidentCommissioner: "WAKIL PRESIDEN KOMISARIS",
}

// ReduceBoc groups BOC structs by Jabatan and returns a map of Jabatan to Nama slices.
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

// ExtractValueFromKeyExecutive mengekstrak daftar BOC dari key tertentu pada map executive.
// - keyExecutive: hasil decode JSON (map[string]interface{}) yang memuat array objek pada key yang dicari.
// - keyToFind: salah satu dari konstanta Key* di atas.
// Mengembalikan slice []model.BOC atau error jika tidak ditemukan.
func ExtractValueFromKeyExecutive(keyExecutive map[string]interface{}, keyToFind string) ([]model.BOC, error) {
	jabatan, ok := jabatanMap[keyToFind]
	if !ok {
		return nil, errors.New("unsupported key: " + keyToFind)
	}

	rawVal, exists := keyExecutive[keyToFind]
	if !exists {
		return nil, errors.New("key not found: " + keyToFind)
	}

	arr, ok := rawVal.([]interface{})
	if !ok {
		return nil, errors.New("malformed value for key: " + keyToFind)
	}

	values := make([]model.BOC, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := m["value"].(string)
		lastupdate, _ := m["lastupdate"].(string)

		// only append jika ada nama
		if name != "" {
			values = append(values, model.BOC{
				Jabatan:    jabatan,
				Nama:       name,
				Lastupdate: lastupdate,
			})
		}
	}

	if len(values) == 0 {
		return nil, errors.New("no values found for key: " + keyToFind)
	}
	return values, nil
}
