package errs

// 用户模块 01
const (
	// UserInputValid 用户模块输入错误，这是一个含糊的错误
	UserInputValid = 401001

	// UserInvalidOrPassword 用户不存在或者密码错误
	// 理论上也属于用户输入错误的范畴，但是需要准确定义的原因是因为我们需要关注它
	// 它的出现有可能是有人在暴露破解
	UserInvalidOrPassword = 401002

	// UserInternalServerError 用户模块内部错误
	UserInternalServerError = 501001
)

// 文章模块 02
const (
	ArticleInvalidInput        = 402001
	ArticleInternalServerError = 502001
)

type Code struct {
	Number  int
	Message string
}

var UserInputValidV1 = Code{Number: 401001, Message: "用户输入错误"}
