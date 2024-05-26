package lang

import (
	"errors"
	"github.com/our-mind-game/kpi-architecture-lab3/painter"
	"reflect"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	p := Parser{}

	expected := []painter.Operation{
		painter.OperationFunc(painter.Clear),
		painter.OperationFunc(painter.GreenFill),
		painter.RectOperation{X1: 1, Y1: 2, X2: 3, Y2: 4},
		painter.FigureOperation{CordX: 1, CordY: 2},
		painter.OperationFunc(painter.WhiteFill),
		painter.FigureOperation{CordX: 3, CordY: 4},
		painter.UpdateOp{},
	}

	r := strings.NewReader("reset\n green\n bgrect 7 7 22 22\n figure 14 88\n white\n figure 8 8\n update")
	res, err := p.Parse(r)
	if err != nil {
		t.Fatal("Should parse operations")
	}

	for i := range expected {
		if !reflect.DeepEqual(reflect.TypeOf(expected[i]), reflect.TypeOf(res[i])) {
			t.Fatal("Should parse operations in correct order")
		}
	}
}

func TestParser_InvalidOperationError(t *testing.T) {
	p := Parser{}
	r := strings.NewReader("rect")
	_, err := p.Parse(r)

	if !errors.Is(err, InvalidOperationError) {
		t.Fatal("Should catch error: InvalidOperationError")
	}

}

func TestParser_InvalidParametersError(t *testing.T) {
	p := Parser{}
	r := strings.NewReader("bgrect 100 golang 300 typescript")
	_, err := p.Parse(r)

	if !errors.Is(err, InvalidParameterError) {
		t.Fatal("Should catch error: InvalidParameterError")
	}

}

func TestParser_NotAllParametersError(t *testing.T) {
	p := Parser{}
	r := strings.NewReader("bgrect 200 400 150")
	_, err := p.Parse(r)

	if !errors.Is(err, NotAllParametersError) {
		t.Fatal("Should catch error: NotAllParametersError")
	}

}
