package settings

import (
	"encoding/base64"
	"log"
)

func Base64EncodeString(data []byte) string {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)

	return string(dst)
}

func DecodeString(val string) (string, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(val)))
	n, err := base64.StdEncoding.Decode(dst, []byte(val))
	if err != nil {
		log.Println("Error decodign string")

		return "", err
	}
	return string(dst[:n]), nil
}
