package errors_test

import (
	"fmt"
	"io"
	"testing"

	"golang.org/x/xerrors"

	"github.com/stretchr/testify/suite"

	"github.com/kysre/TurtleMQ/leader/pkg/errors"
)

type ErrorsTestSuite struct {
	suite.Suite
}

func TestErrorsTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}

func (s *ErrorsTestSuite) TestNewShouldContainMessage() {
	err := errors.New("hello")
	s.Contains(err.Error(), "hello")
}

func (s *ErrorsTestSuite) TestNewShouldContainStackTraceInFormat() {
	err := errors.New("hello")
	str := fmt.Sprintf("%+v", err)
	s.Contains(str, "errors_test.go")
}

func (s *ErrorsTestSuite) TestNewShouldContainExtras() {
	err := errors.NewWithExtra("hello", map[string]interface{}{
		"key": "value",
	})

	extras := errors.Extras(err)
	value, ok := extras["key"]
	s.True(ok)
	s.Equal("value", value)
}

func (s *ErrorsTestSuite) TestWrapShouldBeUnwrappable() {
	cause := errors.New("cause")
	err := errors.Wrap(cause, "another error")
	s.Equal(cause, xerrors.Unwrap(err))
}

func (s *ErrorsTestSuite) TestWrapShouldBeIsCompatible() {
	err := errors.Wrap(io.EOF, "hello")
	s.NotEqual(io.EOF, err)
	s.True(xerrors.Is(err, io.EOF))
}

func (s *ErrorsTestSuite) TestWrapShouldMergeExtras() {
	cause := errors.NewWithExtra("cause", map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})
	err := errors.WrapWithExtra(cause, "another err", map[string]interface{}{
		"key2": "anothervalue2",
	})

	extras := errors.Extras(err)

	value1, ok := extras["key1"]
	s.True(ok)
	s.Equal("value1", value1)

	value2, ok := extras["key2"]
	s.True(ok)
	s.Equal("anothervalue2", value2)
}
