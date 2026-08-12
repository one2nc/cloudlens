package model

import (
	"testing"

	"github.com/sahilm/fuzzy"
	"github.com/stretchr/testify/assert"
)

type spyResourceViewerListener struct {
	lines   []string
	matches fuzzy.Matches
	failed  error
	calls   int
}

func (s *spyResourceViewerListener) ResourceChanged(lines []string, matches fuzzy.Matches) {
	s.lines = lines
	s.matches = matches
	s.calls++
}

func (s *spyResourceViewerListener) ResourceFailed(err error) {
	s.failed = err
}

func newTestDescribe(lines []string) (*Describe, *spyResourceViewerListener) {
	d := NewDescribe("ec2", "some/path")
	d.lines = lines
	spy := &spyResourceViewerListener{}
	d.AddListener(spy)
	return d, spy
}

func TestDescribeFilterFindsMatchingLine(t *testing.T) {
	d, spy := newTestDescribe([]string{
		`"InstanceId": "i-123"`,
		`"Name": "cloudlens-demo-web"`,
		`"State": "running"`,
	})

	d.Filter("web")

	assert.Len(t, spy.matches, 1)
	assert.Equal(t, 1, spy.matches[0].Index)
}

func TestDescribeFilterMatchedIndexesListEveryMatchedPosition(t *testing.T) {
	d, spy := newTestDescribe([]string{`"Name": "cloudlens-demo-web"`})

	d.Filter("web")

	assert.Equal(t, []int{24, 25, 26}, spy.matches[0].MatchedIndexes)
}

func TestDescribeFilterIsCaseInsensitive(t *testing.T) {
	d, spy := newTestDescribe([]string{`"Name": "cloudlens-demo-WEB"`})

	d.Filter("web")

	assert.Len(t, spy.matches, 1)
}

func TestDescribeFilterFindsMultipleMatches(t *testing.T) {
	d, spy := newTestDescribe([]string{
		`"Name": "cloudlens-demo-web"`,
		`"State": "running"`,
		`"Tag": "web-tier"`,
	})

	d.Filter("web")

	assert.Len(t, spy.matches, 2)
	assert.Equal(t, 0, spy.matches[0].Index)
	assert.Equal(t, 2, spy.matches[1].Index)
}

func TestDescribeFilterNoMatchesReturnsEmptyNotError(t *testing.T) {
	d, spy := newTestDescribe([]string{`"State": "running"`})

	d.Filter("zzznomatch")

	assert.Empty(t, spy.matches)
	assert.Equal(t, 1, spy.calls)
}

func TestDescribeFilterEmptyQueryReturnsNoMatches(t *testing.T) {
	d, spy := newTestDescribe([]string{`"State": "running"`})

	d.Filter("")

	assert.Empty(t, spy.matches)
}

func TestDescribeFilterFuzzySelectorUsesFuzzyMatch(t *testing.T) {
	d, spy := newTestDescribe([]string{
		`"Name": "cloudlens-demo-web"`,
		`"State": "running"`,
	})

	d.Filter("-f cdw")

	assert.Len(t, spy.matches, 1)
	assert.Equal(t, 0, spy.matches[0].Index)
}

func TestDescribeClearFilterResetsQueryAndReFires(t *testing.T) {
	lines := []string{
		`"Name": "cloudlens-demo-web"`,
		`"State": "running"`,
	}
	d, spy := newTestDescribe(lines)

	d.Filter("web")
	assert.Len(t, spy.matches, 1)

	d.ClearFilter()

	assert.Equal(t, "", d.query)
	assert.Empty(t, spy.matches)
	assert.Equal(t, lines, spy.lines)
	assert.Equal(t, 2, spy.calls)
}
