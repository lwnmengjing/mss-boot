/*
 * @Author: lwnmengjing
 * @Date: 2022/7/18 10:06:11
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2022/7/18 10:06:11
 */

package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/mss-boot-io/mss-boot/pkg/config/source"
)

type Source struct {
	opt source.Options
}

func (s *Source) Open(string) (fs.File, error) {
	return nil, errors.New("method Get not implemented")
}

func (s *Source) ReadFile(name string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), s.opt.Timeout)
	defer cancel()
	object, err := s.opt.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.opt.Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", s.opt.Dir, name)),
	})
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()
	return io.ReadAll(object.Body)
}

// New source
func New(options ...source.Option) (*Source, error) {
	s := &Source{}
	for _, opt := range options {
		opt(&s.opt)
	}
	if s.opt.Timeout == 0 {
		s.opt.Timeout = 5 * time.Second
	}
	fmt.Println(s.opt.Dir, s.opt.ProjectName)
	if s.opt.ProjectName != "" {
		s.opt.Dir = s.opt.Dir[strings.Index(s.opt.Dir, s.opt.ProjectName+"/"):]
	}
	if s.opt.Client != nil {
		return s, nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), s.opt.Timeout)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(s.opt.Region))
	if err != nil {
		return nil, err
	}
	s.opt.Client = s3.NewFromConfig(cfg)
	return s, nil
}
