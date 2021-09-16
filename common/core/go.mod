module www.seawise.com/shrimps/common/core

go 1.16

require (
	www.seawise.com/shrimps/common/exposed v0.0.0-unpublished
	www.seawise.com/shrimps/common/persistance v0.0.0-unpublished
	www.seawise.com/shrimps/common/log v0.0.0-unpublished
)

replace (
	www.seawise.com/shrimps/common/exposed => ../exposed/
    www.seawise.com/shrimps/common/persistance => ../persistance/
	www.seawise.com/shrimps/common/log => ../log/
)
