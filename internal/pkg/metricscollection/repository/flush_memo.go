package repository

import (
	"encoding/json"
	"os"
)

func (r rep) FlushMemo() error {
	file, err := os.OpenFile(r.cfg.StorePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(r.metricDB); err != nil {
		return err
	}

	return nil
}
