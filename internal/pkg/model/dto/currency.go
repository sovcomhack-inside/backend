package dto

type ListCurrenciesRequest struct {
	Code string `query:"code" validate:"required"`
}

type ListCurrenciesResponse struct {
	Currencies []*CurrencyChangeInfo `json:"currencies"`
}

type CurrencyChangeInfo struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	CurrentPrice   float64 `json:"current_price"`
	DayChangePct   float64 `json:"day_change_pct"`
	MonthChangePct float64 `json:"month_change_pct"`
}
