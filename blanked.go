package adcp

import (
	"context"
	"fmt"
)

var (
	// BlankStatus is the command to ask the projector if it is blanked or not
	BlankStatus = []byte("blank ?\r\n")

	// Blank is the command to send to the projector to tell it to blank
	Blank = []byte("blank \"on\"\r\n")

	// Unblank is the command to send to the projector to tell it to unblank
	Unblank = []byte("blank \"off\"\r\n")
)

// GetBlanked asks the projector if it is blanked or not and returns the result
func (p *Projector) GetBlanked(ctx context.Context) (bool, error) {
	var blanked bool

	resp, err := p.SendCommand(ctx, p.Address, BlankStatus)
	if err != nil {
		return false, err
	}

	switch resp {
	case `"on"`:
		blanked = true
	case `"off"`:
		blanked = false
	default:
		return false, fmt.Errorf("unknown blanked state '%s'", resp)
	}

	return blanked, nil
}

// SetBlanked tells the projector to blank or unblank itself
func (p *Projector) SetBlanked(ctx context.Context, blanked bool) error {
	cmd := Unblank
	if blanked {
		cmd = Blank
	}

	resp, err := p.SendCommand(ctx, p.Address, cmd)
	if err != nil {
		return err
	}

	if resp != "ok" {
		return fmt.Errorf("unable to set blanked state to %v: %s", blanked, resp)
	}

	return nil
}
