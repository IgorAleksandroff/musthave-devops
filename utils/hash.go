package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/IgorAleksandroff/musthave-devops/utils/enviroment/clientconfig"
)

func GetHash(value, key string) string {
	if key == clientconfig.DefaultEnvHashKey {
		return ""
	}
	// подписываем алгоритмом HMAC, используя SHA256
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	dst := h.Sum(nil)

	fmt.Printf("%x", dst)
	return fmt.Sprintf("%x", dst)
}
