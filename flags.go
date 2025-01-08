package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"os"
)

var certDomain = flag.String("domain", "", "Domain name to use for renewal")
var tencentSecretId = flag.String("secretid", "", "Tencent cloud secret id")
var tencentSecretKey = flag.String("secretkey", "", "Tencent cloud secret key")
var dryRun = flag.Bool("dryrun", false, "Don't really do anything")

func init() {
	flag.Parse()
}

func envValue(key, fallback string) string {
	if len(fallback) > 0 {
		return fallback
	}

	return os.Getenv(key)
}

func getCertDomain() string {
	return envValue("Le_Domain", *certDomain)
}

func getTencentClient() *common.Credential {
	credential := common.NewCredential(
		envValue("TENCENT_SECRET_ID", *tencentSecretId),
		envValue("TENCENT_SECRET_KEY", *tencentSecretKey))

	return credential
}

func readCertKey() string {
	key, err := os.ReadFile(os.Getenv("CERT_KEY_PATH"))
	if err != nil {
		logrus.Fatal("Error when read key file on %s: %v", os.Getenv("CERT_KEY_PATH"), err)
	}

	return string(key)
}

func readCertFullchain() string {
	data, err := os.ReadFile(os.Getenv("CERT_FULLCHAIN_PATH"))
	if err != nil {
		logrus.Fatal("Error when read cert file on %s: %v", os.Getenv("CERT_FULLCHAIN_PATH"), err)
	}

	return string(data)
}
