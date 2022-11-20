package dto

type ListCurrenciesRequest struct {
	Code string `query:"code" validate:"required"`
}

type ListCurrenciesResponse struct {
	Currencies []*CurrencyChangeInfo `json:"currencies"`
}

type GetCurrencyDataRequest struct {
	Code       string `query:"code" validate:"required"`
	DaysNumber int    `query:"ndays" validate:"required"`
}

type GetCurrencyDataResponse struct {
	Code      string    `json:"code"`
	PriceData []float64 `json:"price_data"`
}

type CurrencyChangeInfo struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	CurrentPrice   float64 `json:"current_price"`
	DayChangePct   float64 `json:"day_change_pct"`
	MonthChangePct float64 `json:"month_change_pct"`
}
