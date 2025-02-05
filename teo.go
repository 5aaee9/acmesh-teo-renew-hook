package main

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func watchHostUpdate(sslClient *ssl.Client, id uint64) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		req := ssl.NewDescribeHostUpdateRecordDetailRequest()
		req.DeployRecordId = common.StringPtr(strconv.FormatUint(id, 10))
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
