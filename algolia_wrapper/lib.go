package algolia_wrapper

import (
	"os"
	"reflect"

	algolia_Client "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/wisdomenigma/algolia-media/algolia_wrapper/protos"
)

type MediaFormat struct {
	File  *os.File         `json:"File"`
	Class protos.ImageType `json:"Class"`
}

type Algolia_File_Service interface {
	ToConnectAlgolia(credentials *protos.Credentials) *algolia_Client.Client
	Index(client *algolia_Client.Client, name ...string) *algolia_Client.Index
	Put(index *algolia_Client.Index) (*algolia_Client.SaveObjectRes, error)
	Get(index *algolia_Client.Index, objectId ...string) error
}

func NewAlgolia_Object(file *os.File) Algolia_File_Service {
	return &MediaFormat{File: file, Class: protos.ImageType_FILE}
}

func (format *MediaFormat) ToConnectAlgolia(credentials *protos.Credentials) *algolia_Client.Client {

	return algolia_Client.NewClient(credentials.APP_Code, credentials.Algolia_AMI)
}

func (format *MediaFormat) Index(client *algolia_Client.Client, name ...string) *algolia_Client.Index {
	return client.InitIndex(name[0])
}

func (format *MediaFormat) Put(index *algolia_Client.Index) (*algolia_Client.SaveObjectRes, error) {

	if ok, err := index.Exists(); !ok && err != nil {
		return &algolia_Client.SaveObjectRes{}, err
	}

	result, err := index.SaveObject(format.File)
	if err != nil {
		return &algolia_Client.SaveObjectRes{}, err
	}

	if reflect.DeepEqual(result.ObjectID, " ") {
		return &algolia_Client.SaveObjectRes{}, err
	}

	return &result, nil
}

func (format *MediaFormat) Get(index *algolia_Client.Index, objectId ...string) error {
	return index.GetObject(objectId[0], format.File)
}
