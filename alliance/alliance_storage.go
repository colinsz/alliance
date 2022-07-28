package alliance

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/colinsz/alliance/db"

	"github.com/gomodule/redigo/redis"
)

var ItemTypeErr = errors.New("alliance: item type error")
var ItemCountErr = errors.New("alliance: item count error")
var AllianceIsFull = errors.New("alliance: alliance full")
var GridIndexErr = errors.New("alliance: grid index error")
var HasIncreased = errors.New("alliance: has increased")
var LockErr = errors.New("alliance: lock error")
var NotAllianceOwner = errors.New("alliance: not alliance owner")
var UnexpectedItermNumber = errors.New("alliance: unexpect item number")
var AllianceDelete = errors.New("alliance: deleted")

const MaxGridStack = 5
const MaxAllianceGrid = 30

type Alliance struct {
	AllianceId  string     // 公会仓库ID 全局唯一
	OwnerId     string     // 公会仓库会长id
	Items       []ItemDesc // 公会道具
	MaxCapacity int        // 默认为30
	Dismiss     bool       // 是否已经解散
}

const EmptyGridItemType = 0 // 表示公会格子没有道具

type ItemDesc struct {
	// Name     string // 道具名
	ItemType int32 // 道具类型
	Number   int32 // 道具数量
	// Dirty    int32 // 是否更新
}

// 会长修改公会仓库容量
func (as *Alliance) AllianceStorageIncreaseCapacity(allianceid string, count int) error {
	if err := lockAlliance(allianceid); err != nil {
		return LockErr
	}
	defer func() { unlockAlliance(allianceid) }()

	if err := as.getAllianceStorage(allianceid); err != nil {
		return err
	}

	if as.MaxCapacity >= 40 {
		return HasIncreased
	}

	// 修改容量
	as.MaxCapacity += count
	increaseCap := make([]ItemDesc, count)
	for i := 0; i < len(increaseCap); i++ {
		increaseCap[i].ItemType = EmptyGridItemType
	}
	as.Items = append(as.Items, increaseCap...)

	if err := as.setAllianceStorageToDB(); err != nil {
		return RedisNetworkErr
	}

	return nil
}

// 会长销毁公会仓库的某格物品
func (as *Alliance) AllianceStorageDestroyItem(allianceid, playername string, index int) error {
	if err := lockAlliance(allianceid); err != nil {
		return LockErr
	}
	defer func() { unlockAlliance(allianceid) }()

	if err := as.getAllianceStorage(allianceid); err != nil {
		return err
	}

	// 只允许会长可以销毁某一格的所有物品
	if playername == "" || as.OwnerId != playername {
		return NotAllianceOwner
	}

	if err := as.destroyItem(index); err != nil {
		return RedisNetworkErr
	}

	return as.setAllianceStorageToDB()
}

// 整理公会仓库
func (as *Alliance) AllianceStorageClearup(allianceid string) error {
	if err := lockAlliance(allianceid); err != nil {
		return LockErr
	}
	defer func() { unlockAlliance(allianceid) }()

	if err := as.getAllianceStorage(allianceid); err != nil {
		return err
	}

	if err := as.clearup(); err != nil {
		return err
	}

	return as.setAllianceStorageToDB()
}

// 普通玩家向公会仓库新增物品
func (as *Alliance) AllianceStorageAddItem(allianceid string, index int, itemtype int32, num int32) error {
	if err := lockAlliance(allianceid); err != nil {
		return LockErr
	}
	defer func() { unlockAlliance(allianceid) }()

	if err := as.getAllianceStorage(allianceid); err != nil {
		return err
	}

	if err := as.add(index, itemtype, num); err != nil {
		return err
	}

	return as.setAllianceStorageToDB()
}

func (as *Alliance) clearup() error {
	// 功能为合并所有同类的item，并按类型和数量排序 类型小的要再前面，数量多的要在前面，仓库空间要顺序排列中间不能间断
	// 1. 合并：遍历item，把相同的itemtype的道具合并
	itemTypeMap := make(map[int32]int)
	for i := 0; i < len(as.Items); i++ {
		itemType := as.Items[i].ItemType
		if v, ok := itemTypeMap[itemType]; ok {
			if as.Items[v].Number > MaxGridStack {
				return UnexpectedItermNumber
			}
			c := MaxGridStack - as.Items[v].Number
			if as.Items[i].Number < c {
				c = as.Items[i].Number
			}
			as.add(v, itemType, c)
			as.add(i, itemType, -c)
			if as.Items[i].Number == 0 {
				continue
			}
		}
		itemTypeMap[as.Items[i].ItemType] = i
	}

	// 2. 排序: 按照类型/数量的优先级排序
	sort.Slice(as.Items, func(i, j int) bool {
		if as.Items[i].ItemType == EmptyGridItemType {
			return false
		}
		if as.Items[j].ItemType == EmptyGridItemType {
			return true
		}
		if as.Items[i].ItemType == as.Items[j].ItemType {
			return as.Items[i].Number > as.Items[j].Number
		}
		return as.Items[i].ItemType < as.Items[j].ItemType
	})
	return nil
}

