package model

// MonthlyStats represents statistics for pastes in a given month
type MonthlyStats struct {
	TotalPastes      int64   `json:"totalPastes"`
	TotalViews       int64   `json:"totalViews"`
	AvgViewsPerPaste float64 `json:"avgViewsPerPaste"`
	MinViews         int     `json:"minViews"`
	MaxViews         int     `json:"maxViews"`
	ActivePastes     int     `json:"activePastes"`
	ExpiredPastes    int     `json:"expiredPastes"`
}

// ResponseData represents a standard API response
type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}
