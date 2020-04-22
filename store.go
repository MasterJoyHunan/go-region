package main

type Region struct {
	Id       string // 省市区编码
	ParentId string // 所属上级 0:顶级
	Name     string // 省市区名
}
