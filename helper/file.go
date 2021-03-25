package helper

import (
	"encoding/json"
	"os"
	"twfinder/logger"
)

// SaveReplaceJsonFile : save any type to json file
func SaveReplaceJsonFile(data interface{}, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		logger.Errorf("Error occurred during open file %v %v", path, err)
		return err
	}
	defer f.Close()
	jsonToSave, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		logger.Errorf("Error occurred during Marshal json %v", err)
		return err
	}
	_, err = f.Write(jsonToSave)
	if err != nil {
		logger.Errorf("Error occurred during Write to json file %v", err)
		return err
	}
	return nil
}
