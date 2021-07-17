package validate

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/i18n"

	"github.com/getfider/fider/app/pkg/env"
)

var emailRegex = regexp.MustCompile("^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
var hostnameRegex = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)

//Email validates given email address
func Email(ctx context.Context, email string) []string {
	email = strings.ToLower(email)

	if len(email) > 200 {
		return []string{i18n.T(ctx, "validation.maxstringlen",
			i18n.Params{"name": i18n.T(ctx, "property.email")},
			i18n.Params{"len": 200},
		)}
	}

	if !emailRegex.MatchString(email) {
		return []string{i18n.T(ctx, "validation.custom.invalidemail",
			i18n.Params{"email": email},
		)}
	}

	return []string{}
}

//URL validates given URL
func URL(ctx context.Context, rawurl string) []string {
	if len(rawurl) > 300 {
		return []string{i18n.T(ctx, "validation.maxstringlen",
			i18n.Params{"name": "URL"},
			i18n.Params{"len": 300},
		)}
	}

	_, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return []string{i18n.T(ctx, "validation.custom.invalidurl",
			i18n.Params{"url": rawurl},
		)}
	}

	return []string{}
}

//CNAME validates given cname
func CNAME(ctx context.Context, cname string) []string {
	cname = strings.ToLower(cname)

	if !env.IsSingleHostMode() {
		domain := env.MultiTenantDomain()
		if strings.HasSuffix(cname, domain) || cname == domain[1:] {
			return []string{i18n.T(ctx, "validation.custom.invalidcustomdomain",
				i18n.Params{"domain": cname},
			)}
		}
	}

	if len(cname) > 100 {
		return []string{i18n.T(ctx, "validation.maxstringlen",
			i18n.Params{"name": i18n.T(ctx, "property.customdomain")},
			i18n.Params{"len": 100},
		)}
	}

	if !hostnameRegex.MatchString(cname) || !strings.Contains(cname, ".") {
		return []string{i18n.T(ctx, "validation.custom.invalidcustomdomain",
			i18n.Params{"domain": cname},
		)}
	}

	isAvailable := &query.IsCNAMEAvailable{CNAME: cname}
	bus.MustDispatch(ctx, isAvailable)
	if !isAvailable.Result {
		return []string{i18n.T(ctx, "validation.custom.customdomaintaken")}
	}

	return []string{}
}
