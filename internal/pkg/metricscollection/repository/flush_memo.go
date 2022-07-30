package repository

import (
	"encoding/json"
	"os"

	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/entity"
)

func (r rep) FlushMemo(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(r.metricDB); err != nil {
		return err
	}
	r.metricDB = make(map[string]entity.Metrics)
	return nil
}
