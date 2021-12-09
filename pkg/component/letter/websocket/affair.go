package websocket

// import (
// 	"context"
// )

// type Affair struct {
// }

// // NewAffair 创建
// func NewAffair(ctx context.Context) (*Affair, error) {

// 	return &Affair{}, nil
// }

// // RegisterReq req
// type RegisterReq struct {
// 	UserID string
// 	UUID   string
// 	IP     string
// }

// // RegisterResp resp
// type RegisterResp struct{}

// func (a *Affair) Register(ctx context.Context, req *RegisterReq) (*RegisterResp, error) {
// 	return nil, nil
// }

// // ExpireReq ExpireReq
// type ExpireReq struct {
// 	UserID string
// }

// // ExpireResp ExpireResp
// type ExpireResp struct {
// }

// func (a *Affair) Expire(ctx context.Context, req *ExpireReq) (*ExpireResp, error) {
// 	return nil, nil
// }

// // RenewalReq req
// type RenewalReq struct {
// 	UserID string
// 	UUID   string
// 	IP     string
// }

// // RenewalResp resp
// type RenewalResp struct {
// }

// // Renewal renewal
// func (a *Affair) Renewal(ctx context.Context, req *RenewalReq) (*RenewalResp, error) {
// 	return nil, nil
// }

// // SendReq req
// type SendReq struct {
// 	UserID string

// 	// Specific specify the recipient, if empty, send all
// 	Specific []string

// 	Type    string
// 	Content []byte
// }

// // SendParam param
// type SendParam struct {
// 	UserID  string `json:"userID"`
// 	UUID    string `json:"uuid"`
// 	Type    string `json:"type"`
// 	Content []byte `json:"content"`
// }

// // SendResp resp
// type SendResp struct {
// }

// // Send 发送
// func (a *Affair) Send(ctx context.Context, req *SendReq) (*SendResp, error) {
// 	return nil,nil
// }
