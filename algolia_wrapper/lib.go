package algolia_wrapper

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"reflect"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	algolia_Client "github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/wisdomenigma/algolia-media/algolia_wrapper/protos"
)

type MediaFormat struct {
	File       *os.File               `json:"File"`
	Class      protos.ImageType       `json:"Class"`
	Avatar     *image.Image           `json:"Avatar"`
	FormatType protos.ImageFormatType `json:"FormatType"`
}

type Algolia_File_Service interface {
	ToConnectAlgolia(credentials *protos.Credentials) *algolia_Client.Client
	Index(client *algolia_Client.Client, name ...string) *algolia_Client.Index
	Put(index *algolia_Client.Index) (*algolia_Client.SaveObjectRes, error)
	Get(index *algolia_Client.Index, object *algolia_Client.SaveObjectRes) (interface{}, image.Image, error)
	Decode() (image.Image, error)
	Close() error
	Query(index *algolia_Client.Index, result *algolia_Client.SaveObjectRes, doc ...string) (algolia_Client.QueryRes, error)
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

	if strings.Contains(format.File.Name(), ".png") {

		format.FormatType = protos.ImageFormatType_PNG
	} else if strings.Contains(format.File.Name(), ".jpeg") {

		format.FormatType = protos.ImageFormatType_JPEG
	} else if strings.Contains(format.File.Name(), ".gif") {

		format.FormatType = protos.ImageFormatType_GIF
	}

	img, err := format.Decode()
	if err != nil {
		return &algolia_Client.SaveObjectRes{}, err
	}

	format.Avatar = &img

	result, err := index.SaveObject(*format.Avatar)
	if err != nil {
		return &algolia_Client.SaveObjectRes{}, err
	}

	if reflect.DeepEqual(result.ObjectID, " ") {
		return &algolia_Client.SaveObjectRes{}, err
	}

	return &result, nil
}

func (format *MediaFormat) Get(index *algolia_Client.Index, object *algolia_Client.SaveObjectRes) (interface{}, image.Image, error) {

	var inter interface{} = nil
	if strings.Contains(object.ObjectID, " ") {
		return inter, &image.RGBA64{}, errors.New("object id must not be empty")
	}

	obj, err := format.Query(index, object, []string{"MyAvatar"}...)
	if err != nil && !reflect.DeepEqual(obj, object.ObjectID) {
		return inter, &image.RGBA64{}, errors.New("object id is not match")
	}

	return [...]interface{}{
		format.Class,
		format.FormatType,
	}, *format.Avatar, nil
}

func (format *MediaFormat) Decode() (image.Image, error) {

	var img image.Image
	var err error
	if reflect.DeepEqual(format.FormatType, protos.ImageFormatType_PNG) {

		img, err = png.Decode(format.File)
		if err != nil {
			return &image.RGBA{
				Pix:    []uint8{},
				Stride: 0,
				Rect:   image.Rectangle{},
			}, err
		}
	} else if reflect.DeepEqual(format.FormatType, protos.ImageFormatType_JPEG) {

		img, err = jpeg.Decode(format.File)
		if err != nil {
			return &image.RGBA64{
				Pix:    []uint8{},
				Stride: 0,
				Rect:   image.Rectangle{},
			}, err
		}
	} else if reflect.DeepEqual(format.FormatType, protos.ImageFormatType_GIF) {

		img, err = gif.Decode(format.File)
		if err != nil {
			return &image.RGBA{
				Pix:    []uint8{},
				Stride: 0,
				Rect:   image.Rectangle{},
			}, err
		}
	} else {
		return &image.RGBA{Pix: []uint8{},
			Stride: 0,
			Rect:   image.Rectangle{}}, errors.New("image format not supported")
	}

	return img, nil
}

func (format *MediaFormat) Close() error {
	return format.File.Close()
}

func (format *MediaFormat) Query(index *algolia_Client.Index, result *algolia_Client.SaveObjectRes, doc ...string) (algolia_Client.QueryRes, error) {

	if strings.Contains(doc[0], " ") {
		return algolia_Client.QueryRes{}, errors.New("document name should not be empty")
	}

	params := []interface{}{
		opt.AttributesToRetrieve("objectID"),
		opt.HitsPerPage(10),
	}

	return index.Search(result.ObjectID, params...)
}
