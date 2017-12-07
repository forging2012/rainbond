// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>

package model

//GetUserToken GetUserToken
//swagger:parameters getToken
type GetUserToken struct {
	// in: body
	Body struct {
		// eid
		// in: body
		// required: true
		EID string `json:"eid" validate:"eid|required"`
		// 可控范围
		// in: body
		// required: false
		Range string `json:"range" validate:"range"`
		// 有效期
		// in: body
		// required: true
		ValidityPeriod int `json:"validity_period" validate:"validity_period|required"` //1549812345
		// 数据中心标识
		// in: body
		// required: true
		RegionTag  string `json:"region_tag" validate:"region_tag"`
		BeforeTime int    `json:"before_time"`
	}
}

//TokenInfo TokenInfo
type TokenInfo struct {
	EID   string `json:"eid"`
	Token string `json:"token"`
	CA    string `json:"ca"`
	Key   string `json:"key"`
}

//APIManager APIManager
//swagger:parameters addAPIManager deleteAPIManager
type APIManager struct {
	//in: body
	Body struct {
		//api级别
		//in: body
		//required: true
		ClassLevel string `json:"class_level" validate:"class_level|reqired"`
		//uri前部
		//in: body
		//required: true
		Prefix string `json:"prefix" validate:"prefix|required"`
		//完整uri
		//in: body
		//required: false
		URI string `json:"uri;size:256" validate:"uri"`
		//别称
		//in: body
		//required: false
		Alias string `json:"alias;size:64" validate:"alias"`
		//补充信息
		//in:body
		//required: false
		Remark string `json:"remark;size:64" validate:"remark"`
	}
}
