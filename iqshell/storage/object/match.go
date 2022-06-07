package object

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/flow"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"os"
)

const (
	MatchCheckModeFileHash = 0
	MatchCheckModeFileSize = 1
)

type MatchApiInfo struct {
	Bucket    string // 文件所在七牛云的空间名，【必选】
	Key       string // 文件所在七牛云的 Key， 【必选】
	LocalFile string // 本地文件路径；【必选】
	CheckMode int    // 检测模式， 0: 检测 hash，其他检测 size 【可选】
	FileHash  string // 文件 Etag，可以是 etagV1, 也可以是 etagV2；【可选】没有会从服务获取
	FileSize  int64  // 文件大小；【可选】没有会从服务获取
}

func (m *MatchApiInfo) WorkId() string {
	return utils.Md5Hex(fmt.Sprintf("%s:%s:%s:%s", m.Bucket, m.Key, m.LocalFile, m.FileHash))
}

func (m *MatchApiInfo) CheckModeHash() bool {
	return m.CheckMode == MatchCheckModeFileHash
}

type MatchResult struct {
	Exist bool
	Match bool
}

var _ flow.Result = (*MatchResult)(nil)

func (m *MatchResult) IsValid() bool {
	return m != nil
}

func Match(info MatchApiInfo) (match *MatchResult, err *data.CodeError) {
	if len(info.LocalFile) == 0 {
		return &MatchResult{
			Exist: false,
			Match: false,
		}, data.NewEmptyError().AppendDesc("Match Check, file is empty")
	}

	if info.CheckModeHash() {
		return matchHash(info)
	} else {
		return matchSize(info)
	}
}

func matchSize(info MatchApiInfo) (match *MatchResult, err *data.CodeError) {
	match = &MatchResult{
		Exist: false,
		Match: false,
	}

	if info.FileSize <= 0 {
		if stat, sErr := Status(StatusApiInfo{
			Bucket:   info.Bucket,
			Key:      info.Key,
			NeedPart: false,
		}); sErr != nil {
			return match, data.NewEmptyError().AppendDescF("Match check size, get file status error:%v", sErr)
		} else {
			info.FileSize = stat.FSize
			match.Exist = true
		}
	}

	stat, sErr := os.Stat(info.LocalFile)
	if sErr != nil {
		return match, data.NewEmptyError().AppendDescF("Match check size, get local file stat error:%v", sErr)
	}
	if info.FileSize == stat.Size() {
		match.Match = true
		return match, nil
	} else {
		match.Match = false
		return match, data.NewEmptyError().AppendDescF("Match check size, size don't match, file:%s except:%d but:%d")
	}
}

func matchHash(info MatchApiInfo) (result *MatchResult, err *data.CodeError) {
	result = &MatchResult{
		Exist: false,
		Match: false,
	}

	var serverObjectStat *StatusResult
	if len(info.FileHash) == 0 {
		if stat, sErr := Status(StatusApiInfo{
			Bucket:   info.Bucket,
			Key:      info.Key,
			NeedPart: true,
		}); sErr != nil {
			return result, data.NewEmptyError().AppendDescF("Match Check, get file status error:%v", sErr)
		} else {
			info.FileHash = stat.Hash
			serverObjectStat = &stat
			result.Exist = true
		}
	}

	hashFile, oErr := os.Open(info.LocalFile)
	if oErr != nil {
		return result, data.NewEmptyError().AppendDescF("Match check hash, get local file error:%v", oErr)
	}

	// 计算本地文件 hash
	var hash string
	if utils.IsSignByEtagV2(info.FileHash) {
		log.DebugF("Match check hash: get etag by v2 for key:%s", info.Key)
		if serverObjectStat == nil {
			if stat, sErr := Status(StatusApiInfo{
				Bucket:   info.Bucket,
				Key:      info.Key,
				NeedPart: true,
			}); sErr != nil {
				return result, data.NewEmptyError().AppendDescF("Match check hash, etag v2, get file status error:%v", sErr)
			} else {
				serverObjectStat = &stat
			}
		}
		if h, eErr := utils.EtagV2(hashFile, serverObjectStat.Parts); eErr != nil {
			return result, data.NewEmptyError().AppendDescF("Match check hash, get file etag v2 error:%v", eErr)
		} else {
			hash = h
		}
		log.DebugF("Match check hash, get etag by v2 for key:%s hash:%s", info.Key, hash)
	} else {
		log.DebugF("Match check hash, get etag by v1 for key:%s", info.Key)
		if h, eErr := utils.EtagV1(hashFile); eErr != nil {
			return result, data.NewEmptyError().AppendDescF("Match check hash, get file etag v1 error:%v", eErr)
		} else {
			hash = h
		}
		log.DebugF("Match check hash, get etag by v1 for key:%s hash:%s", info.Key, hash)
	}
	log.DebugF("Match check hash,       server hash, key:%s hash:%s", info.Key, hash)
	if hash != info.FileHash {
		return result, data.NewEmptyError().AppendDescF("Match check hash, file hash doesn't match for key:%s, local file hash:%s server file hash:%s", info.Key, hash, info.FileHash)
	}

	result.Match = true
	return result, nil
}
