package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
)

func encryptData(plainText, key string) (string, error) {
    block, err := aes.NewCipher([]byte(key))
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

    cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
    return base64.StdEncoding.EncodeToString(cipherText), nil
}

func decryptData(encryptedText, key string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(encryptedText)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return "", errors.New("invalid ciphertext")
    }

    nonce, cipherText := data[:nonceSize], data[nonceSize:]
    plainText, err := gcm.Open(nil, nonce, cipherText, nil)
    if err != nil {
        return "", err
    }

    return string(plainText), nil
}