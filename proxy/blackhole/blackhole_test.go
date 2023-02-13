package blackhole_test

import (
	"context"
	"testing"

	"github.com/dharak36/xray-core/common"
	"github.com/dharak36/xray-core/common/buf"
	"github.com/dharak36/xray-core/common/serial"
	"github.com/dharak36/xray-core/proxy/blackhole"
	"github.com/dharak36/xray-core/transport"
	"github.com/dharak36/xray-core/transport/pipe"
)

func TestBlackholeHTTPResponse(t *testing.T) {
	handler, err := blackhole.New(context.Background(), &blackhole.Config{
		Response: serial.ToTypedMessage(&blackhole.HTTPResponse{}),
	})
	common.Must(err)

	reader, writer := pipe.New(pipe.WithoutSizeLimit())

	var mb buf.MultiBuffer
	var rerr error
	go func() {
		b, e := reader.ReadMultiBuffer()
		mb = b
		rerr = e
	}()

	link := transport.Link{
		Reader: reader,
		Writer: writer,
	}
	common.Must(handler.Process(context.Background(), &link, nil))
	common.Must(rerr)
	if mb.IsEmpty() {
		t.Error("expect http response, but nothing")
	}
}
