package adcp

import (
	"context"
	"fmt"
	"strconv"
)

var (
	volumeStatus = []byte("volume ?\r\n")
	muteStatus   = []byte("muting ?\r\n")
)

// GetVolumeByBlock returns the volume level of the projector
func (p *Projector) GetVolumeByBlock(ctx context.Context, block string) (int, error) {
	resp, err := p.SendCommand(ctx, p.Address, volumeStatus)
	if err != nil {
		return -1, err
	}

	volume, err := strconv.Atoi(resp)
	if err != nil {
		return -1, err
	}

	return adcpToNormalVolume(volume), nil
}

// SetVolumeByBlock sets the volume level of the projector
func (p *Projector) SetVolumeByBlock(ctx context.Context, block string, level int) error {
	level = normalToAdcpVolume(level)

	cmd := []byte(fmt.Sprintf("volume %v\r\n", level))

	resp, err := p.SendCommand(ctx, p.Address, cmd)
	if err != nil {
		return err
	}

	if resp != "ok" {
		return fmt.Errorf("unable to set volume to %v: %s", level, resp)
	}

	return nil
}

// GetMutedByBlock returns whether the projector is muted or not
func (p *Projector) GetMutedByBlock(ctx context.Context, block string) (bool, error) {
	resp, err := p.SendCommand(ctx, p.Address, muteStatus)
	if err != nil {
		return false, err
	}

	var muted bool
	switch resp {
	case `"on"`:
		muted = true
	case `"off"`:
		muted = false
	default:
		return false, fmt.Errorf("unknown muted state '%s'", resp)
	}

	return muted, nil
}

// SetMutedByBlock sets the muted status of the projector
func (p *Projector) SetMutedByBlock(ctx context.Context, block string, muted bool) error {
	var str string
	switch muted {
	case true:
		str = "on"
	case false:
		str = "off"
	}

	cmd := []byte(fmt.Sprintf("muting \"%s\"\r\n", str))
	resp, err := p.SendCommand(ctx, p.Address, cmd)
	if err != nil {
		return err
	}

	if resp != "ok" {
		return fmt.Errorf("unable to set muted state to %v: %s", muted, resp)
	}

	return nil
}

// the volume level that the projectors put out is only really
// useful from 0-50(ish). above 50 or so, the volume seems to stay
// somewhat constant. these functions map the given 0-100 volume
// to the min and the max that we set.

const (
	minAdcp = 0
	maxAdcp = 50

	adcpConversion = 100 / maxAdcp
)

func normalToAdcpVolume(level int) int {
	switch {
	case level >= 0 && level <= 100:
		return level / adcpConversion
	case level < 0:
		return minAdcp
	case level > 100:
		return maxAdcp
	default:
		return level
	}
}

func adcpToNormalVolume(level int) int {
	switch {
	case level >= minAdcp && level <= maxAdcp:
		return level * adcpConversion
	case level < minAdcp:
		return minAdcp
	case level > maxAdcp:
		return maxAdcp
	default:
		return level
	}
}
