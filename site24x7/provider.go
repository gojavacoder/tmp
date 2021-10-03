package site24x7

import (
	"os"
	"time"

	"github.com/Site24x7/terraform-provider-site24x7/backoff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	log "github.com/sirupsen/logrus"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"oauth2_client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SITE24X7_OAUTH2_CLIENT_ID", nil),
				Description: "OAuth2 Client ID",
			},
			"oauth2_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SITE24X7_OAUTH2_CLIENT_SECRET", nil),
				Description: "OAuth2 Client Secret",
			},
			"oauth2_refresh_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SITE24X7_OAUTH2_REFRESH_TOKEN", nil),
				Description: "OAuth2 Refresh Token",
			},
			"data_center": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Site24x7 data center.",
			},
			"retry_min_wait": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Minimum wait time in seconds before retrying failed API requests.",
			},
			"retry_max_wait": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Maximum wait time in seconds before retrying failed API requests (exponential backoff).",
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4,
				Description: "Maximum number of retries for Site24x7 API errors until giving up",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"site24x7_website_monitor":   resourceSite24x7WebsiteMonitor(),
			"site24x7_ssl_monitor":       resourceSite24x7SSLMonitor(),
			"site24x7_rest_api_monitor":  resourceSite24x7RestApiMonitor(),
			"site24x7_amazon_monitor":    resourceSite24x7AmazonMonitor(),
			"site24x7_monitor_group":     resourceSite24x7MonitorGroup(),
			"site24x7_url_action":        resourceSite24x7URLAction(),
			"site24x7_threshold_profile": resourceSite24x7ThresholdProfile(),
			"site24x7_user_group":        resourceSite24x7UserGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	tfLog := os.Getenv("TF_LOG")
	if tfLog == "DEBUG" || tfLog == "TRACE" {
		log.SetLevel(log.DebugLevel)
	}
	dataCenter := GetDataCenter(d.Get("data_center").(string))
	log.Println("dataCenter GetAPIBaseURL : ", dataCenter.GetAPIBaseURL())
	log.Println("dataCenter GetTokenURL : ", dataCenter.GetTokenURL())
	config := Config{
		ClientID:     d.Get("oauth2_client_id").(string),
		ClientSecret: d.Get("oauth2_client_secret").(string),
		RefreshToken: d.Get("oauth2_refresh_token").(string),
		APIBaseURL:   dataCenter.GetAPIBaseURL(),
		TokenURL:     dataCenter.GetTokenURL(),
		RetryConfig: &backoff.RetryConfig{
			MinWait:    time.Duration(d.Get("retry_min_wait").(int)) * time.Second,
			MaxWait:    time.Duration(d.Get("retry_max_wait").(int)) * time.Second,
			MaxRetries: d.Get("max_retries").(int),
		},
	}

	return New(config), nil
}
