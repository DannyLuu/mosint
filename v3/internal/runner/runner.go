/*
Copyright © 2023 github.com/DannyLuu
*/
package runner

import (
    "strings"
	"github.com/DannyLuu/mosint/v3/pkg/dns"
	"github.com/DannyLuu/mosint/v3/pkg/scrape/googlesearch"
	"github.com/DannyLuu/mosint/v3/pkg/services/breachdirectory"
	"github.com/DannyLuu/mosint/v3/pkg/services/emailrep"
	"github.com/DannyLuu/mosint/v3/pkg/services/haveibeenpwned"
	"github.com/DannyLuu/mosint/v3/pkg/services/hunter"
	"github.com/DannyLuu/mosint/v3/pkg/services/intelx"
	"github.com/DannyLuu/mosint/v3/pkg/services/ipapi"
	"github.com/DannyLuu/mosint/v3/pkg/services/psbdmp"
	"github.com/DannyLuu/mosint/v3/pkg/social/instagram"
	"github.com/DannyLuu/mosint/v3/pkg/social/spotify"
	"github.com/DannyLuu/mosint/v3/pkg/social/twitter"
	"github.com/gammazero/workerpool"
)

type Runner struct {
	Email            string
	DnsC             *dns.Dns
	GoogleSearchC    *googlesearch.GoogleSearch
	BreachDirectoryC *breachdirectory.BreachDirectory
	HaveibeenpwnedC  *haveibeenpwned.HaveIBeenPwned
	EmailRepC        *emailrep.Emailrep
	HunterC          *hunter.Hunter
	IntelxC          *intelx.Intelx
	IpApiC           *ipapi.Ipapi
	PsbdmpC          *psbdmp.Psbdmp
	InstagramC       *instagram.Instagram
	SpotifyC         *spotify.Spotify
	TwitterC         *twitter.Twitter
}

func New(email string) *Runner {
	return &Runner{
		Email:            email,
		DnsC:             dns.New(),
		GoogleSearchC:    googlesearch.New(),
		BreachDirectoryC: breachdirectory.New(),
		HaveibeenpwnedC:  haveibeenpwned.New(),
		EmailRepC:        emailrep.New(),
		HunterC:          hunter.New(),
		IntelxC:          intelx.New(),
		IpApiC:           ipapi.New(),
		PsbdmpC:          psbdmp.New(),
		InstagramC:       instagram.New(),
		SpotifyC:         spotify.New(),
		TwitterC:         twitter.New(),
	}
}

func (c *Runner) Start() {
	email := c.Email

	// Extract the part of the email before the '@'
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		// Handle invalid email case
		println("Invalid email format")
		return
	}
	emailBeforeAt := email[:atIndex]
	
	runners := []func(string){
		c.DnsC.Resolver,
		c.GoogleSearchC.Search,
		func(s string) { c.GoogleSearchC.Search(emailBeforeAt) }, // Search email minus the '@'
		c.BreachDirectoryC.Lookup,
		c.HaveibeenpwnedC.Lookup,
		c.EmailRepC.Lookup,
		c.HunterC.Lookup,
		c.IntelxC.Search,
		c.IpApiC.GetInfo,
		c.PsbdmpC.Search,
		c.InstagramC.Check,
		c.SpotifyC.Check,
		c.TwitterC.Check,
	}

	wp := workerpool.New(12)
	for _, runner := range runners {
		runner := runner
		wp.Submit(func() {
			runner(email)
		})
	}
	wp.StopWait()
}

func (c *Runner) Print() {
	println()
	c.EmailRepC.Print()

	println()
	c.HunterC.Print()

	println()
	c.GoogleSearchC.Print()

	println()
	c.InstagramC.Print()
	c.SpotifyC.Print()
	c.TwitterC.Print()

	println()
	c.PsbdmpC.Print()

	println()
	c.IntelxC.Print()

	println()
	c.BreachDirectoryC.Print()

	println()
	c.HaveibeenpwnedC.Print()

	println()
	c.IpApiC.Print()

	println()
	c.DnsC.Print()
}
