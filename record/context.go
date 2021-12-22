package record

import "context"

type RecordContext interface {
	// closed when work done
	Done() <-chan struct{}

	// explains the reason that the context is Done()
	Err() error

	// get stream URL from this context
	GetStreamUrl() string

	// get streamer's screen ID from this context
	GetStreamer() string
}

type recordContextImpl struct {
	ctx context.Context
}

type contextKey string

const (
	streamerKey  = contextKey("streamer")
	streamUrlKey = contextKey("streamUrl")
)

func newRecordContext(streamer, streamUrl string) (RecordContext, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, streamUrlKey, streamUrl)
	ctx = context.WithValue(ctx, streamerKey, streamer)
	return &recordContextImpl{ctx}, cancelFunc
}

func (ctxImpl *recordContextImpl) Done() <-chan struct{} {
	return ctxImpl.ctx.Done()
}

func (ctxImpl *recordContextImpl) Err() error {
	return ctxImpl.ctx.Err()
}

func (ctxImpl *recordContextImpl) GetStreamUrl() string {
	return ctxImpl.ctx.Value(streamUrlKey).(string)
}

func (ctxImpl *recordContextImpl) GetStreamer() string {
	return ctxImpl.ctx.Value(streamerKey).(string)
}
