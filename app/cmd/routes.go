package cmd

import (
	"net/http"
	"strings"
	"time"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/web"
)

func routes(r *web.Engine) *web.Engine {
	r.Worker().Use(middlewares.WorkerSetup())

	r.Use(middlewares.CatchPanic())

	r.NotFound(func(c *web.Context) error {
		next := func(c *web.Context) error {
			return c.NotFound()
		}
		next = middlewares.Tenant()(next)
		next = middlewares.WebSetup()(next)
		return next(c)
	})

	r.Use(middlewares.Secure())
	r.Use(middlewares.Compress())

	assets := r.Group()
	{
		assets.Use(middlewares.CORS())
		assets.Use(middlewares.ClientCache(365 * 24 * time.Hour))
		assets.Get("/static/favicon", handlers.Favicon())
		assets.Static("/assets/*filepath", "dist")
	}

	r.Use(middlewares.Session())

	r.Get("/_health", handlers.Health())
	r.Get("/robots.txt", handlers.RobotsTXT())
	r.Post("/_api/log-error", handlers.LogError())

	r.Use(middlewares.Maintenance())
	r.Use(middlewares.WebSetup())
	r.Use(middlewares.Tenant())
	r.Use(middlewares.User())

	r.Get("/privacy", handlers.LegalPage("Privacy Policy", "privacy.md"))
	r.Get("/terms", handlers.LegalPage("Terms of Service", "terms.md"))

	r.Post("/_api/tenants", handlers.CreateTenant())
	r.Get("/_api/tenants/:subdomain/availability", handlers.CheckAvailability())
	r.Get("/signup", handlers.SignUp())
	r.Get("/oauth/:provider", handlers.SignInByOAuth())
	r.Get("/oauth/:provider/callback", handlers.OAuthCallback())

	//Starting from this step, a Tenant is required
	r.Use(middlewares.RequireTenant())

	r.Get("/sitemap.xml", handlers.Sitemap())

	tenantAssets := r.Group()
	{
		tenantAssets.Use(middlewares.ClientCache(5 * 24 * time.Hour))
		tenantAssets.Get("/static/avatars/letter/:id/:name", handlers.LetterAvatar())
		tenantAssets.Get("/static/avatars/gravatar/:id/:name", handlers.Gravatar())

		tenantAssets.Use(middlewares.ClientCache(30 * 24 * time.Hour))
		tenantAssets.Get("/static/favicon/*bkey", handlers.Favicon())
		tenantAssets.Get("/static/images/*bkey", handlers.ViewUploadedImage())
		tenantAssets.Get("/static/custom/:md5.css", func(c *web.Context) error {
			return c.Blob(http.StatusOK, "text/css", []byte(c.Tenant().CustomCSS))
		})
	}

	r.Get("/_design", handlers.Page("Design System", "A preview of Fider UI elements", "DesignSystem.page"))
	r.Get("/signup/verify", handlers.VerifySignUpKey())
	r.Get("/signout", handlers.SignOut())
	r.Get("/oauth/:provider/token", handlers.OAuthToken())
	r.Get("/oauth/:provider/echo", handlers.OAuthEcho())

	//If tenant is pending, block it from using any other route
	r.Use(middlewares.BlockPendingTenants())

	r.Get("/signin", handlers.SignInPage())
	r.Get("/not-invited", handlers.NotInvitedPage())
	r.Get("/signin/verify", handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))
	r.Get("/invite/verify", handlers.VerifySignInKey(enum.EmailVerificationKindUserInvitation))
	r.Post("/_api/signin/complete", handlers.CompleteSignInProfile())
	r.Post("/_api/signin", handlers.SignInByEmail())

	//Block if it's a locked tenant with a non-administrator user
	r.Use(middlewares.BlockLockedTenants())

	//Block if it's private tenant with unauthenticated user
	r.Use(middlewares.CheckTenantPrivacy())

	r.Get("/", handlers.Index())
	r.Get("/posts/:number", handlers.PostDetails())
	r.Get("/posts/:number/:slug", handlers.PostDetails())

	/*
	** This is a temporary redirect and should be removed in the future
	** START
	 */
	r.Get("/ideas/:number", func(c *web.Context) error {
		return c.PermanentRedirect(strings.Replace(c.Request.URL.Path, "/ideas/", "/posts/", 1))
	})
	r.Get("/ideas/:number/*all", func(c *web.Context) error {
		return c.PermanentRedirect(strings.Replace(c.Request.URL.Path, "/ideas/", "/posts/", 1))
	})
	/*
	** END
	 */

	ui := r.Group()
	{
		//From this step, a User is required
		ui.Use(middlewares.IsAuthenticated())

		ui.Get("/settings", handlers.UserSettings())
		ui.Get("/notifications", handlers.Notifications())
		ui.Get("/notifications/:id", handlers.ReadNotification())
		ui.Get("/change-email/verify", handlers.VerifyChangeEmailKey())

		ui.Delete("/_api/user", handlers.DeleteUser())
		ui.Post("/_api/user/regenerate-apikey", handlers.RegenerateAPIKey())
		ui.Post("/_api/user/settings", handlers.UpdateUserSettings())
		ui.Post("/_api/user/change-email", handlers.ChangeUserEmail())
		ui.Post("/_api/notifications/read-all", handlers.ReadAllNotifications())
		ui.Get("/_api/notifications/unread/total", handlers.TotalUnreadNotifications())

		// From this step, only Collaborators and Administrators are allowed
		ui.Use(middlewares.IsAuthorized(enum.RoleCollaborator, enum.RoleAdministrator))

		// locale is forced to English for administrative pages.
		// This is meant to be removed when all pages are translated.
		ui.Use(middlewares.SetLocale("en"))

		ui.Get("/admin", handlers.GeneralSettingsPage())
		ui.Get("/admin/advanced", handlers.AdvancedSettingsPage())
		ui.Get("/admin/privacy", handlers.Page("Privacy · Site Settings", "", "PrivacySettings.page"))
		ui.Get("/admin/invitations", handlers.Page("Invitations · Site Settings", "", "Invitations.page"))
		ui.Get("/admin/members", handlers.ManageMembers())
		ui.Get("/admin/tags", handlers.ManageTags())
		ui.Get("/admin/authentication", handlers.ManageAuthentication())
		ui.Get("/_api/admin/oauth/:provider", handlers.GetOAuthConfig())

		//From this step, only Administrators are allowed
		ui.Use(middlewares.IsAuthorized(enum.RoleAdministrator))

		ui.Get("/admin/export", handlers.Page("Export · Site Settings", "", "Export.page"))
		ui.Get("/admin/export/posts.csv", handlers.ExportPostsToCSV())
		ui.Get("/admin/export/backup.zip", handlers.ExportBackupZip())
		ui.Post("/_api/admin/settings/general", handlers.UpdateSettings())
		ui.Post("/_api/admin/settings/advanced", handlers.UpdateAdvancedSettings())
		ui.Post("/_api/admin/settings/privacy", handlers.UpdatePrivacy())
		ui.Post("/_api/admin/oauth", handlers.SaveOAuthConfig())
		ui.Post("/_api/admin/roles/:role/users", handlers.ChangeUserRole())
		ui.Put("/_api/admin/users/:userID/block", handlers.BlockUser())
		ui.Delete("/_api/admin/users/:userID/block", handlers.UnblockUser())
	}

	api := r.Group()
	{
		api.Get("/api/v1/posts", apiv1.SearchPosts())
		api.Get("/api/v1/tags", apiv1.ListTags())
		api.Get("/api/v1/posts/:number", apiv1.GetPost())
		api.Get("/api/v1/posts/:number/comments", apiv1.ListComments())
		api.Get("/api/v1/posts/:number/comments/:id", apiv1.GetComment())

		//From this step, a User is required
		api.Use(middlewares.IsAuthenticated())

		api.Post("/api/v1/posts", apiv1.CreatePost())
		api.Put("/api/v1/posts/:number", apiv1.UpdatePost())
		api.Post("/api/v1/posts/:number/comments", apiv1.PostComment())
		api.Put("/api/v1/posts/:number/comments/:id", apiv1.UpdateComment())
		api.Delete("/api/v1/posts/:number/comments/:id", apiv1.DeleteComment())
		api.Post("/api/v1/posts/:number/votes", apiv1.AddVote())
		api.Delete("/api/v1/posts/:number/votes", apiv1.RemoveVote())
		api.Post("/api/v1/posts/:number/subscription", apiv1.Subscribe())
		api.Delete("/api/v1/posts/:number/subscription", apiv1.Unsubscribe())

		//From this step, only Collaborators and Administrators are allowed
		api.Use(middlewares.IsAuthorized(enum.RoleCollaborator, enum.RoleAdministrator))

		// locale is forced to English for administrative API.
		// This is meant to be removed when all API are translated.
		ui.Use(middlewares.SetLocale("en"))

		api.Get("/api/v1/users", apiv1.ListUsers())
		api.Get("/api/v1/posts/:number/votes", apiv1.ListVotes())
		api.Post("/api/v1/invitations/send", apiv1.SendInvites())
		api.Post("/api/v1/invitations/sample", apiv1.SendSampleInvite())
		api.Put("/api/v1/posts/:number/status", apiv1.SetResponse())
		api.Post("/api/v1/posts/:number/tags/:slug", apiv1.AssignTag())
		api.Delete("/api/v1/posts/:number/tags/:slug", apiv1.UnassignTag())

		//From this step, only Administrators are allowed
		api.Use(middlewares.IsAuthorized(enum.RoleAdministrator))

		api.Post("/api/v1/users", apiv1.CreateUser())
		api.Delete("/api/v1/posts/:number", apiv1.DeletePost())
		api.Post("/api/v1/tags", apiv1.CreateEditTag())
		api.Put("/api/v1/tags/:slug", apiv1.CreateEditTag())
		api.Delete("/api/v1/tags/:slug", apiv1.DeleteTag())
	}

	return r
}
