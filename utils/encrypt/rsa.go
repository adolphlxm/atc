package encrypt

import (
	"encoding/pem"
	"errors"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
)

/************************************/
/********** AES encryption **********/
/************************************/
//
// 基于PKCS#1规范
// PublicKey和PrivateKey两个类型分别代表公钥和私钥
// OpenSSl生成publickey,privatekey
//	$ openssl
//	$ genrsa -out rsa_private_key.pem 2048 // 生成私钥
//	$ pkcs8 -topk8 -inform PEM -in rsa_private_key.pem -outform PEM –nocrypt // 第二句命令：把RSA私钥转换成PKCS8格式。提示输入密码，密码为空（直接回车）就行；
//	$ rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem // 生成公钥
func RsaEncrypt(origData, publicKey []byte) ([]byte, error){
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error.")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}


func RsaDecrypt(ciphertext, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil  {
		return nil, errors.New("private key error.")
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, private, ciphertext)
}
