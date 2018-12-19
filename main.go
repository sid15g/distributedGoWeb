package distributedGoWeb

func main() {
	s := server{ip:"127.0.0.1", port:int16(8080)}
	s.servAt("/getPartitions");
}