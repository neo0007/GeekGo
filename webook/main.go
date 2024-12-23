package main

func main() {
	r := InitWebServer()
	//r.Run(":8081")
	err := r.Run(":8081")
	if err != nil {
		panic(err)
	}
}