func (as *Alliance) add(index int, itemtype int32, num int32) error {
	// 道具增减操作
	// TODO： 是否考虑index之前的空格
	if index >= len(as.Items) {
		return AllianceIsFull
	}

	// 判断下当前格子的状态，是否为空
	d := &(as.Items[index])
	if d.ItemType == EmptyGridItemType {
		d.ItemType = itemtype
	}

	// 当前格子不符合顺延，顺延到下一个格子查找
	if d.ItemType != itemtype {
		return as.add(index+1, itemtype, num)
	}

	// 不允许出现负数
	c := d.Number + num
	if c < 0 {
		return ItemCountErr
	}

	// 当前格子满了，顺延到下一个格子查找
	if c > MaxGridStack {
		if err := as.add(index+1, itemtype, c-MaxGridStack); err != nil {
			return err
		}
		d.Number = MaxGridStack
		// d.Dirty = 1
		return nil
	}

	// 放入当前格子
	d.Number = c
	// d.Dirty = 1
	if d.Number == 0 {
		d.ItemType = EmptyGridItemType
	}

	return nil
}

func (as *Alliance) destroyItem(index int) error {
	return as.add(index, as.Items[index].ItemType, -as.Items[index].Number)
}

func cachekey(allianceid string) string {
	return fmt.Sprintf("alliance:storage:%s", allianceid)
}

// 读取仓库物品
func (as *Alliance) getAllianceStorage(allianceid string) error {
	err := as.getAllianceStorageFromCache(allianceid)
	if err == nil && !as.Dismiss {
		return nil
	}

	if as.Dismiss {
		return AllianceDelete
	}

	// err = as.getAllianceStorageFromDB(allianceid)
	// if err != nil {
	// 	return err
	// }

	// return as.setAllianceStorageToCache()
	return RedisNetworkErr
}

func (as *Alliance) getAllianceStorageFromCache(allianceid string) error {
	p := db.GetRedisPool()
	conn := p.Get()

	b, err := redis.Bytes(conn.Do("GET", cachekey(allianceid)))
	// if err != nil && err != redis.ErrNil {
	if err != nil {
		return err
	}

	return json.Unmarshal(b, as)
}

func (as *Alliance) setAllianceStorageToCache() error {
	p := db.GetRedisPool()
	conn := p.Get()

	b, _ := json.Marshal(as)
	_, err := conn.Do("set", cachekey(as.AllianceId), b)
	return err
}

func (as *Alliance) delAllianceStorageToCache() error {
	p := db.GetRedisPool()
	conn := p.Get()

	_, err := conn.Do("del", cachekey(as.AllianceId))
	return err
}

func (as *Alliance) setAllianceStorageToDB() error {
	// 存入数据库

	// 删除缓存：等待下次读取的时候更新缓存
	// TODO 如果MySQL是读写分离的话，追加个延时删除
	// return as.delAllianceStorageToCache()

	// 直接cache
	if err := as.setAllianceStorageToCache(); err != nil {
		return RedisNetworkErr
	}
	return nil
}

// getAllianceStorageFromDB
func (as *Alliance) getAllianceStorageFromDB(allianceid string) error {
	// 数据库读取
	return errors.New("alliance: no db yet")
}

func lockkey(allianceid string) string {
	return fmt.Sprintf("alliance:storagelock:%s", allianceid)
}

// 加锁
func lockAlliance(allianceid string) error {
	p := db.GetRedisPool()
	conn := p.Get()

	ok, err := redis.String(conn.Do("set", lockkey(allianceid), "1", "nx", "ex", 3))
	if err != nil || ok != "OK" {
		return errors.New("alliance: lock failed")
	}
	return nil
}

// 解锁
func unlockAlliance(allianceid string) error {
	p := db.GetRedisPool()
	conn := p.Get()
	_, err := conn.Do("del", lockkey(allianceid))
	return err
}
