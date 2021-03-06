package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type RsaServer struct {
	publicKey  []byte // 公钥
	privateKey []byte // 私钥
}

func NewRsaServer(publickKeyFile, privateKeyFile string) (*RsaServer, error) {
	publicFile, err := os.Open(publickKeyFile)
	defer publicFile.Close()
	if err != nil {
		return nil, err
	}
	info, _ := publicFile.Stat()
	publicByte := make([]byte, info.Size())
	_, err = publicFile.Read(publicByte)
	if err != nil {
		return nil, err
	}

	privateFile, err := os.Open(privateKeyFile)

	defer privateFile.Close()
	if err != nil {
		return nil, err
	}
	info, _ = privateFile.Stat()
	privateByte := make([]byte, info.Size())
	_, err = privateFile.Read(privateByte)
	if err != nil {
		return nil, err
	}
	return &RsaServer{publicKey: publicByte, privateKey: privateByte}, nil
}

func (r *RsaServer) Encrypt(content string) (encryptStr string, err error) {
	block, _ := pem.Decode(r.publicKey)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	encryptByte, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(content)) //RSA算法加密
	if err != nil {
		return "", err
	}
	if len(encryptByte) == 0 {
		return "", fmt.Errorf("encrpt failed")
	}
	encryptStr = string(encryptByte)
	return
}

func (r *RsaServer) Decryption(cipherText string) (decryptStr string, err error) {
	block, _ := pem.Decode(r.privateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(cipherText))

	if err != nil {
		return "", err
	}
	return string(plainText), nil
}
