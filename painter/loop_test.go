package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Receiver(t *testing.T) {
	var l Loop
	var tr TestReceiver
	l.Receiver = &tr
	l.Start(mockScreen{})

	l.Post(OperationFunc(Clear))
	l.Post(OperationFunc(WhiteFill))
	l.Post(OperationFunc(GreenFill))
	l.Post(UpdateOp{})

	if tr.LastTexture != nil {
		t.Fatal("Should get texture the last")
	}

	l.StopAndWait()

	if tr.LastTexture == nil {
		t.Fatal("Should get texture")
	}
}

func TestLoop_Operations(t *testing.T) {
	var l Loop
	var tr TestReceiver
	l.Receiver = &tr
	l.Start(mockScreen{})

	l.Post(OperationFunc(GreenFill))
	l.Post(OperationFunc(WhiteFill))
	l.Post(OperationFunc(GreenFill))
	l.Post(UpdateOp{})

	if len(l.mq.ops) != 4 {
		t.Fatal("Should return correct number of operations: 4")
	}

	l.StopAndWait()

	tx, _ := tr.LastTexture.(*mockTexture)
	if tx.FillCnt != 1 {
		t.Fatal("Should call Fill once")
	}

	if len(l.mq.ops) != 0 {
		t.Fatal("Should return correct number of operations: 0")
	}
}

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr TestReceiver
	)
	l.Receiver = &tr
	l.Start(mockScreen{})
	var testOps []string

	l.Post(OperationFunc(func(t screen.Texture) {
		testOps = append(testOps, "firstCommand")

		l.Post(OperationFunc(func(t screen.Texture) {
			testOps = append(testOps, "thirdCommand")
		}))
	}))

	l.Post(OperationFunc(func(t screen.Texture) {
		testOps = append(testOps, "secondCommand")
	}))

	l.StopAndWait()

	if !reflect.DeepEqual(testOps, []string{"firstCommand", "secondCommand", "thirdCommand"}) {
		t.Fatal("Should return correct operations order")
	}
}

type TestReceiver struct {
	LastTexture screen.Texture
}

func (tr *TestReceiver) Update(t screen.Texture) {
	tr.LastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("Implement me")
}
func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}
func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("Implement me")
}

type mockTexture struct {
	FillCnt int
}

func (m *mockTexture) Release() {}
func (m *mockTexture) Size() image.Point {
	return image.Pt(400, 400)
}
func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Pt(400, 400)}
}
func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	panic("Implement me")
}
func (m *mockTexture) Fill(dp image.Rectangle, src color.Color, op draw.Op) {
	m.FillCnt++
}
