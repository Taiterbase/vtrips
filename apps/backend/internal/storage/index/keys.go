package index

func MakePK(clientID string) string {
	return "client_id#" + clientID
}

func MakeSK(tripID string) string {
	return "trip_id#" + tripID
}
