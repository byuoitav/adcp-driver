package adcp

import (
	"context"
	"encoding/json"
	"strings"
)

// HardwareInfo contains the common information for device hardware information
type HardwareInfo struct {
	Hostname              string           `json:"hostname,omitempty"`
	ModelName             string           `json:"model_name,omitempty"`
	SerialNumber          string           `json:"serial_number,omitempty"`
	BuildDate             string           `json:"build_date,omitempty"`
	FirmwareVersion       string           `json:"firmware_version,omitempty"`
	ProtocolVersion       string           `json:"protocol_version,omitempty"`
	NetworkInfo           NetworkInfo      `json:"network_information,omitempty"`
	FilterStatus          string           `json:"filter_status,omitempty"`
	WarningStatus         []string         `json:"warning_status,omitempty"`
	ErrorStatus           []string         `json:"error_status,omitempty"`
	PowerStatus           string           `json:"power_status,omitempty"`
	PowerSavingModeStatus string           `json:"power_saving_mode_status,omitempty"`
	TimerInfo             []map[string]int `json:"timer_info,omitempty"`
	Temperature           string           `json:"temperature,omitempty"`
}

// NetworkInfo contains the network information for the device
type NetworkInfo struct {
	IPAddress  string   `json:"ip_address,omitempty"`
	MACAddress string   `json:"mac_address,omitempty"`
	Gateway    string   `json:"gateway,omitempty"`
	DNS        []string `json:"dns,omitempty"`
}

var (
	modelName   = []byte("modelname ?\r\n")
	ipAddr      = []byte("ipv4_ip_address ?\r\n")
	gateway     = []byte("ipv4_default_gateway ?\r\n")
	dns         = []byte("ipv4_dns_server1 ?\r\n")
	dns2        = []byte("ipv4_dns_server2 ?\r\n")
	macAddr     = []byte("mac_address ?\r\n")
	serialNum   = []byte("serialnum ?\r\n")
	filter      = []byte("filter_status ?\r\n")
	powerStatus = []byte("power_status ?\r\n")
	warnings    = []byte("warning ?\r\n")
	errors      = []byte("error ?\r\n")
	timer       = []byte("timer ?\r\n")
)

// GetInfo returns the hardware information of the projector
func (p *Projector) GetInfo(ctx context.Context, addr string) (interface{}, error) {
	var info HardwareInfo

	// model name
	resp, err := p.SendCommand(ctx, addr, modelName)
	if err != nil {
		return info, err
	}

	info.ModelName = strings.Trim(resp, "\"")

	// ip address
	resp, err = p.SendCommand(ctx, addr, ipAddr)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.IPAddress = strings.Trim(resp, "\"")

	// gateway
	resp, err = p.SendCommand(ctx, addr, gateway)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.Gateway = strings.Trim(resp, "\"")

	// dns
	resp, err = p.SendCommand(ctx, addr, dns)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.DNS = append(info.NetworkInfo.DNS, strings.Trim(resp, "\""))

	resp, err = p.SendCommand(ctx, addr, dns2)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.DNS = append(info.NetworkInfo.DNS, strings.Trim(resp, "\""))

	// mac address
	resp, err = p.SendCommand(ctx, addr, macAddr)
	if err != nil {
		return info, err
	}

	info.NetworkInfo.MACAddress = strings.Trim(resp, "\"")

	// serial number
	resp, err = p.SendCommand(ctx, addr, serialNum)
	if err != nil {
		return info, err
	}

	info.SerialNumber = strings.Trim(resp, "\"")

	// filter status
	resp, err = p.SendCommand(ctx, addr, filter)
	if err != nil {
		return info, err
	}

	info.FilterStatus = strings.Trim(resp, "\"")

	// power status
	resp, err = p.SendCommand(ctx, addr, powerStatus)
	if err != nil {
		return info, err
	}

	info.PowerStatus = strings.Trim(resp, "\"")

	// warnings
	resp, err = p.SendCommand(ctx, addr, warnings)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.WarningStatus)
	if err != nil {
		return info, err
	}

	// errors
	resp, err = p.SendCommand(ctx, addr, errors)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.ErrorStatus)
	if err != nil {
		return info, err
	}

	// timer info
	resp, err = p.SendCommand(ctx, addr, timer)
	if err != nil {
		return info, err
	}

	err = json.Unmarshal([]byte(resp), &info.TimerInfo)
	if err != nil {
		return info, err
	}

	return info, nil
}
