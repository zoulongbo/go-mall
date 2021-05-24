package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//一致性hash算法 虚拟的0-(2^32-1)圆环上  uint32最大 2^32  节点排序二分查找 必需实现len、less、Swap
type uints []uint32

//切片长度
func (x uints) Len() int {
	return len(x)
}

//比较
func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

//切片元素交换
func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}


type Consistent struct {
	//key 为hash值， key的内容为节点的信息
	circle map[uint32]string
	//已经排序的节点hash切片
	sortedHashes uints
	//虚拟节点 用来平衡 数据分布
	VirtualNode int
	//读写锁 map在数据量大的时候肯能不准
	sync.RWMutex
}

func NewConsistent() *Consistent  {
	return &Consistent{
		//初始化 虚拟圆
		circle:       make(map[uint32]string),
		//节点个数
		VirtualNode:  20,
	}
}

//生成key
func (c *Consistent) generateKey(element string, index int) string  {
	return element + strconv.Itoa(index)
}

//生成hash key
func (c *Consistent) hashKey(key string) uint32  {
	if len(key) < 64 {
		var strCatch [64]byte
		copy(strCatch[:], key)
		//使用IEEE 多项式 返回数据的crc-32校验
		return crc32.ChecksumIEEE(strCatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}
//对外 添加节点 加锁
func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}


//删除节点 加锁
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

//向虚拟环中添加 节点
func (c *Consistent) add(element string) {
	//添加虚拟节点 设置副本
	for i:=0;i < c.VirtualNode ; i++  {
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	//更新排序
	c.updateSortedHashes()
}


//删除节点
func (c *Consistent) remove(element string) {
	//添加虚拟节点 设置副本
	for i:=0;i < c.VirtualNode ; i++  {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	//更新排序
	c.updateSortedHashes()
}

//更新排序
func (c *Consistent) updateSortedHashes()  {
	hashes := c.sortedHashes[:0]
	//判断切片容量 是否过大 过大则重制
	if cap(c.sortedHashes) / (c.VirtualNode * 4) > len(c.circle) {
		hashes = nil
	}

	//重新排序
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	//方便二分查找
	sort.Sort(hashes)
	c.sortedHashes = hashes
}
//获取节点信息
func (c *Consistent) Get(name string) (string, error)  {
	c.RLock()
	defer c.RUnlock()
	if len(c.circle ) <= 0 {
		return "", errors.New("节点数据不存在")
	}
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
//二分查找
func (c *Consistent) search(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//二分查找 搜索满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	//如果超出 就置0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}