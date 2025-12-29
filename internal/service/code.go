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

var (
	ErrCodeSendTooFrequently  = repository.ErrCodeSendTooFrequently
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(c *gin.Context, biz string, phone string, inputCode string) error
}

type codeService struct {
	r      repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(r repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{r: r, smsSvc: smsSvc}
}

// 发送验证码 biz:业务类型 code:验证码
func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.r.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

// 校验验证码
func (svc *codeService) Verify(c *gin.Context, biz string, phone string,
	inputCode string) error {
	return svc.r.Verify(c, biz, phone, inputCode)
}

// 生成验证码
func (svc *codeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
