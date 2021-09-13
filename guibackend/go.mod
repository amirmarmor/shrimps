module www.seawise.com/shrimps/guibackend

go 1.16

require (
	github.com/labstack/echo/v4 v4.4.0
	www.seawise.com/shrimps/common/core v0.0.0-unpublished
	www.seawise.com/shrimps/common/exposed v0.0.0-unpublished
    www.seawise.com/shrimps/common/persistance v0.0.0-unpublished
	www.seawise.com/shrimps/common/log v0.0.0-unpublished
)

replace (
	www.seawise.com/shrimps/common/core => ../common/core
    www.seawise.com/shrimps/common/exposed => ../common/exposed
    www.seawise.com/shrimps/common/persistance => ../common/persistance
   	www.seawise.com/shrimps/common/log => ../common/log

)
