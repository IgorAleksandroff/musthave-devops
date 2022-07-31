package repository

import (
	"log"
	"time"
)

func (r rep) MemSync() {
	go func() {
		ticker := time.NewTicker(r.cfg.StoreInterval)
		if r.cfg.StoreInterval == 0 {
			ticker.Stop()
		}
		defer ticker.Stop()
		for {
			select {
			case <-r.ctx.Done():
				err := r.FlushMemo()
				if err != nil {
					log.Printf("can't save metrics, %s", err.Error())
				}
				return
			case <-ticker.C:
				err := r.FlushMemo()
				if err != nil {
					log.Printf("can't save metrics, %s", err.Error())
				}
			}
		}
	}()
}
