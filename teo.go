package main

import (
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"time"
)

func watchHostUpdate(sslClient *ssl.Client, id string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		req := ssl.NewDescribeHostUpdateRecordDetailRequest()
		req.DeployRecordId = common.StringPtr(id)
		res, err := sslClient.DescribeHostUpdateRecordDetail(req)
		if err != nil {
			logrus.Infof("Error on fetch host update detail: %v", err)
			continue
		}

		total := *res.Response.TotalCount
		done := *res.Response.SuccessTotalCount + *res.Response.FailedTotalCount
		if done >= total {
			return
		}
	}
}
