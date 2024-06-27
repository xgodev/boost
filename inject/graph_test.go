package inject

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type NewGraphFromEntriesTestSuite struct {
	suite.Suite
	testData map[string][]annotation.Entry
}

func TestNewGraphFromEntriesTestSuite(t *testing.T) {
	suite.Run(t, new(NewGraphFromEntriesTestSuite))
}

func (suite *NewGraphFromEntriesTestSuite) SetupSuite() {

	suite.testData = make(map[string][]annotation.Entry)

	files, err := filepath.Glob("testdata/inject/mkgraph/*.yaml")
	suite.Require().NoError(err)

	for _, file := range files {
		id := filepath.Base(file)
		entries, err := suite.loadEntriesFromYAML(file)
		suite.Require().NoError(err)

		suite.testData[id] = entries
	}
}

func (suite *NewGraphFromEntriesTestSuite) loadEntriesFromYAML(filename string) ([]annotation.Entry, error) {
	var entries []annotation.Entry

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (suite *NewGraphFromEntriesTestSuite) TestNewGraphFromEntries() {
	testCases := []struct {
		name      string
		id        string // ID do cenário de teste
		expectErr bool
		errType   error
	}{
		{
			name:      "valid",
			id:        "1_valid.yaml",
			expectErr: false,
		},
		{
			name:      "provider not found",
			id:        "2_provider_notfound.yaml",
			expectErr: true,
			errType:   errors.NotFoundf("tipo de erro esperado"),
		},
		{
			name:      "provide index not found",
			id:        "3_provide_index_notfound.yaml",
			expectErr: true,
			errType:   errors.NotValidf("tipo de erro esperado"),
		},
		{
			name:      "inject index not found",
			id:        "4_inject_index_notfound.yaml",
			expectErr: true,
			errType:   errors.NotValidf("tipo de erro esperado"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			entries := suite.testData[tc.id]
			_, err := NewGraphFromEntries(context.Background(), entries)

			if tc.expectErr {
				suite.Error(err)
				suite.IsType(tc.errType, err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
