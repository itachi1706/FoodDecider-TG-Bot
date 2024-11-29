package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RandomOrgResponse struct {
	Random struct {
		Data []int `json:"data"`
	} `json:"random"`
}

func GetTrueRandomNumber(min, max int64) (int, error) {
	url := fmt.Sprintf("https://www.random.org/integers/?num=1&min=%d&max=%d&col=1&base=10&format=plain&rnd=new", min, max)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var result int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result, nil
}
