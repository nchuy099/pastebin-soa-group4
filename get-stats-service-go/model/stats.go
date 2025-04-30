package model

// MonthlyStats represents statistics for pastes in a given month
type TimeView struct {
	Time  string `json:"time"`
	Views int64  `json:"views"`
}

type Stats struct {
	PasteID    string     `json:"pasteId,omitempty"`
	TimeViews  []TimeView `json:"timeViews"`
	TotalViews int64      `json:"totalViews"`
	Timezone   string     `json:"timezone,omitempty"`
}

// MonthlyStats represents detailed statistics for pastes in a given month
type MonthlyStats struct {
	TotalViews       int64 `json:"totalViews"`
	AvgViewsPerPaste int64 `json:"avgViewsPerPaste"`
	MinViews         int64 `json:"minViews"`
	MaxViews         int64 `json:"maxViews"`
}

// ResponseData represents a standard API response
type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}
