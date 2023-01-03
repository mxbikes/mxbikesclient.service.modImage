package handler

import (
	"context"

	"github.com/gogo/status"
	"github.com/google/uuid"
	"github.com/mxbikes/mxbikesclient.service.modImage/models"
	"github.com/mxbikes/mxbikesclient.service.modImage/repository"
	protobuffer "github.com/mxbikes/protobuf/modImage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

type ModImage struct {
	protobuffer.UnimplementedModImageServiceServer
	minio  repository.ModImageMinioRepository
	logger logrus.Logger
}

// Return a new handler
func New(minio repository.ModImageMinioRepository, logger logrus.Logger) *ModImage {
	return &ModImage{minio: minio, logger: logger}
}

func (e *ModImage) GetModImagesByModID(ctx context.Context, req *protobuffer.GetModImagesByModIDRequest) (*protobuffer.GetModImagesByModIDResponse, error) {
	// Check if valid uuid
	modID, err := uuid.Parse(req.ModID)
	if err != nil {
		e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.ModImage_GetModImagesByModID"}).Errorf("request ModID is not a valid UUID: {%s}", req.ModID)
		return nil, status.Error(codes.Internal, "Error request value ModID, is not a valid UUID!")
	}

	modImages, err := e.minio.GetModImagesByModID(ctx, modID.String())
	if err != nil {
		return nil, err
	}

	e.logger.WithFields(logrus.Fields{"prefix": "SERVICE.ModImage_GetModImagesByModID"}).Infof("mod images fetched with modID: {%s} ", req.ModID)

	return &protobuffer.GetModImagesByModIDResponse{ModImage: models.ModImagesToProto(modImages)}, nil
}
