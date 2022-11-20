package dto

type ListCurrenciesRequest struct {
	Code string `query:"code" validate:"required"`
}

type ListCurrenciesResponse struct {
	Currencies []*CurrencyChangeInfo `json:"currencies"`
}

type GetCurrencyDataRequest struct {
	Code             string `query:"code" validate:"required"`
	BaseCurrencyCode string `query:"base" validate:"required"`
	DaysNumber       int    `query:"ndays" validate:"required"`
}

type PriceData struct {
	Price float64 `json:"price"`
	Date  string  `json:"date"`
}

type GetCurrencyDataResponse struct {
	Code      string      `json:"code"`
	PriceData []PriceData `json:"price_data"`
}

type CurrencyChangeInfo struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	CurrentPrice   float64 `json:"current_price"`
	DayChange      float64 `json:"day_change"`
	DayChangePct   float64 `json:"day_change_pct"`
	MonthChangePct float64 `json:"month_change_pct"`
	MonthChange    float64 `json:"month_change"`
}
