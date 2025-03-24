package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func main() {
	cred := getTencentClient()
	sslClient, err := ssl.NewClient(cred, regions.Shanghai, profile.NewClientProfile())
	if err != nil {
		logrus.Fatal("Error on create ssl client: %v", err)
	}

	describeRequest := ssl.NewDescribeCertificatesRequest()
	describeRequest.SearchKey = common.StringPtr(getCertDomain())
	describeRequest.Limit = common.Uint64Ptr(1000)
	describeRequest.CertificateType = common.StringPtr("SVR")

	res, err := sslClient.DescribeCertificates(describeRequest)
	if err != nil {
		logrus.Fatalf("Error on describe certificate: %s", err)
	}

	var certId string
	for _, cert := range res.Response.Certificates {
		if *cert.Domain == getCertDomain() {
			certId = *cert.CertificateId
		} else {
			logrus.Debugf("Cert domain %s not equal target %s, skipped", *cert.Domain, getCertDomain())
		}
	}

	if certId == "" {
		logrus.Fatal("No certificate found")

		if *dryRun {
			return
		}

		createReq := ssl.NewUploadCertificateRequest()
		createReq.CertificatePrivateKey = common.StringPtr(readCertKey())
		createReq.CertificatePublicKey = common.StringPtr(readCertFullchain())
		createReq.Repeatable = common.BoolPtr(false)
		res, err := sslClient.UploadCertificate(createReq)
		if err != nil {
			logrus.Fatalf("create cert failed: %v", err)
		}

		logrus.Infof("create new cert with id: %s", *res.Response.CertificateId)
		return
	}

	logrus.Infof("Tencent SSL certificate id: %s", certId)
	if *dryRun {
		return
	}

	updateRequest := ssl.NewUpdateCertificateInstanceRequest()
	updateRequest.OldCertificateId = common.StringPtr(certId)
	updateRequest.CertificatePrivateKey = common.StringPtr(readCertKey())
	updateRequest.CertificatePublicKey = common.StringPtr(readCertFullchain())
	updateRequest.ResourceTypes = common.StringPtrs([]string{
		"teo",
	})
	updateRequest.ExpiringNotificationSwitch = common.Uint64Ptr(0)
	updateRequest.Repeatable = common.BoolPtr(false)
	updateRequest.AllowDownload = common.BoolPtr(true)
	var updateResponse *ssl.UpdateCertificateInstanceResponse

	for {
		var err error
		updateResponse, err = sslClient.UpdateCertificateInstance(updateRequest)
		if err != nil {
			logrus.Fatalf("Error on update certificate: %s", err)
		}

		logrus.Debugf("UpdateCertificateInstance: %v", updateResponse.ToJsonString())
		if *updateResponse.Response.DeployRecordId != 0 {
			break
		}

		time.Sleep(time.Second)
	}

	if *updateResponse.Response.DeployStatus != 1 {
		logrus.Fatalf("UpdateCertificateInstance: bad status => %s", updateResponse.ToJsonString())
	}

	deleteRequest := ssl.NewDeleteCertificateRequest()
	deleteRequest.CertificateId = common.StringPtr(certId)
	_, err = sslClient.DeleteCertificate(deleteRequest)
	if err != nil {
		logrus.Fatalf("Error on delete certificate: %s", err)
	}
}
