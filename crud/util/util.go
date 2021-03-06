package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// DecodeReqBodyToMap decode request body to map , and get d_version
func DecodeReqBodyToMap(c *gin.Context) (map[string]interface{}, string, error) {
	respBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, "", err
	}
	change := make(map[string]interface{})
	if err := json.Unmarshal(respBytes, &change); err != nil {
		return nil, "", err
	}
	dVersionInterface, ok := change["d_version"]
	if !ok {
		return nil, "", errors.New("d_version is required ")
	}
	dVersion, ok := dVersionInterface.(string)
	if !ok {
		return nil, "", errors.New("d_version invalid ")
	}
	delete(change, "d_version")

	return change, dVersion, nil
}
