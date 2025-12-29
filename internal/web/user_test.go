package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/domain"
	"webook/internal/service"
	svcmocks "webook/internal/service/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

//func TestMock(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	usersvc := svcmocks.NewMockUserService(ctrl)
//	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).
//		Return(errors.New("mock error"))
//
//	err := usersvc.SignUp(context.Background(), domain.User{
//		Email: "124@qq.com",
//	})
//
//	t.Log(err)
//}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		mock     func(controller *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(controller *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(controller)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "12345678*a",
				}).Return(nil)

				return usersvc
			},
			reqBody: `{
						"email": "123@qq.com",
						"password": "12345678*a",
						"confirmPassword": "12345678*a"
					   }`,
			wantCode: http.StatusOK,
			wantBody: "注册成功......",
		},
		{
			name: "参数不对，bind失败",
			mock: func(controller *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(controller)

				return usersvc
			},
			reqBody: `{
						"email": "123@qq.com",
						"password": "12345678*a",
					   }`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(controller *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(controller)

				return usersvc
			},
			reqBody: `{
						"email": "123@q",
						"password": "12345678*a",
						"confirmPassword": "12345678*a"
					   }`,
			wantCode: http.StatusOK,
			wantBody: "邮箱格式不对......",
		},
		{
			name: "两次输入密码不匹配",
			mock: func(controller *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(controller)

				return usersvc
			},
			reqBody: `{
						"email": "123@qq.com",
						"password": "12345678*a",
						"confirmPassword": "12345678*a0"
					   }`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致......",
		},
		{
			name: "密码格式不对",
			mock: func(controller *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(controller)

				return usersvc
			},
			reqBody: `{
						"email": "123@qq.com",
						"password": "123",
						"confirmPassword": "123"
					   }`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含特殊数字与字符......",
		},
		{
			name: "邮箱冲突",
			mock: func(controller *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(controller)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "12345678*a",
				}).Return(service.ErrUserDuplicateEmail)

				return usersvc
			},
			reqBody: `{
						"email": "123@qq.com",
						"password": "12345678*a",
						"confirmPassword": "12345678*a"
					   }`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突......",
		},
		{
			name: "系统异常",
			mock: func(controller *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(controller)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "12345678*a",
				}).Return(errors.New("系统异常"))

				return usersvc
			},
			reqBody: `{
						"email": "123@qq.com",
						"password": "12345678*a",
						"confirmPassword": "12345678*a"
					   }`,
			wantCode: http.StatusOK,
			wantBody: "系统异常......",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup",
				bytes.NewBuffer([]byte(tc.reqBody)))
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())

		})
	}
}
