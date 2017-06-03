package encrypt

import (
	"crypto/des"
	"errors"
)

/************************************/
/********** ESC encryption **********/
/************************************/

// DES encrypt
// Note:
//	#秘钥 长度必须为24byte,否则直接返回错误
//	PHP中只要8byte以上即可
//	Java中必须是24byte以上，内部会取前24byte(相当于就是24byte)
//	key := []byte{0xD5, 0x92, 0x86, 0x02, 0x2A, 0x0B, 0x3E, 0x64}
// Use:
//	enData := []byte("hello") #加密字符串
//	out, _ := encrypt.DesEncrypt(enData, key) #调用加密
//	data["encrypt"] = out
//
//	out, _ = encrypt.DesDecrypt(out, key) #调用解密
//	data["decrypt"] = string(out)
func DesEncrypt(data, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

// DES decrypt
func DesDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return out, nil
}
