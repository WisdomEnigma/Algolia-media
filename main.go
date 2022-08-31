package main

import (
	"log"
	"os"

	"github.com/wisdomenigma/algolia-media/algolia_wrapper"
	"github.com/wisdomenigma/algolia-media/algolia_wrapper/protos"
)

func main() {

	file, err := os.OpenFile("metamask.png", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}

	wrapper := algolia_wrapper.NewAlgolia_Object(file)

	client := wrapper.ToConnectAlgolia(&protos.Credentials{
		APP_Code:    "XNZCWV197B",
		Algolia_AMI: "37695cff1374a0e0f2929012fd65b5ef",
	})

	index := wrapper.Index(client, []string{"MyAvatar"}...)
	object, err := wrapper.Put(index)
	if err != nil {
		log.Println("Error putting image:", err)
		return
	}

	result, err := wrapper.Query(index, object, []string{"MyAvatar"}...)
	if err != nil {
		log.Println("Error finding image :", err)
		return
	}

	log.Println("Result:", result.Query)

	resultSet, _, err := wrapper.Get(index, object)
	if err != nil {
		log.Println("Error getting image:", err)
		return
	}

	log.Println("Result :", resultSet)

	if err = wrapper.Close(); err != nil {
		return
	}

}
