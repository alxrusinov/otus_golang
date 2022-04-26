package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type StringFrequency struct {
	Value     string
	Frequency int
}

func prepareString(str string) []string {
	return strings.Fields(str)
}

func getFrequencyMap(sl []string) map[string]*StringFrequency {
	res := make(map[string]*StringFrequency)

	for _, str := range sl {
		if _, ok := res[str]; ok {
			res[str].Frequency++
			continue
		}
		res[str] = &StringFrequency{
			Value:     str,
			Frequency: 1,
		}
	}
	return res
}

func getFrequency(mp map[string]*StringFrequency) []*StringFrequency {
	res := []*StringFrequency{}

	for key := range mp {
		res = append(res, mp[key])
	}
	sort.Slice(res, func(prev, next int) bool {
		if res[prev].Frequency == res[next].Frequency {
			return res[prev].Value < res[next].Value
		}
		return res[prev].Frequency > res[next].Frequency
	})

	return res
}

func takeFirst(sl []*StringFrequency, num int) []string {
	res := []string{}

	for key := range sl {
		if key+1 > num {
			break
		}
		res = append(res, sl[key].Value)
	}

	return res
}

func Top10(str string) []string {
	sl := prepareString(str)
	mp := getFrequencyMap(sl)
	fullRes := getFrequency(mp)
	res := takeFirst(fullRes, 10)
	return res
}
