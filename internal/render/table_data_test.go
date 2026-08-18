package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestTableData() *TableData {
	return &TableData{
		Header: Header{
			{Name: "Name"},
			{Name: "State"},
			{Name: "Secret", Hide: true},
		},
		RowEvents: RowEvents{
			NewRowEvent(EventAdd, Row{ID: "i-1", Fields: Fields{"web-server", "running", "shh"}}),
			NewRowEvent(EventAdd, Row{ID: "i-2", Fields: Fields{"db-server", "stopped", "shh"}}),
			NewRowEvent(EventAdd, Row{ID: "i-3", Fields: Fields{"worker", "running", "shh"}}),
		},
	}
}

func rowIDs(td *TableData) []string {
	ids := make([]string, len(td.RowEvents))
	for i, re := range td.RowEvents {
		ids[i] = re.Row.ID
	}
	return ids
}

func TestTableDataFilterEmptyQueryReturnsAllRows(t *testing.T) {
	td := newTestTableData()

	got := td.Filter("")

	assert.Len(t, got.RowEvents, 3)
	assert.Equal(t, []string{"i-1", "i-2", "i-3"}, rowIDs(got))
}

func TestTableDataFilterMatchesAcrossVisibleColumns(t *testing.T) {
	td := newTestTableData()

	got := td.Filter("running")

	assert.Len(t, got.RowEvents, 2)
	assert.Equal(t, []string{"i-1", "i-3"}, rowIDs(got))
}

func TestTableDataFilterIgnoresHiddenColumns(t *testing.T) {
	td := newTestTableData()

	got := td.Filter("shh")

	assert.Len(t, got.RowEvents, 0)
}

func TestTableDataFilterNoMatchesReturnsEmpty(t *testing.T) {
	td := newTestTableData()

	got := td.Filter("zzznomatch")

	assert.Len(t, got.RowEvents, 0)
}

func TestTableDataFilterInvertShowsNonMatchingRows(t *testing.T) {
	td := newTestTableData()

	got := td.Filter("!running")

	assert.Len(t, got.RowEvents, 1)
	assert.Equal(t, []string{"i-2"}, rowIDs(got))
}

func TestTableDataFilterFuzzyMatchesAcrossVisibleColumns(t *testing.T) {
	td := newTestTableData()

	got := td.Filter("-f stp")

	assert.Len(t, got.RowEvents, 1)
	assert.Equal(t, []string{"i-2"}, rowIDs(got))
}

func TestTableDataFilterFuzzyIgnoresHiddenColumns(t *testing.T) {
	td := newTestTableData()

	got := td.Filter("-f shh")

	assert.Len(t, got.RowEvents, 0)
}

func TestTableDataFilterDoesNotMutateOriginal(t *testing.T) {
	td := newTestTableData()

	td.Filter("running")

	assert.Len(t, td.RowEvents, 3)
	assert.Equal(t, []string{"i-1", "i-2", "i-3"}, rowIDs(td))
}
