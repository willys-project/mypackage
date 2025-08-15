package model

// APIResponse adalah struktur untuk response API umum
type APIResponse struct {
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
	URL        string `json:"url"`
	LastUpdate string `json:"lastUpdate"`
}

// CustomError represents a custom error structure.
type CustomError struct {
	message string
}

// DataTrade adalah struktur untuk data perdagangan saham harian
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

type ObjKey struct {
	SecCode      string
	Granularity  string
	StartDate    string
	EndDate      string
	AppName      string
	CacheVersion int
}
