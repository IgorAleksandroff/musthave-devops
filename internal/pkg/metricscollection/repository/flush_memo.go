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

	//var restoredMetrics map[string]entity.Metrics
	//decoder := json.NewDecoder(file)
	//if err = decoder.Decode(&restoredMetrics); err != nil {
	//	return err
	//}
	//log.Println(restoredMetrics)
	//
	//for name, metric := range r.metricDB {
	//	restoredMetrics[name] = metric
	//}

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(r.metricDB); err != nil {
		return err
	}

	//r.metricDB = make(map[string]entity.Metrics)
	return nil
}
