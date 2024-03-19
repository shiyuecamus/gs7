// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package util

import (
	"math"
)

type ComItem struct {
	// Index 数据索引
	Index int
	// Index 原始数据大小
	RawSize int
	// SplitOffset 分割点，即数据偏移索引
	SplitOffset int
	// RipeSize 分割后的数据大小
	RipeSize int
	// ExtraSize 额外需要的数据大小
	ExtraSize int
	// Threshold 阀值
	Threshold int
}

func (i *ComItem) TotalLength() int {
	return int(math.Max(float64(i.RipeSize+i.ExtraSize), float64(i.Threshold)))
}

type ComGroup struct {
	Items []*ComItem
}

func WriteRecombination(src []uint16, targetSize int, extraSize int) []*ComGroup {
	sum := 0
	var number int
	group := &ComGroup{
		Items: make([]*ComItem, 0),
	}
	groupList := []*ComGroup{
		group,
	}
	for i := 0; i < len(src); i++ {
		number = int(src[i])
		offset := 0
		for number > 0 {
			item := &ComItem{
				Index:       i,
				RawSize:     int(src[i]),
				SplitOffset: offset,
				RipeSize:    0,
				ExtraSize:   extraSize,
				Threshold:   0,
			}
			if sum+number+item.ExtraSize > targetSize {
				item.RipeSize = targetSize - sum - item.ExtraSize
			} else {
				item.RipeSize = number
			}
			sum += item.TotalLength()
			number -= item.RipeSize
			offset += item.RipeSize
			group.Items = append(group.Items, item)
			if sum+extraSize >= targetSize {
				group = &ComGroup{
					Items: make([]*ComItem, 0),
				}
				groupList = append(groupList, group)
				sum = 0
			}
		}
	}
	return Filter[*ComGroup](groupList, func(v *ComGroup) bool {
		return len(v.Items) != 0
	})
}

func ReadRecombination(src []uint16, targetSize int, extraSize int, threshold int) []*ComGroup {
	sum := 0
	var number int
	group := &ComGroup{
		Items: make([]*ComItem, 0),
	}
	groupList := []*ComGroup{
		group,
	}

	for i := 0; i < len(src); i++ {
		number = int(src[i])
		offset := 0
		for number > 0 {
			item := &ComItem{
				Index:       i,
				RawSize:     int(src[i]),
				SplitOffset: offset,
				RipeSize:    0,
				ExtraSize:   extraSize,
				Threshold:   threshold,
			}
			if sum+number+item.ExtraSize > targetSize {
				item.RipeSize = targetSize - sum - item.ExtraSize
			} else {
				item.RipeSize = number
			}
			sum += item.TotalLength()
			number -= item.RipeSize
			offset += item.RipeSize
			group.Items = append(group.Items, item)
			if sum+threshold >= targetSize {
				group = &ComGroup{
					Items: make([]*ComItem, 0),
				}
				groupList = append(groupList, group)
				sum = 0
			}
		}
	}
	return Filter[*ComGroup](groupList, func(v *ComGroup) bool {
		return len(v.Items) != 0
	})
}
