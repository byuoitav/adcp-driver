package adcp

import (
	"context"
	"fmt"
)

var (
	// PowerStatus gets the projector's power status
	PowerStatus = []byte("power_status ?\r\n")

	// PowerOn powers on the projector
	PowerOn = []byte("power \"on\"\r\n")

	// PowerStandby powers off the projector
	PowerStandby = []byte("power \"standby\"\r\n")
)

// GetPower returns the status of the projector
func (p *Projector) GetPower(ctx context.Context) (bool, error) {
	state := false

	resp, err := p.SendCommand(ctx, p.Address, PowerStatus)
	if err != nil {
		return false, err
	}
	switch resp {
	case `"standby"`:
	case `"startup"`:
		state = true
	case `"on"`:
		state = true
	case `"cooling1"`:
	case `"cooling2"`:
	case `"saving_cooling1"`:
	case `"saving_cooling2"`:
	case `"saving_standby"`:
	default:
		return state, fmt.Errorf("unknown power state '%s'", resp)
	}

	return state, nil
}

// SetPower sets the status of the projector
func (p *Projector) SetPower(ctx context.Context, power bool) error {
	var cmd []byte
	switch {
	case power:
		cmd = PowerOn
	case !power:
		cmd = PowerStandby
	default:
		return fmt.Errorf("unable to set power state to %q: must be %q or %q", power, "on", "standby")
	}

	_, err := p.SendCommand(ctx, p.Address, cmd)
	return err
}
