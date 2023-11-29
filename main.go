package main

import (
	"aci-adr-go-base/service/dal"
	"os"
)

func main() {
	dal.InitMongo(os.Getenv("MONGODB_URI"))
}
