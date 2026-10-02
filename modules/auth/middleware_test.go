package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/farbeyka/go-admin/context"
	"github.com/farbeyka/go-admin/modules/config"
	"github.com/farbeyka/go-admin/modules/constant"
	"github.com/farbeyka/go-admin/plugins/admin/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	config.Initialize(&config.Config{
		UrlPrefix:    "admin",
		InfoLogOff:   true,
		ErrorLogOff:  true,
		AccessLogOff: true,
	})
	os.Exit(m.Run())
}

func TestDefaultInvokerAuthFailCallback(t *testing.T) {
	const requestURI = "/admin/info/normal_manager?filter=%22%3Cscript%3E&sort=name"
	loginURL := "/admin/login?ref=" + url.QueryEscape(requestURI)
	invoker := DefaultInvoker(nil)

	for _, tc := range []struct {
		name    string
		method  string
		pjax    bool
		cookie  bool
		referer string
	}{
		{name: "navigation_without_cookie", method: http.MethodGet},
		{name: "navigation_with_stale_cookie", method: http.MethodGet, cookie: true, referer: "https://other.example/previous"},
		{name: "pjax_without_cookie", method: http.MethodGet, pjax: true},
		{name: "pjax_with_stale_cookie", method: http.MethodGet, pjax: true, cookie: true, referer: "/admin/previous"},
		{name: "post_with_stale_cookie", method: http.MethodPost, cookie: true},
		{name: "pjax_post_without_cookie", method: http.MethodPost, pjax: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, requestURI, nil)
			if tc.pjax {
				req.Header.Set(constant.PjaxHeader, "true")
			}
			if tc.cookie {
				req.AddCookie(&http.Cookie{Name: DefaultCookieKey, Value: "expired-session"})
			}
			if tc.referer != "" {
				req.Header.Set("Referer", tc.referer)
			}
			ctx := context.NewContext(req)
			invoker.authFailCallback(ctx)

			cookies := ctx.Response.Cookies()
			require.Len(t, cookies, 1)
			assert.Equal(t, DefaultCookieKey, cookies[0].Name)
			assert.Empty(t, cookies[0].Value)
			assert.Equal(t, "/", cookies[0].Path)
			assert.Equal(t, -1, cookies[0].MaxAge)
			assert.True(t, cookies[0].Expires.Before(time.Now()))
			assert.True(t, cookies[0].HttpOnly)

			defer ctx.Response.Body.Close()
			body, err := io.ReadAll(ctx.Response.Body)
			require.NoError(t, err)
			if tc.pjax {
				assert.Equal(t, http.StatusOK, ctx.Response.StatusCode)
				assert.Equal(t, "text/html; charset=utf-8", ctx.Response.Header.Get("Content-Type"))
				assert.Empty(t, ctx.Response.Header.Get("Location"))
				assert.Equal(t, "<script>\n\twindow.location.replace(\""+loginURL+"\")\n</script>", string(body))
			} else {
				assert.Equal(t, http.StatusFound, ctx.Response.StatusCode)
				assert.Equal(t, loginURL, ctx.Response.Header.Get("Location"))
				assert.Empty(t, body)
			}
		})
	}
}

func TestCheckPermissions(t *testing.T) {

	user := models.UserModel{
		Permissions: []models.PermissionModel{
			{
				Name:       "/",
				Slug:       "/",
				HttpMethod: []string{"GET"},
				HttpPath:   []string{"/"},
			}, {
				Name:       "/info/user",
				Slug:       "/",
				HttpMethod: []string{"GET"},
				HttpPath:   []string{"/info/user"},
			}, {
				Name:       "/info/user/edit",
				Slug:       "/",
				HttpMethod: []string{"GET"},
				HttpPath:   []string{"/info/user/edit"},
			}, {
				Name:       "/info/normal_manager?id=2",
				Slug:       "/",
				HttpMethod: []string{"GET"},
				HttpPath:   []string{"/info/normal_manager?id=2"},
			}, {
				Name:       "/info/normal_manager/edit?id=2",
				Slug:       "/",
				HttpMethod: []string{"GET"},
				HttpPath:   []string{"/info/normal_manager/edit?id=2"},
			}, {
				Name:       "/info/user_list?user_type=10",
				Slug:       "/",
				HttpMethod: []string{"GET"},
				HttpPath:   []string{"/info/user_list?user_type=10"},
			}, {
				Name:       "/info/user_list?user_type=20",
				Slug:       "/",
				HttpMethod: []string{"GET"},
				HttpPath:   []string{"/info/user_list?user_type=20"},
			}, {
				Name:       "/delete/user",
				Slug:       "/",
				HttpMethod: []string{"POST"},
				HttpPath:   []string{"/delete/user"},
			},
		},
	}

	param := make(url.Values)

	assert.Equal(t, CheckPermissions(user, "/admin/", "GET", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin", "GET", param), true)
	assert.Equal(t, CheckPermissions(user, "/", "GET", param), false)
	assert.Equal(t, CheckPermissions(user, "/admin", "POST", param), false)
	assert.Equal(t, CheckPermissions(user, "/admin/info/users", "GET", param), false)
	assert.Equal(t, CheckPermissions(user, "/admin/info/user", "GET", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin/info/user", "get", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin/info/normal_manager/edit?__goadmin_edit_pk=2&__columns=id,roles,created_at,updated_at", "get", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin/info/normal_manager/edit?__goadmin_edit_pk=2", "get", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin/info/normal_manager/edit?__goadmin_edit_pk=3&__columns=id,roles,created_at,updated_at", "get", param), false)
	assert.Equal(t, CheckPermissions(user, "/admin/info/normal_manager/edit?__columns=id,roles,created_at,updated_at&id=3", "get", param), false)
	assert.Equal(t, CheckPermissions(user, "/admin/info/user", "post", param), false)
	assert.Equal(t, CheckPermissions(user, "/admin/info/user/edit?id=3", "get", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin/logout?j=asdf", "post", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin/info/user_list?user_type=20", "get", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin/info/user_list?__goadmin_edit_pk=3&user_type=20", "get", param), true)
	assert.Equal(t, CheckPermissions(user, "/admin/delete/user", "post", param), true)
}
