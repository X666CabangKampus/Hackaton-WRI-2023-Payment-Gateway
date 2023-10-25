package modules

func LoadConfig() map[string]string {
	mapConfig := make(map[string]string)


	// Local
	mapConfig["rabbitUser"] = "hackathon23"
	mapConfig["rabbitPass"] = "kitaharusbisa"
	mapConfig["rabbitHost"] = "103.175.217.181"
	mapConfig["rabbitPort"] = "5672"
	mapConfig["rabbitVHost"] = "hackathon23"

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

