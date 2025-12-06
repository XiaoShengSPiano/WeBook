package cache

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

var (
	ErrCodeSendTooFrequently  = errors.New("code send too frequently")
	ErrCodeVerifyTooManyTimes = errors.New("code verify too many")
	ErrUnknowForCode          = errors.New("unknow code")
)

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{client: client}
}

func (c *CodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *CodeCache) SetCode(ctx context.Context, biz string, phone string,
	code string) error {
	key := c.key(biz, phone)
	res, err := c.client.Eval(ctx, luaSetCode, []string{key}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 设置成功
		return nil
	case -1:
		// 发送太频繁
		return ErrCodeSendTooFrequently
	default:
		return ErrUnknowForCode
	}
}

func (c *CodeCache) VerifyCode(ctx context.Context, biz string, phone string, inputCode string) error {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case -1:
		// 验证次数太多, 如果频繁出现这个错误，需要进行警告
		return ErrCodeVerifyTooManyTimes
	case -2:
		// 验证码不存在
		return ErrUnknowForCode

	}
	return ErrUnknowForCode
}
