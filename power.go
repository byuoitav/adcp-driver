package adcp

import (
	"context"
	"fmt"
	"strings"
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
func (p *Projector) GetPower(ctx context.Context, addr string) (string, error) {
	var state string

	resp, err := p.SendCommand(ctx, addr, PowerStatus)
	if err != nil {
		return "", err
	}
	switch resp {
	case `"standby"`:
		state = "standby"
	case `"startup"`:
		state = "on"
	case `"on"`:
		state = "on"
	case `"cooling1"`:
		state = "standby"
	case `"cooling2"`:
		state = "standby"
	case `"saving_cooling1"`:
		state = "standby"
	case `"saving_cooling2"`:
		state = "standby"
	case `"saving_standby"`:
		state = "standby"
	default:
		return "", fmt.Errorf("unknown power state '%s'", resp)
	}

	return state, nil
}

// SetPower sets the status of the projector
func (p *Projector) SetPower(ctx context.Context, addr, power string) error {
	var cmd []byte
	switch {
	case strings.EqualFold(power, "on"):
		cmd = PowerOn
	case strings.EqualFold(power, "standby"):
		cmd = PowerStandby
	default:
		return fmt.Errorf("unable to set power state to %q: must be %q or %q", power, "on", "standby")
	}

	_, err := p.SendCommand(ctx, addr, cmd)
	return err
}
