package adcp

import (
	"context"
	"fmt"
	"strings"
)

var (
	// InputStatus is the command to ask the projector what input it is on
	InputStatus = []byte("input ?\r\n")

	// ActiveSignal is the command to ask the projector if it has an active input signal
	ActiveSignal = []byte("signal ?\r\n")
)

// GetInput returns the current input that the projector is set to
func (p *Projector) GetInput(ctx context.Context, addr string) (string, error) {
	resp, err := p.SendCommand(ctx, addr, InputStatus)
	if err != nil {
		return "", err
	}

	return strings.Trim(resp, "\""), nil
}

// SetInput sets the current input of the projector to the given input
func (p *Projector) SetInput(ctx context.Context, addr, input string) error {
	cmd := []byte(fmt.Sprintf("input \"%s\"\r\n", input))
	resp, err := p.SendCommand(ctx, addr, cmd)
	if err != nil {
		return err
	}
	if resp != "ok" {
		return fmt.Errorf("unable to set input to %v: %s", input, resp)
	}

	return nil
}

// GetActiveSignal checks to see if the projector has an active input signal and returns the result
func (p *Projector) GetActiveSignal(ctx context.Context, addr string) (bool, error) {
	var active bool
	resp, err := p.SendCommand(ctx, addr, ActiveSignal)
	if err != nil {
		return false, err
	}

	switch resp {
	case `"Invalid"`:
		active = false
	case "ok":
		active = false
	default:
		active = true
	}

	return active, nil
}
