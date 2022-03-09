package workspace

import (
	"errors"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/qshell/v2/iqshell/common/config"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
)

func defaultConfig() *config.Config {
	return &config.Config{
		Credentials: &auth.Credentials{
			AccessKey: "",
			SecretKey: nil,
		},
		UseHttps: data.NewBool(true),
		Hosts: &config.Hosts{
			//Rs:  []string{"rs.qiniu.com"},
			//Rsf: []string{"rsf.qiniu.com"},
			//Api: []string{"api.qiniu.com"},
			UC: []string{"uc.qbox.me"},
		},
		Up: &config.Up{
			PutThreshold: data.NewInt64(1024 * 1024 * 4),
			Policy: &storage.PutPolicy{
				Scope:               "",
				Expires:             0,
				IsPrefixalScope:     0,
				InsertOnly:          0,
				DetectMime:          0,
				FsizeMin:            0,
				FsizeLimit:          0,
				MimeLimit:           "",
				ForceSaveKey:        false,
				SaveKey:             "",
				CallbackFetchKey:    0,
				CallbackURL:         "",
				CallbackHost:        "",
				CallbackBody:        "",
				CallbackBodyType:    "",
				ReturnURL:           "",
				ReturnBody:          "",
				PersistentOps:       "",
				PersistentNotifyURL: "",
				PersistentPipeline:  "",
				EndUser:             "",
				DeleteAfterDays:     0,
				FileType:            0,
			},
			LogSetting: &config.LogSetting{
				LogLevel:  data.NewString(config.InfoKey),
				LogFile:   nil,
				LogRotate: data.NewInt(7),
				LogStdout: data.NewBool(true),
			},
			Tasks: &config.Tasks{
				ConcurrentCount:       data.NewInt(3),
				StopWhenOneTaskFailed: data.NewBool(false),
			},
			Retry: &config.Retry{
				Max:      data.NewInt(1),
				Interval: data.NewInt(1000),
			},
		},
		Download: &config.Download{
			LogSetting: &config.LogSetting{
				LogLevel:  data.NewString(config.InfoKey),
				LogFile:   nil,
				LogRotate: data.NewInt(7),
				LogStdout: data.NewBool(true),
			},
			Tasks: &config.Tasks{
				ConcurrentCount:       data.NewInt(3),
				StopWhenOneTaskFailed: data.NewBool(false),
			},
			Retry: &config.Retry{
				Max:      data.NewInt(1),
				Interval: data.NewInt(1000),
			},
		},
	}
}

func checkConfig(cfg *config.Config) (err error) {
	// host
	configHostCount := 0
	if len(cfg.Hosts.Api) > 0 {
		configHostCount += 1
	}
	if len(cfg.Hosts.Rs) > 0 {
		configHostCount += 1
	}
	if len(cfg.Hosts.Rsf) > 0 {
		configHostCount += 1
	}
	if len(cfg.Hosts.Io) > 0 {
		configHostCount += 1
	}
	if len(cfg.Hosts.Up) > 0 {
		configHostCount += 1
	}
	if configHostCount != 0 && configHostCount != 5 {
		err = errors.New("hosts: api/rs/rsf/io/up should config all")
	}
	return
}
