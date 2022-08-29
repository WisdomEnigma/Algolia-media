package main

import (
	"fmt"
	"os"

	"github.com/wisdomenigma/algolia-media/algolia_wrapper"
	"github.com/wisdomenigma/algolia-media/algolia_wrapper/protos"
)

func main() {

	file, err := os.OpenFile("metamask.png", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}

	if err := file.Close(); err != nil {
		return
	}

	wrapper := algolia_wrapper.NewAlgolia_Object(file)

	client := wrapper.ToConnectAlgolia(&protos.Credentials{
		APP_Code:    "XNZCWV197B",
		Algolia_AMI: "37695cff1374a0e0f2929012fd65b5ef",
	})

	index := wrapper.Index(client, []string{"MyAvatar"}...)
	object, _ := wrapper.Put(index)

	fmt.Println("Success : ", object)
	// if err = wrapper.Get(index, object.ObjectID); err != nil {
	// 	log.Fatalln("Object is not found")
	// }
}
