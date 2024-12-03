package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RandomOrgResponse struct {
	Random struct {
		Data []int `json:"data"`
	} `json:"random"`
}

func GetTrueRandomNumber(min, max int64) (int, error) {
	if (min < 0) || (max < 0) {
		log.Println("min and max must be positive")
		return 0, fmt.Errorf("min and max must be positive")
	}

	if min > max {
		log.Println("min must be less than or equal to max")
		return 0, fmt.Errorf("min must be less than or equal to max")
	}

	if (min == 0) && (max == 0) {
		log.Println("min and max must not be both 0")
		return 0, fmt.Errorf("min and max must not be both 0")
	}

	if min == max {
		return int(min), nil
	}

	url := fmt.Sprintf("https://www.random.org/integers/?num=1&min=%d&max=%d&col=1&base=10&format=plain&rnd=new", min, max)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error getting random number: ", err)
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var result int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("error decoding random number: %w", err)
		return 0, err
	}
	return result, nil
}
