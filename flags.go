package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

var certDomain = flag.String("domain", "", "Domain name to use for renewal")
var tencentSecretId = flag.String("secretid", "", "Tencent cloud secret id")
var tencentSecretKey = flag.String("secretkey", "", "Tencent cloud secret key")
var dryRun = flag.Bool("dryrun", false, "Don't really do anything")
var debug = flag.Bool("verbose", false, "Enable verbose log")

func init() {
	flag.Parse()

	if *debug || len(os.Getenv("TEO_HOOK_VERBOSE")) > 0 {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func envValue(keys []string, fallback string) string {
	if len(fallback) > 0 {
		return fallback
	}

	for _, key := range keys {
		v := os.Getenv(key)
		if len(v) > 0 {
			return v
		}
	}

	return ""
}

func getCertDomain() string {
	return envValue([]string{"Le_Domain"}, *certDomain)
}

func getTencentClient() *common.Credential {
	credential := common.NewCredential(
		envValue([]string{"Tencent_SecretId", "TENCENT_SECRET_ID"}, *tencentSecretId),
		envValue([]string{"Tencent_SecretKey", "TENCENT_SECRET_KEY"}, *tencentSecretKey))

	return credential
}

func readCertKey() string {
	key, err := os.ReadFile(os.Getenv("CERT_KEY_PATH"))
	if err != nil {
		logrus.Fatalf("Error when read key file on %s: %v", os.Getenv("CERT_KEY_PATH"), err)
	}

	return string(key)
}

func readCertFullchain() string {
	data, err := os.ReadFile(os.Getenv("CERT_FULLCHAIN_PATH"))
	if err != nil {
		logrus.Fatalf("Error when read cert file on %s: %v", os.Getenv("CERT_FULLCHAIN_PATH"), err)
	}

	return string(data)
}
