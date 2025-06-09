package accrual

import (
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"time"
)

const (
	requestErrorTimeout      = 10 * time.Second
	orderIsProcessingTimeout = 2 * time.Second
)

type Worker struct {
	accrual      usecase.Accrual
	rateLimiter  int
	requestChan  chan string
	responseChan chan model.Accrual
}

func NewWorker(
	accrual usecase.Accrual,
	rateLimiter int,
	requestChan chan string,
	responseChan chan model.Accrual,
) Worker {
	return Worker{
		accrual:      accrual,
		rateLimiter:  rateLimiter,
		requestChan:  requestChan,
		responseChan: responseChan,
	}
}

func (w *Worker) Run() {
	for i := 0; i < w.rateLimiter; i++ {
		go func() {
			for req := range w.requestChan {
				resp, err := w.accrual.Status(req)

				if err != nil {
					go func(orderNumber string) {
						time.Sleep(requestErrorTimeout)
						w.requestChan <- orderNumber
					}(req)
					continue
				}

				if resp.Status == model.AccrualStatusRegistered || resp.Status == model.AccrualStatusProcessing {
					go func(orderNumber string) {
						time.Sleep(orderIsProcessingTimeout)
						w.requestChan <- orderNumber
					}(req)
					continue
				}

				w.responseChan <- resp
			}
		}()
	}
}
