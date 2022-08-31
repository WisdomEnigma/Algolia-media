# Algolia-media


Algolia media is a wrapper library which is used to store media files. Algolia provide link to store data on data base, which is good methodology , however I want to store media files in decode form. Later programmer will encode the data and used in applicaton.


# Example

    file, err := os.OpenFile("metamask.png", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}

	wrapper := algolia_wrapper.NewAlgolia_Object(file)

	client := wrapper.ToConnectAlgolia(&protos.Credentials{
		APP_Code:    "...",
		Algolia_AMI: "...",
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

    <!-- Result set return meta data , decode return media object -->
	resultSet, decode , err := wrapper.Get(index, object)
	if err != nil {
		log.Println("Error getting image:", err)
		return
	}

	log.Println("Metadata :", resultSet, "Decode:", decode[0:20])