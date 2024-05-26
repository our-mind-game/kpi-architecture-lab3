package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// State зберігає поточний стан фігури
type State struct {
	backgroundColor color.Color
	hasRect         bool
	rect            RectOperation
	hasFigure       bool
	figures         []FigureOperation
}

var state = State{
	color.White,
	false,
	RectOperation{},
	false,
	[]FigureOperation{},
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
type UpdateOp struct{}

func (op UpdateOp) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), state.backgroundColor, screen.Src)

	if state.hasRect {
		rect := state.rect
		t.Fill(image.Rect(rect.X1, rect.Y1, rect.X2, rect.Y2), color.Black, screen.Src)
	}

	if state.hasFigure {
		for _, figure := range state.figures {
			figureColor := color.RGBA{R: 255, G: 30, B: 30, A: 1}
			t.Fill(image.Rect(figure.CordX-200, figure.CordY+70, figure.CordX+200, figure.CordY-70), figureColor, screen.Src)
			t.Fill(image.Rect(figure.CordX-70, figure.CordY+200, figure.CordX+70, figure.CordY-200), figureColor, screen.Src)
		}
	}
	return true
}

// OperationFunc перетворює функцію оновлення текстури в Operation
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує текстуру у білий колір. Може бути використана як Operation через OperationFunc(WhiteFill)
func WhiteFill(t screen.Texture) {
	state.backgroundColor = color.White
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути використана як Operation через OperationFunc(GreenFill)
func GreenFill(t screen.Texture) {
	state.backgroundColor = color.RGBA{G: 0xff, A: 0xff}
}

// Clear очищує текстуру та зафарбовує у чорний колір. Може бути викоистана як Operation через OperationFunc(Clear)
func Clear(t screen.Texture) {
	state.backgroundColor = color.Black
	state.hasRect = false
	state.rect = RectOperation{}
	state.hasFigure = false
	state.figures = []FigureOperation{}
}

type RectOperation struct{ X1, Y1, X2, Y2 int }

func (rectOp RectOperation) Do(t screen.Texture) bool {
	state.rect = RectOperation{rectOp.X1, rectOp.Y1, rectOp.X2, rectOp.Y2}
	state.hasRect = true
	return false
}

type FigureOperation struct{ CordX, CordY int }

func (f FigureOperation) Do(t screen.Texture) bool {
	state.figures = append(state.figures, FigureOperation{f.CordX, f.CordY})
	state.hasFigure = true
	return false
}

type MoveOperation struct{ X, Y int }

func (m MoveOperation) Do(t screen.Texture) bool {
	for i := range state.figures {
		state.figures[i].CordX += m.X
		state.figures[i].CordY += m.Y
	}
	return false
}
