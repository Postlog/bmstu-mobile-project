package image

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	extensionPNG = "png"
)

type Repository struct {
	folderPath string
}

func New(folderPath string) (*Repository, error) {
	if !filepath.IsAbs(folderPath) {
		return nil, fmt.Errorf("provided path must be absolute")
	}

	return &Repository{
		folderPath: folderPath,
	}, nil
}

func (r Repository) Count(_ context.Context) (int, error) {
	files, err := os.ReadDir(r.folderPath)
	if err != nil {
		return 0, fmt.Errorf("error reading dir: %w", err)
	}

	count := 0
	for _, file := range files {
		name := file.Name()
		ext := filepath.Ext(name)
		if strings.TrimPrefix(ext, ".") == extensionPNG {
			count++
		}
	}

	return count, nil
}

func (r Repository) Save(_ context.Context, pngImageBytes []byte) (uuid.UUID, error) {
	config, err := png.DecodeConfig(bytes.NewReader(pngImageBytes))
	if err != nil {
		var formatErr png.FormatError
		if errors.Is(err, io.ErrUnexpectedEOF) || errors.As(err, &formatErr) {
			return uuid.Nil, ErrImageNotInPNGFormat
		}

		return uuid.Nil, fmt.Errorf("read image header: %w", err)
	}

	if config.Width*config.Height > MaxImagePixelsCount {
		return uuid.Nil, ErrImageTooLarge
	}

	nameUUID := uuid.New()
	fileName := fmt.Sprintf("%s.%s", nameUUID.String(), extensionPNG)
	filePath := filepath.Join(r.folderPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create file \"%s\": %w", filePath, err)
	}

	bytesWrote, err := file.Write(pngImageBytes)
	if err != nil {
		return uuid.Nil, fmt.Errorf("write image bytes to file: %w", err)
	}

	if bytesWrote < len(pngImageBytes) {
		return uuid.Nil, errors.New("bytes wrote is less than image bytes")
	}

	return nameUUID, nil
}

func (r Repository) Get(_ context.Context, id uuid.UUID) (Image, error) {
	fileName := fmt.Sprintf("%s.%s", id.String(), extensionPNG)
	filePath := filepath.Join(r.folderPath, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Image{}, ErrImageNotExist
		}

		return Image{}, fmt.Errorf("open image \"%s\": %w", filePath, err)
	}

	config, err := png.DecodeConfig(file)
	if err != nil {
		return Image{}, fmt.Errorf("read image header: %w", err)
	}

	_, _ = file.Seek(0, 0)

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return Image{}, fmt.Errorf("read image: %w", err)
	}

	return Image{
		Bytes:  imageBytes,
		Width:  config.Width,
		Height: config.Height,
	}, nil
}
