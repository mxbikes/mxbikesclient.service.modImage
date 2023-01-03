package models

import (
	protobuffer "github.com/mxbikes/protobuf/modImage"
)

type ModImage struct {
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
	Url    string `json:"url"`
}

func ModImageToProto(modImage *ModImage) *protobuffer.ModImage {
	return &protobuffer.ModImage{
		Name:   modImage.Name,
		Bucket: modImage.Bucket,
		Url:    modImage.Url,
	}
}

func ModImagesToProto(modImages []*ModImage) []*protobuffer.ModImage {
	orders := make([]*protobuffer.ModImage, 0, len(modImages))
	for _, projection := range modImages {
		orders = append(orders, ModImageToProto(projection))
	}
	return orders
}
