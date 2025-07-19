package config

import "os"

// Tidak perlu client SDK, hanya environment variables
func GetImageKitKeys() (publicKey, privateKey, urlEndpoint string) {
    return os.Getenv("IMAGEKIT_PUBLIC_KEY"),
        os.Getenv("IMAGEKIT_PRIVATE_KEY"),
        os.Getenv("IMAGEKIT_URL_ENDPOINT")
}