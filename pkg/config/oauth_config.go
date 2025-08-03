package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthProviderConfig struct {
	Config      *oauth2.Config
	UserInfoURL string
}

var OauthProviders = map[string]OAuthProviderConfig{
	"google": {
		Config: &oauth2.Config{
			ClientID:     "GOOGLE_CLIENT_ID",     // <-- подставь из .env
			ClientSecret: "GOOGLE_CLIENT_SECRET", // <-- подставь из .env
			RedirectURL:  "https://your-domain.com/api/v1/auth/callback/google",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		UserInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
	},
	"yandex": {
		Config: &oauth2.Config{
			ClientID:     "YANDEX_CLIENT_ID",
			ClientSecret: "YANDEX_CLIENT_SECRET",
			RedirectURL:  "https://your-domain.com/api/v1/auth/callback/yandex",
			Scopes:       []string{"login:email", "login:info"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://oauth.yandex.ru/authorize",
				TokenURL: "https://oauth.yandex.ru/token",
			},
		},
		UserInfoURL: "https://login.yandex.ru/info?format=json",
	},
	"vk": {
		Config: &oauth2.Config{
			ClientID:     "VK_CLIENT_ID",
			ClientSecret: "VK_CLIENT_SECRET",
			RedirectURL:  "https://your-domain.com/api/v1/auth/callback/vk",
			Scopes:       []string{"email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://oauth.vk.com/authorize",
				TokenURL: "https://oauth.vk.com/access_token",
			},
		},
		UserInfoURL: "https://api.vk.com/method/users.get?fields=uid,first_name,last_name&v=5.131",
	},
}
