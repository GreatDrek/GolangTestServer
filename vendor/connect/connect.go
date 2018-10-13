package connect

var allConnect map[string]string = make(map[string]string)

func AddUser(name string) {
	allConnect[name] = "test"
}

func ReturnAllConnect() map[string]string {
	return allConnect
}
