package model

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

type CustomError struct {
	message string
}

type DataTrade struct {
	Date      string `json:"date"`
	Open      int64  `json:"open"`
	High      int64  `json:"high"`
	Low       int64  `json:"low"`
	Close     int64  `json:"close"`
	Volume    int64  `json:"volume"`
	Value     int64  `json:"value"`
	Frequency int64  `json:"frequency"`
}
