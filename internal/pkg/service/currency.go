package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sovcomhack-inside/internal/pkg/model"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
	"github.com/spf13/viper"
)

type CurrencyService interface {
	ListCurrencies(ctx context.Context, forCurrencyCode string) ([]*dto.CurrencyChangeInfo, error)
	GetCurrencyData(ctx context.Context, forCurrencyCode, base string, ndays int) (*dto.GetCurrencyDataResponse, error)
}

func (svc *service) ListCurrencies(ctx context.Context, forCurrencyCode string) ([]*dto.CurrencyChangeInfo, error) {
	currencyItems := make([]*dto.CurrencyChangeInfo, 0, len(model.CurrencyCodeToName)-1)
	for code, name := range model.CurrencyCodeToName {
		if code != forCurrencyCode {
			currencyItems = append(currencyItems, &dto.CurrencyChangeInfo{
				Code: code,
				Name: name,
			})
		}
	}
	err := findCurrentPrices(ctx, forCurrencyCode, currencyItems)
	if err != nil {
		return nil, err
	}
	return currencyItems, nil
}

func (svc *service) GetCurrencyData(ctx context.Context, forCurrencyCode, base string, ndays int) (*dto.GetCurrencyDataResponse, error) {
	var currencyData []dto.PriceData

	dateNow := time.Now()
	dateDaysBefore := dateNow.AddDate(0, 0, -ndays)

	url := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/timeseries?start_date=%s&end_date=%s", dateDaysBefore.Format("2006-01-02"), dateNow.Format("2006-01-02"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", viper.GetString("service.exchange_rates_secret"))

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	var x map[string]interface{}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &x)
	if err != nil {
		return nil, err
	}
	data := x["rates"].(map[string]interface{})

	i := 0
	for i < ndays {
		dateStr := dateDaysBefore.Format("2006-01-02")
		if _, ok := data[dateStr]; !ok {
			dateDaysBefore = dateDaysBefore.AddDate(0, 0, 1)
			continue
		}
		currentData := data[dateStr].(map[string]interface{})
		basePrice := currentData[base].(float64)
		currencyData = append(currencyData, dto.PriceData{
			Price: basePrice / currentData[forCurrencyCode].(float64),
			Date:  dateStr,
		})
		dateDaysBefore = dateDaysBefore.AddDate(0, 0, 1)
		i++
	}
	return &dto.GetCurrencyDataResponse{
		Code:      forCurrencyCode,
		PriceData: currencyData,
	}, nil
}

func findCurrentPrices(ctx context.Context, base string, items []*dto.CurrencyChangeInfo) error {
	dateNow := time.Now().Format("2006-01-02")
	dateDayBefore := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	dateMonthBefore := time.Now().AddDate(0, -1, 0).Format("2006-01-02")

	url := fmt.Sprintf(
		"https://api.apilayer.com/exchangerates_data/fluctuation?start_date=%s&end_date=%s",
		dateDayBefore,
		dateNow,
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("apikey", viper.GetString("service.exchange_rates_secret"))

	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	var x map[string]interface{}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &x)

	url = fmt.Sprintf(
		"https://api.apilayer.com/exchangerates_data/fluctuation?start_date=%s&end_date=%s",
		dateMonthBefore,
		dateNow,
	)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("apikey", viper.GetString("service.exchange_rates_secret"))

	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	var y map[string]interface{}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &y)
	if err != nil {
		return err
	}

	for i := range items {
		dayBeforeMap := x["rates"].(map[string]interface{})[items[i].Code].(map[string]interface{})
		monthBeforeMap := y["rates"].(map[string]interface{})[items[i].Code].(map[string]interface{})
		baseCurrencyMap := x["rates"].(map[string]interface{})[base].(map[string]interface{})

		items[i].CurrentPrice = baseCurrencyMap["end_rate"].(float64) / dayBeforeMap["end_rate"].(float64)
		items[i].DayChangePct = -dayBeforeMap["change_pct"].(float64)
		items[i].DayChange = dayBeforeMap["change"].(float64)
		items[i].MonthChangePct = -monthBeforeMap["change_pct"].(float64)
		items[i].MonthChange = -monthBeforeMap["change"].(float64)
	}

	return nil
}
