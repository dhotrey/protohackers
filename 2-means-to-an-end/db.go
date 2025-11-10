package main

import (
	"slices"
)

type db struct {
	kv   map[int32]int32
	keys []int32
}

func initDb() db {
	return db{
		kv:   map[int32]int32{},
		keys: []int32{},
	}
}

func (d *db) Add(key, value int32) bool {
	_, exists := d.kv[key]
	if exists {
		return false // only one value can exist for a timestamp
	}

	d.kv[key] = value
	d.keys = append(d.keys, key)
	return true
}

func (d *db) Query(minTime, maxTime int32) int32 {
	slices.Sort(d.keys)
	var sum int64
	items := 0

	for _, time := range d.keys {
		if time >= minTime && time <= maxTime {
			sum += int64(d.kv[time])
			items++
		}
		if time > maxTime {
			break // since array is sorted in ascending order we can break as soon as we cross maxTime
		}
	}
	if items == 0 {
		return 0
	}
	return int32(sum / int64(items))
}
