package enviroment

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	ServerEnvServerURL        = "ADDRESS"
	ServerEnvStoreInterval    = "STORE_INTERVAL"
	ServerEnvStoreFile        = "STORE_FILE"
	ServerEnvRestore          = "RESTORE"
	ServerEnvHashKey          = "KEY"
	ServerEnvDB               = "DATABASE_DSN"
	ServerEnvPrivateCryptoKey = "CRYPTO_KEY"
	ServerEnvPublicCfgPath    = "CONFIG"
	ServerEnvSubnet           = "TRUSTED_SUBNET"

	ServerDefaultServerURL     = "localhost:8080"
	ServerDefaultStoreInterval = 300 * time.Second
	ServerDefaultStoreFile     = "/tmp/devops-metrics-db.json"
	ServerDefaultRestore       = true
	ServerDefaultGRPSSocket    = ":3200"
	ServerDefaultString        = ""
)

type ServerConfig struct {
	Host          string
	HashKey       string
	CryptoKeyPath string
	subnet        *net.IPNet
	GRPSSocket    string
}

type config struct {
	StoreInterval time.Duration
	StorePath     string
	Restore       bool
	AddressDB     string
	ServerConfig
}

func NewServerConfig() config {
	log.Printf("os: %+v", os.Args)

	cfg := config{
		ServerConfig: ServerConfig{
			Host:          ServerDefaultServerURL,
			HashKey:       ServerDefaultString,
			CryptoKeyPath: ServerDefaultString,
			GRPSSocket:    ServerDefaultGRPSSocket,
		},
		StoreInterval: ServerDefaultStoreInterval,
		StorePath:     ServerDefaultStoreFile,
		Restore:       ServerDefaultRestore,
		AddressDB:     ServerDefaultString,
	}

	hostFlag := flag.String("a", ServerDefaultServerURL, "адрес и порт сервера")
	storeIntervalFlag := flag.Duration("i", ServerDefaultStoreInterval, "интервал времени в секундах, по истечении которого текущие показания сервера сбрасываются на диск")
	storePathFlag := flag.String("f", ServerDefaultStoreFile, "строка, имя файла, где хранятся значения")
	restoreFlag := flag.Bool("r", ServerDefaultRestore, "булево значение (true/false), определяющее, загружать или нет начальные значения")
	hashKeyFlag := flag.String("k", ServerDefaultString, "ключ подписи метрик")
	addressDBflag := flag.String("d", ServerDefaultString, "адрес подключения к БД")
	cryptoKeyFlag := flag.String("crypto-key", ServerDefaultString, "путь до файла с приватным ключом")
	cfgPathFlag := flag.String("c", ServerDefaultString, "путь до json файла конфигурации сервера")
	subnetFlag := flag.String("t", ServerDefaultString, "доверенная подсеть CIDR")
	socketFlag := flag.String("s", ServerDefaultGRPSSocket, "если не указан порт gRPC сервера, то используется HTTP клиент")

	flag.Parse()

	cfgJSONPath := GetEnvString(ServerEnvPublicCfgPath, *cfgPathFlag)
	if cfgJSONPath != ServerDefaultString {
		updateServerConfigByJSON(cfgJSONPath, &cfg)
	}

	// update Client config by flags
	if hostFlag != nil && isFlagPassed("a") {
		cfg.Host = *hostFlag
	}
	if restoreFlag != nil && isFlagPassed("r") {
		cfg.Restore = *restoreFlag
	}
	if storeIntervalFlag != nil && isFlagPassed("i") {
		cfg.StoreInterval = *storeIntervalFlag
	}
	if storePathFlag != nil && isFlagPassed("f") {
		cfg.StorePath = *storePathFlag
	}
	if addressDBflag != nil && isFlagPassed("d") {
		cfg.AddressDB = *addressDBflag
	}
	if cryptoKeyFlag != nil && isFlagPassed("crypto-key") {
		cfg.CryptoKeyPath = *cryptoKeyFlag
	}
	subnetString := cfg.subnet.String()
	if subnetFlag != nil && isFlagPassed("t") {
		subnetString = *subnetFlag
	}
	if socketFlag != nil && isFlagPassed("s") {
		cfg.GRPSSocket = *socketFlag
	}

	// update Client config by env, default is flag or json parameter
	cfg.Host = GetEnvString(ServerEnvServerURL, cfg.Host)
	cfg.HashKey = GetEnvString(ServerEnvHashKey, *hashKeyFlag)
	cfg.CryptoKeyPath = GetEnvString(ServerEnvPrivateCryptoKey, cfg.CryptoKeyPath)
	cfg.StoreInterval = GetEnvDuration(ServerEnvStoreInterval, cfg.StoreInterval)
	cfg.StorePath = GetEnvString(ServerEnvStoreFile, cfg.StorePath)
	cfg.Restore = GetEnvBool(ServerEnvRestore, cfg.Restore)
	cfg.AddressDB = GetEnvString(ServerEnvDB, cfg.AddressDB)
	subnetString = GetEnvString(ServerEnvSubnet, subnetString)

	if _, v, err := net.ParseCIDR(subnetString); err != nil {
		log.Printf(errorParseServerJSONConfig, subnetString, err)
	} else {
		cfg.subnet = v
	}

	log.Printf("Parsed Server config: %+v", cfg)

	return cfg
}

func (s ServerConfig) GetTrustedIPMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s.subnet != nil {
				clientIPString := r.Header.Get("X-Real-IP")
				clientIP := net.ParseIP(clientIPString)

				if !s.subnet.Contains(clientIP) {
					log.Println("request from an untrusted address:", clientIPString)
					http.Error(w, "ip isn't part of a trusted subnet", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (s ServerConfig) GetTrustedIPInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if s.subnet != nil {
			var clientIP net.IP
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				values := md.Get("X-Real-IP")
				if len(values) > 0 {
					clientIP = net.ParseIP(values[0])
				}
			}
			if !s.subnet.Contains(clientIP) {
				log.Println("request from an untrusted address:", clientIP)
				return nil, status.Error(codes.Unauthenticated, "ip isn't part of a trusted subnet")
			}
		}

		return handler(ctx, req)
	}
}
