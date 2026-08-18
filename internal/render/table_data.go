package render

import (
	"regexp"
	"strings"
	"sync"

	"github.com/sahilm/fuzzy"
)

const filterFieldSpacer = " "

var tableDataFuzzyRx = regexp.MustCompile(`\A\-f`)

type TableData struct {
	Header    Header
	RowEvents RowEvents
	mx        sync.RWMutex
}

// NewTableData returns a new table.
func NewTableData() *TableData {
	return &TableData{}
}

// Empty checks if there are no entries.
func (t *TableData) Empty() bool {
	t.mx.RLock()
	defer t.mx.RUnlock()

	return len(t.RowEvents) == 0
}

// Count returns the number of entries.
func (t *TableData) Count() int {
	t.mx.RLock()
	defer t.mx.RUnlock()

	return len(t.RowEvents)
}

// IndexOfHeader return the index of the header.
func (t *TableData) IndexOfHeader(h string) int {
	return t.Header.IndexOf(h, false)
}

// Customize returns a new model with customized column layout.
func (t *TableData) Customize(cols []string, wide bool) *TableData {
	res := TableData{
		Header:    t.Header.Customize(cols, wide),
	}
	ids := t.Header.MapIndices(cols, wide)
	res.RowEvents = t.RowEvents.Customize(ids)

	return &res
}


// Clear clears out the entire table.
func (t *TableData) Clear() {
	t.Header, t.RowEvents = Header{}, RowEvents{}
}

func (t *TableData) Filter(q string) *TableData {
	td := t.Clone()
	if q == "" {
		return td
	}
	if tableDataFuzzyRx.MatchString(q) {
		td.RowEvents = t.fuzzyFilter(strings.TrimSpace(q[2:]))
		return td
	}
	invert := strings.HasPrefix(q, "!")
	if invert {
		q = q[1:]
	}
	rr, err := t.rxFilter(q, invert)
	if err != nil {
		return td
	}
	td.RowEvents = rr

	return td
}

func (t *TableData) visibleFields(re RowEvent) string {
	ff := make([]string, 0, len(re.Row.Fields))
	for i, f := range re.Row.Fields {
		if i < len(t.Header) && t.Header[i].Hide {
			continue
		}
		ff = append(ff, f)
	}

	return strings.Join(ff, filterFieldSpacer)
}

func (t *TableData) rxFilter(q string, invert bool) (RowEvents, error) {
	rx, err := regexp.Compile(`(?i)` + q)
	if err != nil {
		return nil, err
	}
	rr := make(RowEvents, 0, len(t.RowEvents))
	for _, re := range t.RowEvents {
		match := rx.MatchString(t.visibleFields(re))
		if match != invert {
			rr = append(rr, re)
		}
	}

	return rr, nil
}

func (t *TableData) fuzzyFilter(q string) RowEvents {
	ff := make([]string, len(t.RowEvents))
	for i, re := range t.RowEvents {
		ff[i] = t.visibleFields(re)
	}
	matches := fuzzy.Find(q, ff)
	rr := make(RowEvents, 0, len(matches))
	for _, m := range matches {
		rr = append(rr, t.RowEvents[m.Index])
	}

	return rr
}

// Clone returns a copy of the table.
func (t *TableData) Clone() *TableData {
	return &TableData{
		Header:    t.Header.Clone(),
		RowEvents: t.RowEvents.Clone(),
	}
}

// SetHeader sets table header.
func (t *TableData) SetHeader(h Header) {
	t.Header = h
}

// Update computes row deltas and update the table data.
func (t *TableData) Update(rows Rows) {
	empty := t.Empty()
	kk := make(map[string]struct{}, len(rows))
	t.mx.Lock()
	{
		for _, row := range rows {
			kk[row.ID] = struct{}{}
			if empty {
				t.RowEvents = append(t.RowEvents, NewRowEvent(EventAdd, row))
				continue
			}
			t.RowEvents = append(t.RowEvents, NewRowEvent(EventAdd, row))
		}
	}
	t.mx.Unlock()

	if !empty {
		t.Delete(kk)
	}
}

// Delete removes items in cache that are no longer valid.
func (t *TableData) Delete(newKeys map[string]struct{}) {
	t.mx.Lock()
	{
		var victims []string
		for _, re := range t.RowEvents {
			if _, ok := newKeys[re.Row.ID]; !ok {
				victims = append(victims, re.Row.ID)
			}
		}
		for _, id := range victims {
			t.RowEvents = t.RowEvents.Delete(id)
		}
	}
	t.mx.Unlock()
}
