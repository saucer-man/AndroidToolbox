package app

// func init() {
// 	path := "./log/androidtoolbox.log"
// 	/* 日志轮转相关函数
// 	`WithLinkName` 为最新的日志建立软连接
// 	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
// 	WithMaxAge 和 WithRotationCount二者只能设置一个
// 	 `WithMaxAge` 设置文件清理前的最长保存时间
// 	 `WithRotationCount` 设置文件清理前最多保存的个数
// 	*/
// 	// 下面配置日志每隔24 小时轮转一个新文件，保留最近 180 小时的日志文件，多余的自动清理掉。
// 	writer, _ := rotatelogs.New(
// 		path+".%Y%m%d%H%M",
// 		rotatelogs.WithLinkName(path),
// 		rotatelogs.WithMaxAge(time.Duration(180)*time.Hour),
// 		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
// 	)
// 	log.SetOutput(writer)
// 	//log.SetFormatter(&log.JSONFormatter{})
// }
