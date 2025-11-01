package main

import (
	"emailn/internal/infra/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	println("Started worker")
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.NewDb()
	repository := database.CampaignRepository{Db: db}
	campaigns, _ := repository.GetCampaignsToBeSent()

	if err != nil {
		println(err.Error())
	}

	println(len(campaigns))

	for _, campaign := range campaigns {
		println(campaign.ID)
	}

}
