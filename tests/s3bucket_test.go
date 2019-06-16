package tests

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pixelguy95/deploy-static-web/deploy"
)

func TestBucketCreationAndDeletion(t *testing.T) {
	path := filepath.Join("testdata", "config.json")
	t.Log(path)
	config := deploy.LoadConfigurations(path)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	config.DomainName = strconv.FormatInt(r1.Int63(), 16) + config.DomainName
	sess, err := deploy.GetNewSession(config)

	if err != nil {
		t.Error(err, "Session creation failed")
	}

	s3Session := s3.New(sess)
	err = deploy.CreateBucket(config, s3Session)

	if err != nil {
		t.Error("Bucket creation failed")
	}

	err = deploy.DestroyBucket(config, s3Session)
	if err != nil {
		t.Error("Bucket destruction failed")
	}

}
