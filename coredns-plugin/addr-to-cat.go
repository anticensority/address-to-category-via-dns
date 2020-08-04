// Package example is a CoreDNS plugin that prints "example" to stdout on every packet received.
//
// It serves as an example CoreDNS plugin with numerous code comments.
package addrToCat

import (
	"context"
	"fmt"
	"io"
	"os"
	"net"

	"github.com/coredns/coredns/plugin"
	//"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/miekg/dns"
	"github.com/coredns/coredns/request"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("addr2cat")

// Example is an example plugin to show how to write a plugin.
type AddrToCat struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (addrToCat AddrToCat) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// This function could be simpler. I.e. just fmt.Println("example") here, but we want to show
	// a slightly more complex example as to make this more interesting.
	// Here we wrap the dns.ResponseWriter in a new ResponseWriter and call the next plugin, when the
	// answer comes back, it will print "example".

	// Debug log that we've have seen the query. This will only be shown when the debug plugin is loaded.
	log.Debug("Received response")
	state := request.Request{W: w, Req: r}
	fmt.Fprintln(out, "NAME:" + state.Name())
	fmt.Fprintln(out, "QNAME:" + state.QName())

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true

	ip := "127.0.0.1"
	var rr dns.RR

	rr = new(dns.A)
	rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
	rr.(*dns.A).A = net.ParseIP(ip).To4()

	a.Extra = []dns.RR{rr}

	w.WriteMsg(a)

	return 0, nil
	// Wrap.
	//pw := NewResponsePrinter(w)

	// Export metric with the server label set to the current server handling the request.
	//requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	// Call next plugin (if any).
	//return plugin.NextOrFailure(addrToCat.Name(), addrToCat.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (addrToCat AddrToCat) Name() string { return "addr2cat" }

// ResponsePrinter wrap a dns.ResponseWriter and will write example to standard output when WriteMsg is called.
type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "example" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	fmt.Fprintln(out, "addr2cat")
	return r.ResponseWriter.WriteMsg(res)
}

// Make out a reference to os.Stdout so we can easily overwrite it for testing.
var out io.Writer = os.Stdout
