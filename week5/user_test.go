package web

import (
	"awesomeProject/webook/internal/domain"
	"awesomeProject/webook/internal/service"
	svcmocks "awesomeProject/webook/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Huawei@123",
				}).Return(nil)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "Huawei@123",
	"confirmPassword": "Huawei@123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对,bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "Huawei@123",
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "12312",
				//	Password: "Huawei@123",
				//}).Return(nil)
				return usersvc
			},
			reqBody: `
{
	"email": "123@q",
	"password": "Huawei@123",
	"confirmPassword": "Huawei@123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "您的邮箱格式不对",
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "123@qq.com",
				//	Password: "Huawei@123",
				//}).Return(nil)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "Huawei@123",
	"confirmPassword": "Huawei@1213"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				//usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
				//	Email:    "123@qq.com",
				//	Password: "Huawei@123",
				//}).Return(nil)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "Huawei",
	"confirmPassword": "Huawei"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位,包括数字、特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Huawei@123",
				}).Return(service.ErrUserDuplicateEmail)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "Huawei@123",
	"confirmPassword": "Huawei@123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "Huawei@123",
				}).Return(errors.New("随便一个error"))
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "Huawei@123",
	"confirmPassword": "Huawei@123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			// 用不上 codeSvc
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			// 这里你就可以继续使用 req

			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}
}

func TestUserHandler_LoginSMS(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBody  string
		wantCode int
		wantBody Result
	}{
		{
			name: "验证码校验通过",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codesvc := svcmocks.NewMockCodeService(ctrl)
				usersvc := svcmocks.NewMockUserService(ctrl)
				codesvc.EXPECT().Verify(gomock.Any(), biz, "11111111", "11").Return(true, nil)
				usersvc.EXPECT().FindOrCreate(gomock.Any(), "11111111").Return(domain.User{
					Phone: "11111111",
				}, nil)
				return usersvc, codesvc
			},
			reqBody: `
{
	"phone": "11111111",
	"code": "11"
}
`,
			wantCode: http.StatusOK,
			wantBody: Result{
				Code: 5,
				Msg:  "验证码校验通过",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl))
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			// 这里你就可以继续使用 req

			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			server.ServeHTTP(resp, req)
			print(11111)

			var result Result
			err = json.Unmarshal([]byte(resp.Body.String()), &result)
			if err != nil {
				t.Errorf("无法解析 JSON: %v", err)
			}
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, result)

		})
	}
}
