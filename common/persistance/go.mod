module www.seawise.com/shrimps/common/persistance

go 1.16

require (
	github.com/go-redis/redis/v8 v8.11.3
	www.seawise.com/shrimps/common/log v0.0.0-unpublihed
)
	
replace (
	www.seawise.com/shrimps/common/log => ../log
)
