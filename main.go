package main

import "APIGenerator/rooter"

func main() {
	r := rooter.SetupRouter()
	r.Run(":8000")
}
