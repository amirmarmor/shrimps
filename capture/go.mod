module www.seawise.com/shrimps/capture

go 1.16

require (
	gocv.io/x/gocv v0.28.0
	www.seawise.com/shrimps/common/core v0.0.0-unpublished
	www.seawise.com/shrimps/common/log v0.0.0-unpublished
	www.seawise.com/shrimps/common/persistance v0.0.0-unpublished
)

replace (
	www.seawise.com/shrimps/common/core => ../common/core
	www.seawise.com/shrimps/common/exposed => ../common/exposed
	www.seawise.com/shrimps/common/log => ../common/log
	www.seawise.com/shrimps/common/persistance => ../common/persistance
)
