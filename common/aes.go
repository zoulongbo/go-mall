package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//aes加密 高级加密标准 advanced encryption standard
//16,24,32为字符串 分别对应 AES-128 ,AES-192 , AES-256
var AesPwdKey = []byte("AES@---xxPwd0001")


//PKCS7 填充模式
func AesPKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(cipherText, padText...)
}
//填充的反向操作，删除填充字符串
func AesPKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unPadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unPadding)], nil
	}
}

//实现加密
func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = AesPKCS7Padding(origData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypt := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(encrypt, origData)
	return encrypt, nil
}

//实现解密
func AesDeCrypt(decrypt []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(decrypt))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, decrypt)
	//去除填充字符串
	origData, err = AesPKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

//加密base64
func AesEnPwdCode(pwd []byte) (string, error) {
	result, err := AesEncrypt(pwd, AesPwdKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

//解密
func AesDePwdCode(pwd string) ([]byte, error) {
	//解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil, err
	}
	//执行AES解密
	return AesDeCrypt(pwdByte, AesPwdKey)

}



