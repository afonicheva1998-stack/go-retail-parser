package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Модель данных продукта согласно ТЗ
type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	URL   string  `json:"url"`
}

type Response struct {
	Content []Product `json:"content"`
}

func main() {
	// Поддержка прокси для обхода блокировок (Условие ТЗ)
	proxyAddr := "http://user:pass@ip:port"
	proxyURL, _ := url.Parse(proxyAddr)

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	// Работа через API (не "в лоб"), категория "Молоко"
	apiURL := "https://www.okeydostavka.ru/api/v1/catalog/products/category/moloko-i-slivki"
	req, _ := http.NewRequest("GET", apiURL, nil)

	// Эмуляция браузера
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Обработка защиты WAF/Qrator
	if len(body) > 0 && body[0] != '{' {
		fmt.Println("Доступ ограничен защитой сайта. Требуются резидентские прокси.")
		return
	}

	var data Response
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Ошибка парсинга: %v\n", err)
		return
	}

	// Вывод результатов
	for i, p := range data.Content {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s | %.2f руб.\n", i+1, p.Name, p.Price)
	}
}
