package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"crypto/sha1"
)

func RsaReadPEMPublicKey(pem_str string) (string, string) {
	//pubkey := "\n-----BEGIN RSA PUBLIC KEY-----\nMIGJAoGBANFkG0SbB9bzdLZwDhdACsbLKfz7V7snrN5utn2Ms64iKfgwosNWhxEU\nkbWJeYi8kQC1hLjv8WAZ9hyWN6wMe4ChwNOggKn4XLrx46IHsueosucAcbDrch1h\n87yatm3HwtQEZrBPLZYi2fg8jOHjc4Obf9Du2YV8NFDWgCLinyLFAgMBAAE=\n-----END RSA PUBLIC KEY-----"

	block, rest := pem.Decode([]byte(pem_str))
	if string(rest) != "" {
		return "", "decode failed"
	}
	fmt.Println(len(block.Bytes))
	fmt.Println(block.Headers)

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", "convert failed"+err.Error()
	}
	//switch pub := pub.(type) {
	//case *rsa.PublicKey:
	//	fmt.Println("pub is of type RSA:", pub)
	//case *dsa.PublicKey:
	//	fmt.Println("pub is of type DSA:", pub)
	//case *ecdsa.PublicKey:
	//	fmt.Println("pub is of type ECDSA:", pub)
	//default:
	//	panic("unknown type of public key")
	//}
	der_pub, err := x509.MarshalPKIXPublicKey(pub)

	return hex.EncodeToString(der_pub), string(rest)
}

func RsaReadPEMPrivateKey(pem_str string) (string, string) {
	block, rest := pem.Decode([]byte(pem_str))
	if block == nil {
		return "", "decode failed"
	}
	//fmt.Println(len(block.Bytes))
	//fmt.Println(block.Headers)

	private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", "convert failed,"+err.Error()
	}

	der_private, err := x509.MarshalPKCS8PrivateKey(private)

	return hex.EncodeToString(der_private), string(rest)
}


func RsaCreatePrivateKeyHex() (string, string, error) {
	bits := 1024
	private_key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	private_data, err := x509.MarshalPKCS8PrivateKey(private_key)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	fmt.Println("priv key:", hex.EncodeToString(private_data))
	public_data := &private_key.PublicKey
	der_pub, err := x509.MarshalPKIXPublicKey(public_data)

	if err != nil {
		return "", "", err
	}
	fmt.Println("public key:", hex.EncodeToString(der_pub))
	return hex.EncodeToString(private_data), hex.EncodeToString(der_pub), nil

}

func RsaCreatePrivateKeyPem(keyPemFile string) (string, string, error) {
	bits := 1024
	private_key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	private_data, err := x509.MarshalPKCS8PrivateKey(private_key)

	fmt.Println("priv key:", hex.EncodeToString(private_data))
	public_data := &private_key.PublicKey
	der_pub,_ := x509.MarshalPKIXPublicKey(public_data)

	block := &pem.Block{
		Type:  "公钥",
		Bytes: der_pub,
	}
	file, err := os.Create(keyPemFile)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	err = pem.Encode(file, block)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	der_pub, _ = x509.MarshalPKIXPublicKey(public_data)
	fmt.Println("public key:", hex.EncodeToString(der_pub))
	return hex.EncodeToString(private_data), hex.EncodeToString(der_pub), nil

}

func RsaCreatePrivateKeyBase64() (string, string, error) {
	bits := 1024
	private_key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	private_data, err := x509.MarshalPKCS8PrivateKey(private_key)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", err
	}
	fmt.Println("priv key:", base64.StdEncoding.EncodeToString(private_data))
	public_data := &private_key.PublicKey
	der_pub, err := x509.MarshalPKIXPublicKey(public_data)
	if err != nil {
		return "", "", err
	}
	fmt.Println("public key:", base64.StdEncoding.EncodeToString(der_pub))
	return base64.StdEncoding.EncodeToString(private_data), base64.StdEncoding.EncodeToString(der_pub), nil

}

//（1）加密：采用sha1算法加密后转base64格式
func RsaEncryptWithSha1Base64(originalData, publicKey string) (string, error) {
	key, _ := base64.StdEncoding.DecodeString(publicKey)
	pubKey, _ := x509.ParsePKIXPublicKey(key)
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(originalData))
	return base64.StdEncoding.EncodeToString(encryptedData), err
}

//（2）解密：对采用sha1算法加密后转base64格式的数据进行解密（私钥PKCS1格式）
func RsaDecryptWithSha1Base64(encryptedData, privateKey string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	key, _ := base64.StdEncoding.DecodeString(privateKey)
	prvKey, _ := x509.ParsePKCS1PrivateKey(key)
	originalData, err := rsa.DecryptPKCS1v15(rand.Reader, prvKey, encryptedDecodeBytes)
	return string(originalData), err
}

//（3）签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func RsaSignWithSha1Hex(data string, prvKey string) (string, error) {
	keyByts, err := hex.DecodeString(prvKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyByts)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}
	h := sha1.New()
	h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)
	//fmt.Println("data:", data)
	//fmt.Println("hash:", hash)
	//fmt.Println("hash:", hex.EncodeToString(hash))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hash[:])

	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := hex.EncodeToString(signature)
	return out, nil
}

//（3）签名：采用sha1算法进行签名并输出为base64格式（私钥PKCS8格式）
func RsaSignWithSha1Base64(data string, prvKey string) (string, error) {
	keyByts, err := base64.StdEncoding.DecodeString(prvKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyByts)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		return "", err
	}
	h := sha1.New()
	h.Write([]byte([]byte(data)))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hash[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return "", err
	}
	out := base64.StdEncoding.EncodeToString(signature)
	return out, nil
}

//（4）验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func RsaVerySignWithSha1Hex(originalData, signData, pubKey string) error {
	sign, err := hex.DecodeString(signData)
	if err != nil {
		return err
	}
	public, _ := hex.DecodeString(pubKey)
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

//（4）验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func RsaVerySignWithSha1Base64(originalData, signData, pubKey string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	public, _ := base64.StdEncoding.DecodeString(pubKey)
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

func RsaConvertPrivToPublic(priv_key string) (string, error) {
	priv_data, err := hex.DecodeString(priv_key)
	if err != nil {
		return "", err
	}
	key, err := x509.ParsePKCS8PrivateKey(priv_data)
	if err != nil {
		return "", err
	}
	public_data := &key.(*rsa.PrivateKey).PublicKey
	der_pub, err := x509.MarshalPKIXPublicKey(public_data)
	if err != nil {
		return "", err
	}
	public_data_hex := hex.EncodeToString(der_pub)
	return public_data_hex, nil
}
