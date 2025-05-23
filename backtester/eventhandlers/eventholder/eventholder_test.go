package eventholder

import (
	"errors"
	"testing"

	"github.com/antonk9021/qocryptotrader/backtester/common"
	"github.com/antonk9021/qocryptotrader/backtester/eventtypes/order"
	gctcommon "github.com/antonk9021/qocryptotrader/common"
)

func TestReset(t *testing.T) {
	t.Parallel()
	e := &Holder{Queue: []common.Event{}}
	err := e.Reset()
	if !errors.Is(err, nil) {
		t.Errorf("received '%v' expected '%v'", err, nil)
	}
	if e.Queue != nil {
		t.Error("expected nil")
	}

	e = nil
	err = e.Reset()
	if !errors.Is(err, gctcommon.ErrNilPointer) {
		t.Errorf("received '%v' expected '%v'", err, gctcommon.ErrNilPointer)
	}
}

func TestAppendEvent(t *testing.T) {
	t.Parallel()
	e := Holder{Queue: []common.Event{}}
	e.AppendEvent(&order.Order{})
	if len(e.Queue) != 1 {
		t.Error("expected 1")
	}
}

func TestNextEvent(t *testing.T) {
	t.Parallel()
	e := Holder{Queue: []common.Event{}}
	if ev := e.NextEvent(); ev != nil {
		t.Error("expected not ok")
	}

	e = Holder{Queue: []common.Event{
		&order.Order{},
		&order.Order{},
		&order.Order{},
	}}
	if len(e.Queue) != 3 {
		t.Error("expected 3")
	}
	e.NextEvent()
	if len(e.Queue) != 2 {
		t.Error("expected 2")
	}
}
