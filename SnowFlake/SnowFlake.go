package SnowFlake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

/*
 * 算法解释
 * SnowFlake的结构如下(每部分用-分开):<br>
 * 0 - 0000000000 0000000000 0000000000 0000000000 0 - 00000 - 00000 - 000000000000 <br>
 * 1位标识，由于long基本类型在Java中是带符号的，最高位是符号位，正数是0，负数是1，所以id一般是正数，最高位是0<br>
 * 41位时间截(毫秒级)，注意，41位时间截不是存储当前时间的时间截，而是存储时间截的差值（当前时间截 - 开始时间截)
 * 得到的值），这里的的开始时间截，一般是我们的id生成器开始使用的时间，由我们程序来指定的（如下的epoch属性）。
 * 41位的时间截，可以使用69年，年T = (1L << 41) / (1000L * 60 * 60 * 24 * 365) = 69<br>
 * 10位的数据机器位，可以部署在1024个节点，包括5位datacenterId和5位workerId<br>
 * 12位序列，毫秒内的计数，12位的计数顺序号支持每个节点每毫秒(同一机器，同一时间截)产生4096个ID序号<br>
 * 加起来刚好64位，为一个Long型。<br>
 * SnowFlake的优点是，整体上按照时间自增排序，并且整个分布式系统内不会产生ID碰撞(由数据中心ID和机器ID作区分)，并且效率较高，经测试，SnowFlake每秒能够产生26万ID左右。
 */
const (
	//t := time.Date(2015, 1, 1, 00, 00, 00, 00, time.Local).UnixNano() / 1e6;//获取时间戳 毫秒
	//开始时间戳 2015-1-1
	epoch int64 = 1420041600000
	// 机器id所占的位数
	workerIdBits int64 = 5
	// 数据标识id所占的位数
	datacenterIdBits int64 = 5
	//支持的最大机器id，结果是31 (这个移位算法可以很快的计算出几位二进制数所能表示的最大十进制数)
	maxWorkerId int64 = -1 ^ (-1 << workerIdBits)
	// 支持的最大数据标识id，结果是31
	maxDatacenterId int64 = -1 ^ (-1 << datacenterIdBits)
	//序列在id中占的位数
	sequenceBits int64 = 12
	// 机器ID向左移12位
	workerIdShift int64 = sequenceBits
	// 数据标识id向左移17位(12+5)
	datacenterIdShift int64 = sequenceBits + workerIdBits
	// 时间截向左移22位(5+5+12)
	timestampLeftShift int64 = sequenceBits + workerIdBits + datacenterIdBits
	// 生成序列的掩码，这里为4095 (0b111111111111=0xfff=4095)
	sequenceMask int64 = -1 ^ (-1 << sequenceBits)
)

type IdWorker struct {
	mutex         sync.Mutex // 添加互斥锁 确保并发安全
	lastTimestamp int64      // 上次生成ID的时间截
	workerId      int64      // 工作机器ID(0~31)
	datacenterId  int64      //数据中心ID(0~31)
	sequence      int64      // 毫秒内序列(0~4095)
}

// CreateWorker 创建SnowflakeIdWorker workerId 工作ID (0~31) datacenterId 数据中心ID (0~31)
func CreateWorker(workerId int64, datacenterId int64) (*IdWorker, error) {
	if workerId < 0 || workerId > maxWorkerId {
		return nil, errors.New("worker id excess of quantity")
	}
	if datacenterId < 0 || datacenterId > maxDatacenterId {
		return nil, errors.New("datacenter id excess of quantity")
	}
	// 生成一个新节点
	return &IdWorker{
		lastTimestamp: 0,
		workerId:      workerId,
		datacenterId:  datacenterId,
		sequence:      0,
	}, nil
}

// NextId 获取ID
func (w *IdWorker) NextId() int64 {
	// 保障线程安全 加锁
	w.mutex.Lock()
	// 生成完成后 解锁
	defer w.mutex.Unlock()
	// 获取生成时的时间戳 毫秒
	now := time.Now().UnixNano() / 1e6
	//如果当前时间小于上一次ID生成的时间戳，说明系统时钟回退过这个时候应当抛出异常
	if now < w.lastTimestamp {
		errors.New("Clock moved backwards")
		//根据需要自定义错误码
		return 3001
	}
	if w.lastTimestamp == now {
		w.sequence = (w.sequence + 1) & sequenceMask
		if w.sequence == 0 {
			// 阻塞到下一个毫秒，直到获得新的时间戳
			for now <= w.lastTimestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		w.sequence = 0
	}
	// 将机器上一次生成ID的时间更新为当前时间
	w.lastTimestamp = now
	ID := int64((now-epoch)<<timestampLeftShift | w.datacenterId<<datacenterIdShift | (w.workerId << workerIdShift) | w.sequence)
	return ID
}

// ConvertToBin 将十进制数字转化为二进制字符串
func ConvertToBin(num int64) string {
	s := ""
	if num == 0 {
		return "0"
	}
	// num /= 2 每次循环的时候 都将num除以2  再把结果赋值给 num
	for ; num > 0; num /= 2 {
		lsb := num % 2
		// 将数字强制性转化为字符串
		s = strconv.FormatInt(lsb, 10) + s
	}
	return s
}
