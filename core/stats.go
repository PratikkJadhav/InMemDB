package core

var keyKeyspaceStat [4]map[string]int

func UpdateDBStat(num int, metric string, value int) {
	keyKeyspaceStat[num][metric] = value
}
