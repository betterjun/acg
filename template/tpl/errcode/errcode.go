package errcode

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type APIError interface {
	Code() int
	Message() string
	error
}

type apiErrorImp struct {
	ErrCode    int    `json:"code"`
	ErrMessage string `json:"message"`
}

func (aei *apiErrorImp) Code() int {
	return aei.ErrCode
}

func (aei *apiErrorImp) Message() string {
	return aei.ErrMessage
}

func (aei *apiErrorImp) Error() string {
	return aei.ErrMessage
}

/*
返回的code说明
code 0:正确，其他:错误
*/

// 正常
var Success = &apiErrorImp{0, "ok"}

// token相关
var NoToken = &apiErrorImp{1000, "无令牌请求头"}
var InvalidToken = &apiErrorImp{1001, "无效的登陆会话"}
var TokenExpired = &apiErrorImp{1002, "登陆会话过期"}

// 系统错误
var InvalidParams = &apiErrorImp{1003, "参数无效"}
var RedisError = &apiErrorImp{1004, "nosql数据库错误"}
var MysqlFailed = &apiErrorImp{1005, "sql数据库错误"}
var InternalError = &apiErrorImp{1006, "内部错误"}

var UserAlreadyRegistered = &apiErrorImp{1007, "用户已注册"}
var UserNotExists = &apiErrorImp{1008, "用户不存在"}
var UsernameOrPasswordError = &apiErrorImp{1009, "用户名或密码错误"}
var RecordNotFound = &apiErrorImp{1010, "记录未找到"}

var ExaminationExpired = &apiErrorImp{1011, "考试已结束"}

var ExamQuestionWrong = &apiErrorImp{1012, "不存在的试题类型"}
var AnswerFormatWrong = &apiErrorImp{1013, "答案格式错误"}

var ParamWrong = &apiErrorImp{1014, "参数错误"}
var OldPasswordError = &apiErrorImp{1015, "旧密码错误"}

var SignError = &apiErrorImp{1016, "sign错误"}

var InvalidAccess = &apiErrorImp{1017, "非法访问"}
var ExaminationNotStarted = &apiErrorImp{1018, "考试未开始"}

func Resp(ae APIError, data ...interface{}) gin.H {
	switch len(data) {
	case 0:
		return gin.H{
			"code": ae.Code(),
			"msg":  ae.Message(),
		}
	default:
		if isNil(data[0]) {
			return gin.H{
				"code": ae.Code(),
				"msg":  ae.Message(),
			}
		} else {
			if ae.Code() == Success.Code() {
				return gin.H{
					"code": ae.Code(),
					"msg":  ae.Message(),
					"data": data[0],
				}
			} else {
				return gin.H{
					"code": ae.Code(),
					"msg":  data[0],
				}
			}
		}
	}
}

// golang的接口比较比较坑，需要用反射才能比较
func isNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
