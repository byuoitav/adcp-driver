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

// GetAudioVideoInputs returns the current input that the projector is set to
func (p *Projector) GetAudioVideoInputs(ctx context.Context) (map[string]string, error) {
	toReturn := make(map[string]string)
	resp, err := p.SendCommand(ctx, p.Address, InputStatus)
	if err != nil {
		return toReturn, err
	}

	toReturn[""] = strings.Trim(resp, "\"")
	return toReturn, nil
}

// SetAudioVideoInput sets the current input of the projector to the given input
func (p *Projector) SetAudioVideoInput(ctx context.Context, output, input string) error {
	cmd := []byte(fmt.Sprintf("input \"%s\"\r\n", input))
	resp, err := p.SendCommand(ctx, p.Address, cmd)
	if err != nil {
		return err
	}
	if resp != "ok" {
		return fmt.Errorf("unable to set input to %v: %s", input, resp)
	}

	return nil
}

// GetActiveSignal checks to see if the projector has an active input signal and returns the result
func (p *Projector) GetActiveSignal(ctx context.Context, port string) (bool, error) {
	var active bool
	resp, err := p.SendCommand(ctx, p.Address, ActiveSignal)
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
