package taint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParserSuite struct {
	suite.Suite
}

func (suite *ParserSuite) SetupSuite() {
	Sources = make([]*Data, 0)
	Sinks = make([]*Data, 0)
}

/*
Test object: File: ../../sourcesAndSinks.txt
expect: 40 sources, 5 sinks and 6 source-interfaces, 0 sink-interfaces
*/
func (suite *ParserSuite) TestParser() {
	err := Read("../../sourcesAndSinks.txt")
	assert.Nil(suite.T(), err)
	for _, source := range Sources {
		suite.T().Logf("%s\n", source.String())
	}
	for _, sink := range Sinks {
		suite.T().Logf("%s\n", sink.String())
	}
	assert.Equal(suite.T(), 40, len(Sources))
	assert.Equal(suite.T(), 5, len(Sinks))

	expected := 6
	actual := 0
	for _, source := range Sources {
		if source.IsInterface() {
			actual++
		}
	}
	assert.Equal(suite.T(), expected, actual)

	expected = 0
	actual = 0
	for _, sink := range Sinks {
		if sink.IsInterface() {
			actual++
		}
	}
	assert.Equal(suite.T(), expected, actual)
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}
