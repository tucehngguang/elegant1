package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Time struct {
	year  int
	month time.Month
	day   int
} //到期时间
type Sim struct {
	ID         uint      `gorm:"primary_key"`
	ICCID      string    `gorm:"type:varchar(12)"`
	IMSI       string    `gorm:"type:varchar(12)"`
	MSISDN     string    `gorm:"type:varchar(12)"`
	STATE      string    `gorm:"type:varchar(12)"`
	USEAGE     int       //流量使用量
	Tlimit     int       //流量上限
	Expiration time.Time `gorm:"type:datetime"` //到期时间12121212666454556454
	APN        []APN     `gorm:"many2many:sim_apn;constraint:OnUpdate:CASCADE,OnDelete:CASCADE; "`
}

type APN struct {
	ID         uint      `gorm:"primary_key"`
	ICCID      string    `gorm:"type:varchar(12)"`
	NAME       string    `gorm:"type:varchar(12)"`
	USEAGE     int       //流量使用量
	Tlimit     int       //流量上限
	Expiration time.Time `gorm:"type:datetime"` //到期时间
	Sims       []Sim     `gorm:"many2many: sim_apn;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func printsql(dsn string) {

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败，error=" + err.Error())
	}
	fmt.Printf("请输要查询的卡号\n")
	var iccid string
	var sim Sim
	var apn APN
	var apn1 APN
	fmt.Scan(&iccid)
	if result := db.First(&sim, "icc_id = ?", iccid); result.Error != nil {
		log.Fatal("Error finding product:", result.Error)
	}
	fmt.Printf("SIM: ICCID: %s, IMSI: %s, MSISDN: %s, STATE: %s, USEAGE: %d, Tlimit: %d, Expiration: %v\n",
		sim.ICCID, sim.IMSI, sim.MSISDN, sim.STATE, sim.USEAGE, sim.Tlimit, sim.Expiration)
	if result := db.Where("icc_id = ? AND name = ?", iccid, "apn1").First(&apn); result.Error != nil {
		log.Fatal("Error finding product:", result.Error)
	}
	fmt.Printf("APN1 Details: ICCID: %s Name: %s Usage: %dLimit: %dExpiration: %s\n",
		apn.ICCID, apn.NAME, apn.USEAGE, apn.Tlimit, apn.Expiration)
	if result := db.Where("icc_id = ? AND name = ?", iccid, "apn2").First(&apn1); result.Error != nil {
		log.Fatal("Error finding product:", result.Error)
	}
	fmt.Printf("APN2 Details: ICCID: %s Name: %s Usage: %dLimit: %dExpiration: %s\n",
		apn1.ICCID, apn1.NAME, apn1.USEAGE, apn1.Tlimit, apn1.Expiration)
}

// 已修改
func creatsql(dsn string) {
	var i int
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败，error=" + err.Error())
	}
	for i = 1; i < 1001; i++ {

		// 将字符串转换为整数

		// 将整数加1

		// 将结果转换为字符串
		result := strconv.Itoa(i)
		sim3 := Sim{
			ICCID:      result,
			IMSI:       "31",
			MSISDN:     "tcg",
			STATE:      "ac",
			USEAGE:     150,
			Tlimit:     1000,
			Expiration: time.Now().AddDate(1, 0, 0), // 假设一年后到期
			APN: []APN{
				{
					ICCID:      result,
					NAME:       "apn1",
					USEAGE:     100,
					Tlimit:     500,
					Expiration: time.Now().AddDate(1, 0, 0),
				},
				{
					ICCID:      result,
					NAME:       "apn2",
					USEAGE:     100,
					Tlimit:     500,
					Expiration: time.Now().AddDate(1, 0, 0),
				},
			},
		}
		var err error
		err = db.Create(&sim3).Error
		fmt.Println(err) //创建
	}
}
func menu() {

	fmt.Printf("--------系统介绍-------------\n")
	fmt.Printf("sim卡管理系统\n")
	fmt.Printf("1为状态激活\n")
	fmt.Printf("2是变更卡的流量上限\n")
	fmt.Printf("3是查看卡的使用流量是否达到上限\n")
	fmt.Printf("4是输入卡的到期时间\n")
	fmt.Println("5为查看卡是否到期")
	fmt.Println("6为打印sim卡信息")
	fmt.Println("-----------------------------")
}
func activation(S *Sim) {

	if S.STATE == "未启用" {
		S.STATE = "激活"
		fmt.Printf("%s", S.STATE)
		return
	}
	if S.STATE == "停用" {
		fmt.Printf("sim卡已停用无法激活\n")
		return
	}
	if S.STATE == "激活" {
		fmt.Printf("sim卡已激活无法再次激活\n")
		return
	}
} //激活
func change(S *Sim, a1 int, a2 int, client *redis.Client) {
	hashKey := "SIM1"
	a := []int{a1, a2}

	for i := range S.APN {
		S.APN[i].Tlimit = a[i]

	}
	// maxDataLimit := 0
	// // for _, apn := range S.APN {
	// // 	if apn.Tlimit > maxDataLimit {
	// // 		maxDataLimit = apn.Tlimit
	// // 	}
	// // }
	S.Tlimit = max(a1, a2)
	judge, erro := client.HExists(context.Background(), hashKey, "Tlimit").Result()
	if erro != nil {
		fmt.Println("不存在")
	} //判断是否存在该卡的流量上限
	if judge == true {
		val, err := client.HGet(context.Background(), hashKey, "Tlimit").Result()
		if err != nil {
			panic(err)
		}
		i, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("not")
		} //转换成int
		if S.Tlimit > i {
			m := strconv.Itoa(S.Tlimit)
			err = client.HSet(context.Background(), hashKey, "Tlimit", m).Err()
			if err != nil {
				fmt.Printf("111111")
			} else {
				fmt.Printf("")
			}

		}
		fmt.Println("exist")
	} else {
		fmt.Println("not")
	} //判断是否存在改卡的流量上限
	fmt.Printf("apn1的流量上限为%dKB\n", a1)
	fmt.Printf("apn2的流量上限为%dKB\n", a2)
	fmt.Printf("sim卡的流量上限为%dKB\n", S.Tlimit)
} //查看流量上限
func detection(ua1 int, ua2 int, dsn string) {

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败，error=" + err.Error())
	}
	fmt.Printf("请输入修改的卡号\n")
	var iccid string
	var sim Sim
	var apn APN
	var apn1 APN
	fmt.Scan(&iccid)
	if result := db.First(&sim, "icc_id = ?", iccid); result.Error != nil {
		log.Fatal("Error finding product:", result.Error)
	}
	if result := db.Where("icc_id = ? AND name = ?", iccid, "apn1").First(&apn); result.Error != nil {
		log.Fatal("Error finding product:", result.Error)
	}

	if result := db.Where("icc_id = ? AND name = ?", iccid, "apn2").First(&apn1); result.Error != nil {
		log.Fatal("Error finding product:", result.Error)
	}

	sim.USEAGE = ua1 + ua2
	apn.USEAGE = ua1
	apn1.USEAGE = ua2
	if apn1.USEAGE > apn1.Tlimit || apn.USEAGE > apn.Tlimit || sim.USEAGE > sim.Tlimit {
		sim.STATE = "停用"
		result := db.Model(&Sim{}).Where("icc_id = ?", sim.ICCID).Updates(Sim{
			USEAGE: ua1 + ua2,
			STATE:  "停用",
		})
		if result.Error != nil {
			fmt.Println("连接数据库成功")
		}
		result2 := db.Model(&APN{}).Where("name= ?", "apn1").Where("icc_id = ?", sim.ICCID).Updates(APN{
			USEAGE:     ua1,
			Expiration: time.Now().AddDate(2, 0, 0),
		})

		if result2.Error != nil {
			fmt.Println("连接数据库成功")
		}
		result3 := db.Model(&APN{}).Where("name= ?", "apn2").Where("icc_id = ?", sim.ICCID).Updates(APN{
			USEAGE:     ua2,
			Expiration: time.Now().AddDate(2, 0, 0),
		})

		if result3.Error != nil {
			fmt.Println("连接数据库成功")
		}
		conn, err := rabbitmq.NewConn(
			"amqp://guest:guest@localhost",
			rabbitmq.WithConnectionOptionsLogging,
		)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		publisher, err := rabbitmq.NewPublisher(
			conn,
			rabbitmq.WithPublisherOptionsLogging,
			rabbitmq.WithPublisherOptionsExchangeName("events"),
			rabbitmq.WithPublisherOptionsExchangeDeclare,
		)
		if err != nil {
			log.Fatal(err)
		}
		defer publisher.Close()

		publisher.NotifyReturn(func(r rabbitmq.Return) {
			log.Printf("message returned from server: %s", string(r.Body))
		})

		publisher.NotifyPublish(func(c rabbitmq.Confirmation) {
			log.Printf("message confirmed from server. tag: %v, ack: %v", c.DeliveryTag, c.Ack)
		})

		// block main thread - wait for shutdown signal
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			fmt.Println()
			fmt.Println(sig)
			done <- true
		}()
		s := sim.ICCID + sim.STATE
		fmt.Println("awaiting signal")
		ticker := time.NewTicker(time.Second)
		select {
		case <-ticker.C:
			err = publisher.PublishWithContext(
				context.Background(),
				[]byte(s),
				[]string{"my_routing_key"},
				rabbitmq.WithPublishOptionsContentType("application/json"),
				rabbitmq.WithPublishOptionsMandatory,
				rabbitmq.WithPublishOptionsPersistentDelivery,
				rabbitmq.WithPublishOptionsExchange("events"),
			)
			if err != nil {
				log.Println(err)
			}
		case <-done:
			fmt.Println("stopping publisher")
			return
		}

		fmt.Printf("卡已到期sim卡的状态为%s\n", sim.STATE)
	} else {
		fmt.Printf("还有空余流量\n")

	} //检测是否有空余流量
}
func modifydate(S *Sim, ta1 Time, ta2 Time) {
	fmt.Printf("请输入卡的到期时间\n")
	fmt.Scan(&ta1.year, &ta1.month, &ta1.day)
	fmt.Scan(&ta2.year, &ta2.month, &ta2.day)

	time1 := time.Date(ta1.year, ta1.month, ta1.day, 15, 10, 11, 999000111, time.Local)
	time2 := time.Date(ta2.year, ta2.month, ta2.day, 15, 10, 11, 999000111, time.Local)
	S.APN[0].Expiration = time1
	S.APN[1].Expiration = time2
	if time1.After(time2) {
		S.Expiration = time1
		fmt.Printf("卡的到期时间:%d-%d-%d %d:%d:%d %s\n", time1.Year(), time1.Month(), time1.Day(), time1.Hour(), time1.Minute(), time1.Second(), time1.Weekday().String())
	} else {
		S.Expiration = time2
		fmt.Printf("卡的到期时间:%d-%d-%d %d:%d:%d %s\n", time2.Year(), time2.Month(), time2.Day(), time2.Hour(), time2.Minute(), time2.Second(), time2.Weekday().String())
	}
} //修改卡到期时间
func checkepire(S *Sim, ta1 Time) {
	fmt.Printf("请输入修改的卡到期时间\n")
	fmt.Scan(&ta1.year, &ta1.month, &ta1.day)
	time1 := time.Date(ta1.year, ta1.month, ta1.day, 15, 10, 11, 999000111, time.Local)

	if time1.After(S.Expiration) {
		S.STATE = "停用"
		fmt.Printf("卡已到期sim卡的状态为%s\n", S.STATE)
	} else {
		fmt.Printf("sim卡未到期还能继续使用\n")
	}
}
func printsim(S *Sim) {
	fmt.Printf("Sim iccid为: %v\n", S.ICCID)
	fmt.Printf("Sim imsi为: %v\n", S.IMSI)
	fmt.Printf("Sim msisdn为: %v\n", S.MSISDN)
	fmt.Printf("Sim 状态为: %v\n", S.STATE)
	fmt.Printf("Sim 流量使用量为: %v\n", S.USEAGE)
	fmt.Printf("Sim 流量上限为: %v\n", S.Tlimit)
	fmt.Printf("Sim 到期时间为: %#v\n", S.Expiration)
}

// 定义一个全局变量db，用于后面数据库的读写操作,通常就放在全局里面
var DB *gorm.DB

func main() {
	var choice int
	var la1 int
	var la2 int //流量上限
	var ua1 int
	var ua2 int //使用流量

	var ta1 Time
	var ta2 Time        //到期时间
	username := "root"  //账号
	password := "123"   //密码
	host := "localhost" //数据库地址
	port := "3306"      //端口
	Dnname := "study"   //数据库名

	//root:root@tcp(127.0.0.1:3306)/test？
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, Dnname)
	//连接mysql，获得DB类型实例，用于后面数据库的读写操作
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败，error=" + err.Error())
	}
	// 创建Sim4实例并关联APN

	sim := []Sim{
		{ICCID: "210",
			IMSI:       "31",
			MSISDN:     "123",
			STATE:      "ac",
			USEAGE:     2,
			Tlimit:     6,
			Expiration: time.Now().AddDate(1, 0, 0), // 假设一年后到期
			APN: []APN{
				{
					ICCID:      "210",
					NAME:       "apn1",
					USEAGE:     2,
					Tlimit:     6,
					Expiration: time.Now().AddDate(1, 0, 0),
				},
				{
					ICCID:      "210",
					NAME:       "apn2",
					USEAGE:     3,
					Tlimit:     6,
					Expiration: time.Now().AddDate(1, 0, 0),
				},
			},
		},
		{ICCID: "21",
			IMSI:       "31",
			MSISDN:     "123",
			STATE:      "ac",
			USEAGE:     150,
			Tlimit:     1000,
			Expiration: time.Now().AddDate(1, 0, 0), // 假设一年后到期
			APN: []APN{
				{
					ICCID:      "210",
					NAME:       "apn1",
					USEAGE:     100,
					Tlimit:     500,
					Expiration: time.Now().AddDate(1, 0, 0),
				},
				{
					ICCID:      "210",
					NAME:       "apn2",
					USEAGE:     100,
					Tlimit:     500,
					Expiration: time.Now().AddDate(1, 0, 0),
				},
			},
		},
	}

	// 输出结果以验证
	fmt.Printf("Sim4 ICCID: %s, APNs: %d\n", sim[0].ICCID, len(sim[0].APN))
	for _, apn := range sim[1].APN {
		fmt.Printf("APN NAME: %s, USEAGE: %d\n", apn.NAME, apn.USEAGE)
	}

	DB = db
	DB.AutoMigrate(&Sim{})
	err = db.Create(&sim[0]).Error
	fmt.Println(err) //创建
	menu()           //菜单打印
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.24.225:6379", // Redis服务器地址
		Password: "",                    // 可选：如果有密码的话
		DB:       0,                     // 使用的数据库编号
	}) //redis数据库引入

	// 使用Ping测试连接
	hashKey := "SIM1"

	// 使用HMSet方法设置多个哈希字段及其对应的值
	fieldsAndValues := map[string]interface{}{
		"ICCID":      "xx",
		"IMSI":       "xxx",
		"MSISDN":     "xxxx",
		"State":      "未启用",
		"Tlimit":     "6",
		"Expiration": "2024-3-2",
	}

	erro := client.HMSet(context.Background(), hashKey, fieldsAndValues).Err()
	if erro != nil {

		fmt.Println("Error setting hash fields:", err)
		return
	}

	for {
		fmt.Println("是否开始操作")
		fmt.Printf("1为状态激活\n")
		fmt.Printf("2是变更卡的流量上限\n")
		fmt.Printf("3是查看卡的使用流量是否达到上限\n")
		fmt.Printf("4是输入卡的到期时间\n")
		fmt.Printf("5为查看卡是否到期\n")
		fmt.Printf("6是打印表的信息n")
		fmt.Printf("7为创建一千张表\n")
		fmt.Printf("8为输出指定表的信息\n")
		fmt.Scan(&choice)
		switch choice {
		case 0:
			{
				fmt.Printf("已操作结束\n")
				return

			}
		case 1:
			{
				activation(&sim[0])

			} //激活卡

		case 2:
			{
				fmt.Printf("请输入流量上限\n")
				fmt.Scan(&la1, &la2)
				change(&sim[0], la1, la2, client)

			} //变更卡的流量上限
		case 3:
			{
				fmt.Printf("请输入已使用流量\n")
				fmt.Scan(&ua1, &ua2)
				detection(ua1, ua2, dsn)

			} //查看卡的使用流量是否达到上限
		case 4:
			{
				modifydate(&sim[0], ta1, ta2)

			}
		case 5:
			{
				checkepire(&sim[0], ta1)

			}
		case 6:
			{
				printsim(&sim[0])
			}
		case 7:
			{

				creatsql(dsn) //创建一千张表
			}
		case 8:
			{

				printsql(dsn) //打印表
			}
		}

	}

}
