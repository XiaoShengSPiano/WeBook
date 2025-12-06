package sms

import "context"

type Service interface {
	Send(ctx context.Context, tplID string, args []string, numbers ...string) error
	//  Send1(ctx context.Context, tplID string, args []NamedArgs, numbers ...string) error
}

//type NamedArgs struct {
//	Val  string
//	Name string
//}
