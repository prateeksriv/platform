// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	"io"
)

type ChannelList []*Channel

func (o *ChannelList) ToJson() string {
	if b, err := json.Marshal(o); err != nil {
		return "[]"
	} else {
		return string(b)
	}
}

func (o *ChannelList) Etag() string {

	id := "0"
	var t int64 = 0
	var delta int64 = 0

	for _, v := range *o {
		if v.LastPostAt > t {
			t = v.LastPostAt
			id = v.Id
		}

		if v.UpdateAt > t {
			t = v.UpdateAt
			id = v.Id
		}

	}

	return Etag(id, t, delta, len(*o))
}

// This function was buggy and I had to create the following workaround to make it work.
// I know this is not the correct solution, I just wanted to share this idea or solution.
// I do not know enough about golang coding to design a fix. My intent was to highlight the bug.

func ChannelListFromJson(data io.Reader) *mattermost.ChannelList {

	var ret mattermost.ChannelList
	decoder := json.NewDecoder(data)
	var dat map[string]interface{}
	err := decoder.Decode(&dat)
	cli := dat["channels"].([]interface{})
	for _, ch := range cli {
		var t mattermost.Channel
		for k, v := range ch.(map[string]interface{}) {
			if k == "name" {
				t.Name = v.(string)
			}
			if k == "display_name" {
				t.DisplayName = v.(string)
			}
			if k == "purpose" {
				t.Purpose = v.(string)
			}
			if k == "id" {
				t.Id = v.(string)
			}
		}
		ret = append(ret, &t)
	}
	if err == nil {
		return &ret
	} else {
		return nil
	}
}
