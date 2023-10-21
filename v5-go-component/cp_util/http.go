package cp_util

import (
	"time"

	"net/http"

	"io/ioutil"
)

func Do(req *http.Request) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//if resp.StatusCode != http.StatusOK {
	//	return nil, errors.New(fmt.Sprintf("status code:[%d]", resp.StatusCode))
	//}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
