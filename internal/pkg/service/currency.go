package service

import (
	"context"
	"math"
	"math/rand"

	"github.com/sovcomhack-inside/internal/pkg/model/dto"
)

type CurrencyService interface {
	ListCurrencies(ctx context.Context, forCurrencyCode string) []*dto.CurrencyChangeInfo
	GetCurrencyData(ctx context.Context, forCurrencyCode string, ndays int) *dto.GetCurrencyDataResponse
}

func (svc *service) ListCurrencies(ctx context.Context, forCurrencyCode string) []*dto.CurrencyChangeInfo {
	currencyItems := make([]*dto.CurrencyChangeInfo, 0, len(currencyCodeToName)-1)
	for code, name := range currencyCodeToName {
		if code != forCurrencyCode {
			currencyItems = append(currencyItems, &dto.CurrencyChangeInfo{
				Code:           code,
				Name:           name,
				CurrentPrice:   float64(rand.Intn(100) + 20),
				DayChangePct:   math.Ceil(rand.Float64()+float64(rand.Intn(5))*100) / 100,
				MonthChangePct: math.Ceil(rand.Float64()+float64(rand.Intn(5))*100) / 100,
			})
		}
	}
	//findCurrentPrices(ctx, forCurrencyCode)
	return currencyItems
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

type rate struct {
	changePct float64
}

func findCurrentPrices(ctx context.Context, forCurrencyCode string) {
	//dateNow := time.Now().Format("2006-01-02")
	//dayBefore := time.Now().Sub(time.Date(0, 0, 1, 0, 0, 0, 0, nil))
	//monthBefore := time.Now().Sub(time.Date(0, 1, 0, 0, 0, 0, 0, nil))
	//
	//logger.Error(ctx, forCurrencyCode, dateNow, dayBefore, monthBefore)
	//url := fmt.Sprintf(
	//	"https://api.apilayer.com/exchangerates_data/fluctuation?start_date=%s&end_date=%s",
	//		)
	//
	//client := &http.Client{}
	//req, err := http.NewRequest("GET", url, nil)
	//req.Header.Set("apikey", "zyAc783IAR5xxXdgxohNEaoCMQdH8oFa")
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//res, err := client.Do(req)
	//if res.Body != nil {
	//	defer res.Body.Close()
	//}
	//body, err := ioutil.ReadAll(res.Body)
	//
	//fmt.Println(string(body))
}

var currencyCodeToName = map[string]string{
	"USD": "Доллар США",
	"EUR": "Евро",
	"ILS": "Новый израильский шекель",
	"INR": "Индийская рупия",
	"JPY": "Йена",
	"RUB": "Российский рубль",
	"GBP": "Фунт Стерлингов",
	"KZT": "Тенге",
	"UAH": "Гривна",
	"EGP": "Египетский фунт",
	"CNY": "Китайский юань",
}
