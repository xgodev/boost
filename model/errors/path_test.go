package errors

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PathSuite struct {
	suite.Suite
}

func TestPathSuite(t *testing.T) {
	suite.Run(t, new(PathSuite))
}

func (s *PathSuite) TestSetSourceTrimPrefix() {

	tt := []struct {
		name string
		got  func() string
		want string
	}{
		{
			name: "source trim prefix",
			got: func() string {
				testPrefix := "/usr/share/"
				savePrefix := SetSourceTrimPrefix(testPrefix)
				defer SetSourceTrimPrefix(savePrefix)
				relative := "github.com/foo/bar/rel.go"
				return filepath.Join(testPrefix, relative)
			},
			want: "/usr/share/github.com/foo/bar/rel.go",
		},
	}

	for _, t := range tt {
		s.Run(t.name, func() {
			got := t.got()
			s.Assert().True(reflect.DeepEqual(got, t.want), "got  %v\nwant %v", got, t.want)
		})
	}
}
