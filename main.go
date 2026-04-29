package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	URL   string  `json:"url"`
}

type Response struct {
	Content []Product `json:"content"`
}

func main() {
	// Настройка прокси (требование ТЗ)
	proxyAddr := "http://user:pass@ip:port"
	proxyURL, _ := url.Parse(proxyAddr)

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	// Категории для парсинга (условие ТЗ)
	categories := map[string]string{
		"Молоко": "https://www.okeydostavka.ru/api/v1/catalog/products/category/moloko-i-slivki",
		"Хлеб":   "https://www.okeydostavka.ru/api/v1/catalog/products/category/khleb-i-khlebobulochnye-izdeliia",
	}

	for catName, apiURL := range categories {
		fmt.Printf("\nКатегория: %s\n", catName)

		req, _ := http.NewRequest("GET", apiURL, nil)

		// Эмуляция браузера для обхода базовых проверок
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Referer", "https://www.okeydostavka.ru/")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Ошибка запроса: %v\n", err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Проверка на блокировку WAF (возвращает HTML вместо JSON)
		if len(body) > 0 && body[0] != '{' {
			fmt.Printf("Доступ к '%s' ограничен. Требуются резидентские прокси.\n", catName)
			continue
		}

		var data Response
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Printf("Ошибка обработки JSON: %v\n", err)
			continue
		}

		// Вывод первых 5 позиций
		if len(data.Content) == 0 {
			fmt.Println("Товары не найдены.")
		} else {
			for i, p := range data.Content {
				if i >= 5 {
					break
				}
				fmt.Printf("%d. %s | %.2f руб. | %s\n", i+1, p.Name, p.Price, p.URL)
			}
		}

		// Задержка между запросами
		time.Sleep(2 * time.Second)
	}

}
