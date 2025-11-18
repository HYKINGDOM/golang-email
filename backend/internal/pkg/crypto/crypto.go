package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"

    "github.com/spf13/viper"
)

func getKey() ([]byte, error) {
    key := viper.GetString("security.secret_key")
    if len(key) == 0 {
        return nil, errors.New("未配置secret_key")
    }
    if len(key) < 32 {
        k := make([]byte, 32)
        copy(k, []byte(key))
        key = string(k)
    }
    b := []byte(key)
    if len(b) > 32 {
        b = b[:32]
    }
    return b, nil
}

func EncryptString(plain string) (string, error) {
    key, err := getKey()
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    ct := gcm.Seal(nonce, nonce, []byte(plain), nil)
    return base64.StdEncoding.EncodeToString(ct), nil
}

func DecryptString(cipherText string) (string, error) {
    key, err := getKey()
    if err != nil {
        return "", err
    }
    data, err := base64.StdEncoding.DecodeString(cipherText)
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", errors.New("密文格式错误")
    }
    nonce := data[:nonceSize]
    ciphertext := data[nonceSize:]
    pt, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }
    return string(pt), nil
}