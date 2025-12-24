package util

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

func QrCodeGen(token string, fileName string) error {
	filepath := filepath.Join("upload", fileName)

	err := qrcode.WriteFile(token, qrcode.Medium, 250, filepath)
	return err
}

func DeleteQrImage(fileName string) error {
	filePath := filepath.Join("upload", fileName)
	if removeErr := os.Remove(filePath); removeErr != nil {
		log.Println("failed to remove qrCode file:", removeErr)
		return fiber.ErrInternalServerError
	}

	return nil
}
