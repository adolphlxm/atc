package encrypt

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestRsa(t *testing.T) {
	var publicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCgohnj06UmBTpWmDsbNUyVoFOy
x8c/+rB6nhDGE9BvYTFEwUKabhpp4T0fB2u84yWLsGhzX8Oa1tlFv8mjBRC/r8Dg
IBAYF3XxAIJCn+P25HmYZnJf/gcQkxsR2qP0SzilBV0FggsHTRgbRIr17bJFlrYK
yviKpYYVfCtPfc9uEwIDAQAB
-----END PUBLIC KEY-----`)

	var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCgohnj06UmBTpWmDsbNUyVoFOyx8c/+rB6nhDGE9BvYTFEwUKa
bhpp4T0fB2u84yWLsGhzX8Oa1tlFv8mjBRC/r8DgIBAYF3XxAIJCn+P25HmYZnJf
/gcQkxsR2qP0SzilBV0FggsHTRgbRIr17bJFlrYKyviKpYYVfCtPfc9uEwIDAQAB
AoGAT096BPjso5cVNMg49rgcNgxRBwaIpFg8Zm5Azevej/+m+dGJ9wHB3OEu8nER
itOEKDJTgBOaWosV/9Av2GOZ3qK/CjcuqGUF3A+D0XS97mNwjtFfPTLxd53zsYYE
baYhyrZFWDi7ynsoNWr3YTRtB4/KkPJih17djpHCziguRxECQQDNOJxYRVisWDHn
pTVuyB5K1QA7em42hr4unqkrsGnVsTzAHXm5Di31HUDU36jLgluEUtmyn6O9SYT9
fEbrBVn3AkEAyGEki0s3xUH493pbHZaLW3uYU54IpF50y9jd4vIuxrylG/9NA8Uc
YM28xk853NlYz3SLMw+EVbVcjaxL+LOlxQJAS6XYi/lUDIOeMcOGhMWj1PXbVhF1
WwgkRs8ZkQ9AlBL3T+INopeFfVtBMLcZY5sz3P0lXmDWXMojCcWr5qpcVQJAXbDm
OGckHYx6T6SbQ9tnL5A7qiVDXy93JvUw0nNwoaYFAXE+3ltkqHKqKINUx8msd9vD
Vk2UD8ssCmYcY54EDQJAeOxoNVYi+zjpLHveQSIK/qrOi5yDmK/Seqy4FsfZDBaN
EiC7iaRlDHqdW/P0N4C9TkGli0uwrwOBdYBDGfe2AA==
-----END RSA PRIVATE KEY-----`)
	origData := []byte("Adolph.liu for atc")
	ciphertext, err := RsaEncrypt(origData, publicKey)
	if err != nil {
		t.Errorf("RsaEncrypt err:%v", err.Error())
	}
	fmt.Println("Rsa encrypt base64 :", base64.StdEncoding.EncodeToString(ciphertext))

	data, err := RsaDecrypt(ciphertext, privateKey)
	if err != nil {
		t.Errorf("RsaDecrypt err:%v", err.Error())
	}
	fmt.Println("Rsa decrypt :", string(data))
}
