package linux

import (
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type Sockstat struct {
	// sockets:
	SocketsUsed uint64 `json:"sockets_used" altname:"sockets:used"`

	// TCP:
	TCPInUse     uint64 `json:"tcp_in_use" altname:"TCP:inuse"`
	TCPOrphan    uint64 `json:"tcp_orphan" altname:"TCP:orphan"`
	TCPTimeWait  uint64 `json:"tcp_time_wait" altname:"TCP:tw"`
	TCPAllocated uint64 `json:"tcp_allocated" altname:"TCP:alloc"`
	TCPMemory    uint64 `json:"tcp_memory" altname:"TCP:mem"`

	// UDP:
	UDPInUse  uint64 `json:"udp_in_use" altname:"UDP:inuse"`
	UDPMemory uint64 `json:"udp_memory" altname:"UDP:mem"`

	// UDPLITE:
	UDPLITEInUse uint64 `json:"udplite_in_use" altname:"UDPLITE:inuse"`

	// RAW:
	RAWInUse uint64 `json:"raw_in_use" altname:"RAW:inuse"`

	// FRAG:
	FRAGInUse  uint64 `json:"frag_in_use" altname:"FRAG:inuse"`
	FRAGMemory uint64 `json:"frag_memory" altname:"FRAG:memory"`
}

func ReadSockstat(path string) (*Sockstat, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	// Maps a meminfo metric to its value (i.e. MemTotal --> 100000)
	statMap := make(map[string]uint64)

	var sockstat Sockstat = Sockstat{}

	for _, line := range lines {
		if strings.Index(line, ":") == -1 {
			continue
		}

		statType := line[0 : strings.Index(line, ":")+1]
		stats := pair(strings.Fields(line[strings.Index(line, ":")+1:]))

		if len(stats) < 1 {
			continue
		}

		for _, pair := range stats {
			val, _ := strconv.ParseUint(pair._2, 10, 64)
			statMap[statType+pair._1] = val
		}
	}

	elem := reflect.ValueOf(&sockstat).Elem()
	typeOfElem := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		val, ok := statMap[typeOfElem.Field(i).Tag.Get("altname")]
		if ok {
			elem.Field(i).SetUint(val)
		}
	}

	return &sockstat, nil
}

type StrPair struct {
	_1 string
	_2 string
}

func pair(strs []string) []StrPair {
	pairs := make([]StrPair, 0)

	for i := 1; i < len(strs); i = i + 2 {
		pairs = append(pairs, StrPair{strs[i-1], strs[i]})
	}

	return pairs
}
