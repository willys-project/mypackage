// response/structs.go

package response

// ApiResponse adalah struktur untuk response API umum
type ApiResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// BOC adalah struktur untuk informasi BOC
type BOC struct {
	Jabatan    string `json:"jabatan"`
	Nama       string `json:"nama"`
	Lastupdate string `json:"lastUpdate"`
}

// FinancialData adalah struktur untuk data keuangan
type FinancialData struct {
	Url        string `json:"url"`
	LastUpdate string `json:"lastUpdate"`
}
