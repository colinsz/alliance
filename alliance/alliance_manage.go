package alliance

import (
	"errors"
	"fmt"
	"log"

	"github.com/colinsz/alliance/db"

	"github.com/gomodule/redigo/redis"
)

var AllianceExist = errors.New("alliance: exist")
var AllianceNameErr = errors.New("alliance: id len 0")
var RedisNetworkErr = errors.New("alliance: redis error")

type AllianceManager struct {
}

func (am *AllianceManager) QueryByPlayer(playername string) (string, error) {
	allianceid, err := getPlayerAlliance(playername)
	if err != nil {
		return allianceid, RedisNetworkErr
	}

	// 检查下公会是否被解散
	var existAlliance int
	existAlliance, err = isAllianceList(allianceid)
	if err != nil {
		return allianceid, RedisNetworkErr
	}
	// 如果解散了顺便删除下
	if existAlliance == 0 {
		delAllianceMember(playername)
		allianceid = ""
	}

	return allianceid, nil
}

func (am *AllianceManager) Create(id string, playername string) (*Alliance, error) {
	// 创建
	alliance := createAlliance(id, playername)
	if alliance == nil {
		log.Println("alliance: createAlliance failed")
		return nil, AllianceNameErr
	}

	// 加入公会列表
	if !addAllianceList(id) {
		log.Println("alliance: addAllianceList failed")
		return nil, AllianceExist
	}

	// 设置玩家预公会的映射
	joinAlliance(playername, id)

	// 存起来公会信息
	if err := alliance.setAllianceStorageToDB(); err != nil {
		log.Println("alliance: setAllianceStorageToDB failed")
		return nil, RedisNetworkErr
	}

	return alliance, nil
}

//
func (am *AllianceManager) List() []string {
	return getAllianceList()
}

func (am *AllianceManager) Join(playername string, id string) {
	joinAlliance(playername, id)
}

func (am *AllianceManager) Dismiss(allianceid string, playername string) error {
	if err := lockAlliance(allianceid); err != nil {
		return LockErr
	}
	defer func() { unlockAlliance(allianceid) }()

	as := &Alliance{}
	if err := as.getAllianceStorage(allianceid); err != nil {
		return err
	}

	// 只允许会长才可以解散
	if playername == "" || as.OwnerId != playername {
		return NotAllianceOwner
	}
	as.Dismiss = true

	// 从公会列表中删除
	delAllianceFromList(allianceid)
	// delAllianceMember(id)

	// 公会dismiss标志置位
	return as.setAllianceStorageToDB()
}

func createAlliance(id string, playername string) *Alliance {
	if len(id) == 0 {
		log.Println("alliance: id = 0")
		return nil
	}

	al := &Alliance{
		AllianceId:  id,
		OwnerId:     playername,
		Items:       make([]ItemDesc, MaxAllianceGrid),
		MaxCapacity: MaxAllianceGrid,
		Dismiss:     false,
	}
	for i := 0; i < len(al.Items); i++ {
		al.Items[i].ItemType = EmptyGridItemType
	}

	return al
}

func allianceListKey() string {
	return "alliance:list"
}

func addAllianceList(id string) bool {
	p := db.GetRedisPool()
	conn := p.Get()

	c, err := redis.Int64(conn.Do("sadd", allianceListKey(), id))
	if err != nil || c == 0 {
		log.Printf("err:%+v, c:%+v\n", err, c)
		return false
	}
	return true
}

func delAllianceFromList(id string) bool {
	p := db.GetRedisPool()
	conn := p.Get()

	c, err := redis.Int64(conn.Do("srem", allianceListKey(), id))
	if err != nil || c == 0 {
		return false
	}
	return true
}

func isAllianceList(id string) (int, error) {
	p := db.GetRedisPool()
	conn := p.Get()

	return redis.Int(conn.Do("SISMEMBER", allianceListKey(), id))
}

func getAllianceList() []string {
	p := db.GetRedisPool()
	conn := p.Get()

	r := make([]string, 0)
	cursor := "0"
	for {
		ss, err := redis.Values(conn.Do("sscan", allianceListKey(), cursor))
		if err != nil || len(ss) != 2 {
			return r
		}
		vv := ss[1].([]interface{})
		for i := 0; i < len(vv); i++ {
			r = append(r, string(vv[i].([]byte)))
		}
		cursor = string(ss[0].([]byte))
		if cursor == "0" {
			break
		}
	}

	return r
}

func allianceMemberKey(playername string) string {
	return fmt.Sprintf("alliance:member:%s", playername)
}

func joinAlliance(playername string, id string) {
	p := db.GetRedisPool()
	conn := p.Get()

	conn.Do("set", allianceMemberKey(playername), id)
}

func getPlayerAlliance(playername string) (string, error) {
	p := db.GetRedisPool()
	conn := p.Get()

	return redis.String(conn.Do("GET", allianceMemberKey(playername)))
}

func delAllianceMember(playername string) {
	p := db.GetRedisPool()
	conn := p.Get()

	conn.Do("del", allianceMemberKey(playername))
}
