package async

//
//import (
//	"Neo/Workplace/goland/src/GeekGo/webook/internal/repository"
//	"Neo/Workplace/goland/src/GeekGo/webook/internal/service/sms"
//	"context"
//	"errors"
//	"sync"
//	"time"
//)
//
//type SMSService struct {
//	svc      sms.Service
//	repo     repository.SMSAsyncRepository
//	retryMax int           // 最大重试次数
//	retryGap time.Duration // 重试间隔
//	mu       sync.Mutex    // 保护异步 goroutine 启动的互斥锁
//	started  bool          // 标记是否已经启动异步 goroutine
//}
//
//// StartAsync 启动异步重试任务（保证只启动一次）
//func (s *SMSService) StartAsync() {
//	s.mu.Lock()
//	if s.started {
//		s.mu.Unlock()
//		return
//	}
//	s.started = true
//	s.mu.Unlock()
//
//	go func() {
//		for {
//			time.Sleep(s.retryGap) // 控制重试间隔
//			reqs := s.repo.FindPendingRequests()
//			for _, req := range reqs {
//				err := s.svc.Send(context.Background(), req.Biz, req.Args, req.Numbers...)
//				if err == nil {
//					s.repo.MarkAsSent(req.ID) // 发送成功，标记已发送
//				} else {
//					s.repo.IncrementRetryCount(req.ID)
//					if req.RetryCount+1 >= s.retryMax {
//						s.repo.MarkAsFailed(req.ID) // 达到最大重试次数，标记失败
//					}
//				}
//			}
//		}
//	}()
//}
//
//// Send 发送短信，遇到限流或服务商崩溃时存入数据库，等待异步重试
//func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
//	err := s.svc.Send(ctx, biz, args, numbers...)
//	if err != nil {
//		if s.isRateLimited(err) || s.isServiceDown(err) {
//			s.repo.StoreAsyncRequest(biz, args, numbers...)
//			return errors.New("请求已转为异步处理")
//		}
//	}
//	return err
//}
//
//// isRateLimited 判断是否触发限流
//func (s *SMSService) isRateLimited(err error) bool {
//	return errors.Is(err, sms.ErrRateLimited)
//}
//
//// isServiceDown 判断是否服务商崩溃（自定义机制）
//func (s *SMSService) isServiceDown(err error) bool {
//	// 机制：如果连续 N 次失败，并且失败类型是网络错误/超时，判定为崩溃
//	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
//}
//
//// NewSMSService 创建新的异步容错短信服务
//func NewSMSService(svc sms.Service, repo repository.SMSAsyncRepository, retryMax int, retryGap time.Duration) sms.Service {
//	service := &SMSService{
//		svc:      svc,
//		repo:     repo,
//		retryMax: retryMax,
//		retryGap: retryGap,
//	}
//	service.StartAsync() // 确保异步任务启动
//	return service
//}
