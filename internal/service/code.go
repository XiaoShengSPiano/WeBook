package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"

	"github.com/gin-gonic/gin"
)

const codeTplId = "1877556"

type CodeService struct {
	r   *repository.CodeRepository
	sms sms.Service
}

// 发送验证码 biz:业务类型 code:验证码
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.r.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	err = svc.sms.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

// 校验验证码
func (svc *CodeService) Verify(c *gin.Context, biz string, phone string,
	inputCode string) error {
	return svc.r.Verify(c, biz, phone, inputCode)
}

// 生成验证码
func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
