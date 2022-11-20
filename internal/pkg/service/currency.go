package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/sovcomhack-inside/internal/pkg/model"
	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

type CurrencyService interface {
	ListCurrencies(ctx context.Context, forCurrencyCode string) ([]*dto.CurrencyChangeInfo, error)
	GetCurrencyData(ctx context.Context, forCurrencyCode string, ndays int) *dto.GetCurrencyDataResponse
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
	err := findCurrentPrices(ctx, currencyItems)
	if err != nil {
		return nil, err
	}
	return currencyItems, nil
}

func (svc *service) GetCurrencyData(ctx context.Context, forCurrencyCode string, ndays int) *dto.GetCurrencyDataResponse {
	var currencyData []float64

	if ndays == 1 {
		currencyData = make([]float64, 24)
	} else {
		currencyData = make([]float64, ndays)
	}
	for i := 0; i < len(currencyData); i++ {
		currencyData[i] = math.Ceil(rand.Float64()+float64(rand.Intn(20)+50)*100) / 100
	}
	return &dto.GetCurrencyDataResponse{
		Code:      forCurrencyCode,
		PriceData: currencyData,
	}
}

func findCurrentPrices(ctx context.Context, items []*dto.CurrencyChangeInfo) error {
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
	req.Header.Set("apikey", "zyAc783IAR5xxXdgxohNEaoCMQdH8oFa")

	if err != nil {
		fmt.Println(err)
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
	err = json.Unmarshal(body, &x)

	url = fmt.Sprintf(
		"https://api.apilayer.com/exchangerates_data/fluctuation?start_date=%s&end_date=%s",
		dateMonthBefore,
		dateNow,
	)
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", "zyAc783IAR5xxXdgxohNEaoCMQdH8oFa")
	if err != nil {
		fmt.Println(err)
	}

	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	var y map[string]interface{}

	body, err = io.ReadAll(res.Body)
	err = json.Unmarshal(body, &y)
	if err != nil {
		return err
	}

	for i := range items {
		dayBeforeMap := x["rates"].(map[string]interface{})[items[i].Code].(map[string]interface{})
		monthBeforeMap := y["rates"].(map[string]interface{})[items[i].Code].(map[string]interface{})

		items[i].CurrentPrice = 1 / dayBeforeMap["end_rate"].(float64)
		items[i].MonthChangePct = -monthBeforeMap["change_pct"].(float64)
		items[i].DayChangePct = -dayBeforeMap["change_pct"].(float64)
	}

	return nil
}
