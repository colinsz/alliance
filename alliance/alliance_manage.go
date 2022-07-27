package alliance

import (
	"cangku/db"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type AllianceManager struct {
}

func (am *AllianceManager) QueryByPlayer(playername string) (string, error) {
	allianceid, err := getPlayerAlliance(playername)
	if err != nil {
		return allianceid, err
	}

	// 检查下公会是否被解散
	// 如果解散了顺便删除下
	var existAlliance bool
	existAlliance, err = isAllianceList(allianceid)
	if err != nil {
		return allianceid, err
	}
	if !existAlliance {
		delAllianceMember(playername)
		return "", nil
	}

	return allianceid, nil
}

func (am *AllianceManager) Create(id string, playername string) *Alliance {
	// TODO id - check
	if len(id) == 0 {
		return nil
	}

	// 创建
	alliance := createAlliance(id, playername)
	if alliance == nil {
		return nil
	}

	// 加入公会列表
	if !addAllianceList(id) {
		return nil
	}

	// 设置玩家预公会的映射
	joinAlliance(playername, id)

	// 存起来公会信息
	if err := alliance.setAllianceStorageToDB(); err != nil {
		return nil
	}

	return alliance
}

//
func (am *AllianceManager) List(id string) []string {
	return getAllianceList(id)
}

func (am *AllianceManager) Join(playername string, id string) {
	joinAlliance(playername, id)
}

func (am *AllianceManager) Dismiss(id string) {
	// 从公会列表中删除
	delAllianceFromList(id)
	// delAllianceMember(id)

	// TODO 公会dismiss标志置位
}

func createAlliance(id string, playername string) *Alliance {
	if len(id) == 0 {
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

func isAllianceList(id string) (bool, error) {
	p := db.GetRedisPool()
	conn := p.Get()

	c, err := redis.Int(conn.Do("SISMEMBER", allianceListKey(), id))
	// TODO 如果因网络问题可能导致误认为联盟被解散
	if err != nil && c != 1 {
		return false, err
	}
	return true, nil
}

func getAllianceList(id string) []string {
	p := db.GetRedisPool()
	conn := p.Get()

	r := make([]string, 0)
	cursor := "0"
	for {
		ss, err := redis.Values(redis.Strings(conn.Do("sscan", allianceListKey(), cursor)))
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
