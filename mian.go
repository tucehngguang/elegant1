package main

import (
	"fmt"
	"time"
)

type APN struct {
	name       string
	useage     int  //流量使用量
	Tlimit     int  //流量上限
	Expiration Time //到期时间
}
type Time struct {
	year  int
	month time.Month
	day   int
} //到期时间
type Sim struct {
	iccid      string
	imsi       string
	msisdn     string //三码
	state      string
	useage     int  //流量使用量
	Tlimit     int  //流量上限
	Expiration Time //到期时间12121212666454556454
	apn        [50]APN
} //已修改

func menu() {

	fmt.Println("--------系统介绍-------------")
	fmt.Println("sim卡管理系统")
	fmt.Println("1为状态激活")
	fmt.Println("2是变更卡的流量上限")
	fmt.Println("3是查看卡的使用流量是否达到上限")
	fmt.Println("4是输入卡的到期时间")
	fmt.Println("5为查看卡是否到期")
	fmt.Println("-----------------------------")

}
func activation(S *Sim) {
	if S.state == "未启用" {
		S.state = "激活"
		fmt.Printf("%s", S.state)
		return
	}
	if S.state == "停用" {
		fmt.Printf("sim卡已停用无法激活\n")
		return
	}
	if S.state == "激活" {
		fmt.Printf("sim卡已激活无法再次激活\n")
		return
	}
} //激活
func change(S *Sim, a1 int, a2 int) {

	S.apn[1].Tlimit = a1
	S.apn[2].Tlimit = a2
	S.Tlimit = max(a1, a2)
	fmt.Printf("apn1的流量上限为%dKB\n", a1)
	fmt.Printf("apn2的流量上限为%dKB\n", a2)
	fmt.Printf("sim卡的流量上限为%dKB\n", S.Tlimit)
} //查看流量上限
func detection(S *Sim, ua1 int, ua2 int) {
	fmt.Printf("请输入使用的流量\n")
	S.useage = ua1 + ua2
	S.apn[1].useage = ua1
	S.apn[2].useage = ua2
	if S.apn[1].useage > S.apn[1].Tlimit || S.apn[2].useage > S.apn[2].Tlimit || S.useage > S.Tlimit {
		S.state = "停用"
		fmt.Printf("卡已到期sim卡的状态为%s\n", S.state)

	} else {
		fmt.Printf("还有空余流量\n")

	} //检测是否有空余流量
}
func modifydate(S *Sim, ta1 Time, ta2 Time) {
	fmt.Printf("请输入卡的到期时间\n")
	fmt.Scan(&ta1.year, &ta1.month, &ta1.day)
	fmt.Scan(&ta2.year, &ta2.month, &ta2.day)
	S.apn[1].Expiration = ta1
	S.apn[2].Expiration = ta2
	time1 := time.Date(ta1.year, ta1.month, ta1.day, 15, 10, 11, 999000111, time.Local)
	time2 := time.Date(ta2.year, ta2.month, ta2.day, 15, 10, 11, 999000111, time.Local)
	if time1.After(time2) {
		S.Expiration = ta1
		fmt.Printf("卡的到期时间:%d-%d-%d %d:%d:%d %s\n", time1.Year(), time1.Month(), time1.Day(), time1.Hour(), time1.Minute(), time1.Second(), time1.Weekday().String())
	} else {
		S.Expiration = ta2
		fmt.Printf("卡的到期时间:%d-%d-%d %d:%d:%d %s\n", time2.Year(), time2.Month(), time2.Day(), time2.Hour(), time2.Minute(), time2.Second(), time2.Weekday().String())
	}
} //修改卡到期时间
func checkepire(S *Sim, ta1 Time) {
	fmt.Printf("请输入修改的卡到期时间\n")
	fmt.Scan(&ta1.year, &ta1.month, &ta1.day)
	time1 := time.Date(ta1.year, ta1.month, ta1.day, 15, 10, 11, 999000111, time.Local)
	time2 := time.Date(S.Expiration.year, S.Expiration.month, S.Expiration.day, 15, 10, 11, 999000111, time.Local)
	if time1.After(time2) {
		S.state = "停用"
		fmt.Printf("卡已到期sim卡的状态为%s\n", S.state)
	} else {
		fmt.Printf("sim卡未到期还能继续使用\n")
	}
}
func printsim(S *Sim) {
	fmt.Printf("Sim iccid为: %v\n", S.iccid)
	fmt.Printf("Sim imsi为: %v\n", S.imsi)
	fmt.Printf("Sim msisdn为: %v\n", S.msisdn)
	fmt.Printf("Sim 状态为: %v\n", S.state)
	fmt.Printf("Sim 流量使用量为: %v\n", S.useage)
	fmt.Printf("Sim 流量上限为: %v\n", S.Tlimit)
	fmt.Printf("Sim 到期时间为: %#v\n", S.Expiration)
}
func main() {
	var choice int
	var la1 int
	var la2 int //流量上限
	var ua1 int
	var ua2 int //使用流量
	var sim Sim
	var ta1 Time
	var ta2 Time //到期时间
	sim.iccid = "xx"
	sim.imsi = "xxx"
	sim.msisdn = "xxxx"
	sim.state = "激活" //卡初始状态
	sim.apn[1].name = "apn1"
	sim.apn[2].name = "apn2" //apn名字
	sim.apn[1].Tlimit = 6
	sim.apn[2].Tlimit = 4
	sim.Tlimit = 6 //截止时间
	sim.Expiration.year = 2024
	sim.Expiration.month = 6
	sim.Expiration.day = 30
	sim.apn[1].Expiration.year = 2024
	sim.apn[1].Expiration.month = 6
	sim.apn[1].Expiration.day = 30
	sim.apn[2].Expiration.year = 2024
	sim.apn[2].Expiration.month = 5
	sim.apn[2].Expiration.day = 1
	//卡的截止时间
	menu() //菜单打印
	for {
		fmt.Println("是否开始操作")
		fmt.Println("1为状态激活")
		fmt.Println("2是变更卡的流量上限")
		fmt.Println("3是查看卡的使用流量是否达到上限")
		fmt.Println("4是输入卡的到期时间")
		fmt.Println("5为查看卡是否到期")
		fmt.Println("-----------------------")
		fmt.Scan(&choice)
		switch choice {
		case 0:
			{
				fmt.Printf("已操作结束\n")
				return

			}
		case 1:
			{
				activation(&sim)

			} //激活卡

		case 2:
			{
				fmt.Printf("请输入流量上限\n")
				fmt.Scan(&la1, &la2)
				change(&sim, la1, la2)

			} //变更卡的流量上限
		case 3:
			{
				fmt.Printf("请输入已使用流量\n")
				fmt.Scan(&ua1, &ua2)
				detection(&sim, ua1, ua2)

			} //查看卡的使用流量是否达到上限
		case 4:
			{
				modifydate(&sim, ta1, ta2)

			}
		case 5:
			{
				checkepire(&sim, ta1)

			}
		case 6:
			{
				printsim(&sim)
			}
		}

	}
}
