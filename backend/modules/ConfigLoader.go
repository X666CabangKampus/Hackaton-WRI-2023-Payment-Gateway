package modules

import "os"

func LoadConfig() map[string]string {
	mapConfig := make(map[string]string)

	// Local
	mapConfig["rabbitUser"] = os.Getenv("RABBIT_USERNAME")
	mapConfig["rabbitPass"] = os.Getenv("RABBIT_PASSWORD")
	mapConfig["rabbitHost"] = os.Getenv("RABBIT_HOST")
	mapConfig["rabbitPort"] = os.Getenv("RABBIT_PORT")
	mapConfig["rabbitVHost"] = os.Getenv("RABBIT_VHOST")

	return mapConfig
}

func LoadConfigProduction() map[string]string {
	mapConfig := make(map[string]string)

	// Production
	mapConfig["mongoDBHost"] = "172.31.2.44"
	mapConfig["mongoDBPort"] = "27017"
	mapConfig["logEndpoint"] = "http://localhost:55555/log"

	return mapConfig
}
