package proc

import (
	"regexp"
	"strconv"
	"strings"
)

const PROC_NET_DATA = "/proc/net/dev"

type NetData struct {
	Bytes      int64
	Packets    int64
	Errs       int64
	Drop       int64
	Fifo       int64
	Frame      int64
	Compressed int64
	Multicast  int64
}

type NetDev struct {
	Name         string
	Receive      NetData
	ReceiveDiff  NetData
	Transmit     NetData
	TransmitDiff NetData
}

func netDataDiff(od NetData, nd NetData) NetData {
	diff := NetData{}
	diff.Bytes = nd.Bytes - od.Bytes
	diff.Packets = nd.Packets - od.Packets
	diff.Errs = nd.Errs - od.Errs
	diff.Drop = nd.Drop - od.Drop
	diff.Fifo = nd.Fifo - od.Fifo
	diff.Frame = nd.Frame - od.Frame
	diff.Compressed = nd.Compressed - od.Compressed
	diff.Multicast = nd.Multicast - od.Multicast
	return diff
}

func getDev(name string, devs map[string]NetDev) NetDev {
	return devs[name]
	/*
		for _, d := range devs {
			if d.Name == name {
				return d
			}
		}
		return NetDev{}
	*/
}

//Net struct to contain all network fnformation
type Net struct {
	List map[string]NetDev
}

//Init prepare list of network
func (n *Net) Init() {
	n.Update()
}

//Update refresh list of network
func (n *Net) Update() {
	n.List = readNetDevs(n.List)
}

func readNetDevs(old map[string]NetDev) map[string]NetDev {
	netDevs := map[string]NetDev{}
	devMap, _ := readFileMap([]string{`[\w]+`}, PROC_NET_DATA, `:`)
	for key, value := range devMap {
		vals := regexp.MustCompile(`[\s]+`).Split(strings.TrimSpace(value), -1)
		if len(vals) >= 16 {
			dataReceive := readNetData(vals, 0)
			dataReceiveDiff := netDataDiff(getDev(key, old).Receive, dataReceive)
			dataTransmit := readNetData(vals, 8)
			dataTransmitDiff := netDataDiff(getDev(key, old).Transmit, dataTransmit)

			//netDevs = append(netDevs, NetDev{Name: key, Receive: dataReceive, Transmit: dataTransmit, ReceiveDiff: dataReceiveDiff, TransmitDiff: dataTransmitDiff})
			netDevs[key] = NetDev{Name: key, Receive: dataReceive, Transmit: dataTransmit, ReceiveDiff: dataReceiveDiff, TransmitDiff: dataTransmitDiff}
		}
	}
	return netDevs
}
func readNetData(in []string, off int) NetData {
	netData := NetData{}
	netData.Bytes, _ = strconv.ParseInt(in[0+off], 0, 64)
	netData.Packets, _ = strconv.ParseInt(in[1+off], 0, 64)
	netData.Errs, _ = strconv.ParseInt(in[2+off], 0, 64)
	netData.Drop, _ = strconv.ParseInt(in[3+off], 0, 64)
	netData.Fifo, _ = strconv.ParseInt(in[4+off], 0, 64)
	netData.Frame, _ = strconv.ParseInt(in[5+off], 0, 64)
	netData.Compressed, _ = strconv.ParseInt(in[6+off], 0, 64)
	netData.Multicast, _ = strconv.ParseInt(in[7+off], 0, 64)
	return netData
}
