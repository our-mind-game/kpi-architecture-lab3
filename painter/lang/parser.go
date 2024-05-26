package lang

import (
	"bufio"
	"errors"
	"github.com/our-mind-game/kpi-architecture-lab3/painter"
	"io"
	"strconv"
	"strings"
)

var (
	InvalidOperationError = errors.New("invalid operation specified")
	InvalidParameterError = errors.New("invalid parameter specified")
	NotAllParametersError = errors.New("not all parameters specified")
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	var res []painter.Operation

	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := parseLine(commandLine) // parse the line to get Operation
		if err != nil {
			return nil, err
		}
		res = append(res, op)
	}
	return res, nil
}

func parseLine(lineOp string) (painter.Operation, error) {
	command, parameters, err := getOperationFields(lineOp)

	if err != nil {
		return nil, err
	}

	switch command {
	case "white":
		return painter.OperationFunc(painter.WhiteFill), nil
	case "green":
		return painter.OperationFunc(painter.GreenFill), nil
	case "update":
		return painter.UpdateOp{}, nil
	case "bgrect":
		return painter.RectOperation{
			X1: parameters[0],
			Y1: parameters[1],
			X2: parameters[2],
			Y2: parameters[3],
		}, nil
	case "figure":
		return painter.FigureOperation{
			CordX: parameters[0],
			CordY: parameters[1],
		}, nil
	case "move":
		return painter.MoveOperation{
			X: parameters[0],
			Y: parameters[1],
		}, nil
	case "reset":
		return painter.OperationFunc(painter.Clear), nil
	}
	return nil, InvalidOperationError
}

func getOperationFields(operation string) (string, []int, error) {
	fields := strings.Fields(operation)
	command := fields[0]
	var parameters []int

	for _, field := range fields[1:] {
		res, err := strconv.ParseFloat(field, 32)
		if err != nil {
			return "", nil, InvalidParameterError
		}
		parameters = append(parameters, int(res))
	}

	switch command {
	case "bgrect":
		if len(parameters) < 4 {
			return "", nil, NotAllParametersError
		}
	case "figure", "move":
		if len(parameters) < 2 {
			return "", nil, NotAllParametersError
		}
	}

	return command, parameters, nil
}
