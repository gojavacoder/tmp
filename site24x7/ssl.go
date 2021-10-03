package site24x7

import (
	"github.com/Site24x7/terraform-provider-site24x7/api"
	apierrors "github.com/Site24x7/terraform-provider-site24x7/api/errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var SSLMonitorSchema = map[string]*schema.Schema{
	"display_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"domain_name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"timeout": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  30,
	},
	"protocol": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "HTTPS",
	},
	"port": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  443,
	},
	"expire_days": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  30,
	},
	"http_protocol_version": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "H1.1",
	},
	"ignore_domain_mismatch": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"ignore_trust": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"location_profile_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"location_profile_name": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"notification_profile_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"threshold_profile_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"monitor_groups": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	"user_group_ids": {
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
		Computed: true,
	},
}

func resourceSite24x7SSLMonitor() *schema.Resource {
	return &schema.Resource{
		Create: sslMonitorCreate,
		Read:   sslMonitorRead,
		Update: sslMonitorUpdate,
		Delete: sslMonitorDelete,
		Exists: sslMonitorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: SSLMonitorSchema,
	}
}

func sslMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)

	sslMonitor, err := resourceDataToSSLMonitor(d, client)
	if err != nil {
		return err
	}

	sslMonitor, err = client.SSLMonitors().Create(sslMonitor)
	if err != nil {
		return err
	}

	d.SetId(sslMonitor.MonitorID)

	return nil
}

func sslMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)

	sslMonitor, err := client.SSLMonitors().Get(d.Id())
	if err != nil {
		return err
	}

	updateSSLMonitorResourceData(d, sslMonitor)

	return nil
}

func sslMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)

	sslMonitor, err := resourceDataToSSLMonitor(d, client)
	if err != nil {
		return err
	}

	sslMonitor, err = client.SSLMonitors().Update(sslMonitor)
	if err != nil {
		return err
	}

	d.SetId(sslMonitor.MonitorID)

	// return sslMonitorRead(d, meta)
	return nil
}

func sslMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)

	err := client.SSLMonitors().Delete(d.Id())
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

func sslMonitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(Client)

	_, err := client.SSLMonitors().Get(d.Id())
	if apierrors.IsNotFound(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func resourceDataToSSLMonitor(d *schema.ResourceData, client Client) (*api.SSLMonitor, error) {

	var userGroupIDs []string
	for _, id := range d.Get("user_group_ids").([]interface{}) {
		userGroupIDs = append(userGroupIDs, id.(string))
	}

	var monitorGroups []string
	for _, group := range d.Get("monitor_groups").([]interface{}) {
		monitorGroups = append(monitorGroups, group.(string))
	}

	sslMonitor := &api.SSLMonitor{
		MonitorID:             d.Id(),
		DisplayName:           d.Get("display_name").(string),
		Type:                  string(api.SSL_CERT),
		DomainName:            d.Get("domain_name").(string),
		Protocol:              d.Get("protocol").(string),
		Timeout:               d.Get("timeout").(int),
		Port:                  d.Get("port"),
		ExpireDays:            d.Get("expire_days").(int),
		HTTPProtocolVersion:   d.Get("http_protocol_version").(string),
		IgnoreDomainMismatch:  d.Get("ignore_domain_mismatch").(bool),
		IgnoreTrust:           d.Get("ignore_trust").(bool),
		LocationProfileID:     d.Get("location_profile_id").(string),
		NotificationProfileID: d.Get("notification_profile_id").(string),
		ThresholdProfileID:    d.Get("threshold_profile_id").(string),
		MonitorGroups:         monitorGroups,
		UserGroupIDs:          userGroupIDs,
	}

	if sslMonitor.LocationProfileID == "" {
		locationProfileNameToMatch := d.Get("location_profile_name").(string)
		profile, err := DefaultLocationProfile(client, locationProfileNameToMatch)
		if err != nil {
			return nil, err
		}
		sslMonitor.LocationProfileID = profile.ProfileID
		d.Set("location_profile_id", profile.ProfileID)
	}

	if sslMonitor.NotificationProfileID == "" {
		profile, err := DefaultNotificationProfile(client)
		if err != nil {
			return nil, err
		}
		sslMonitor.NotificationProfileID = profile.ProfileID
		d.Set("notification_profile_id", profile.ProfileID)
	}

	if sslMonitor.ThresholdProfileID == "" {
		profile, err := DefaultThresholdProfile(client, api.SSL_CERT)
		if err != nil {
			return nil, err
		}
		sslMonitor.ThresholdProfileID = profile.ProfileID
		d.Set("threshold_profile_id", profile)
	}

	if len(sslMonitor.UserGroupIDs) == 0 {
		userGroup, err := DefaultUserGroup(client)
		if err != nil {
			return nil, err
		}
		sslMonitor.UserGroupIDs = []string{userGroup.UserGroupID}
		d.Set("user_group_ids", []string{userGroup.UserGroupID})
	}

	return sslMonitor, nil
}

func updateSSLMonitorResourceData(d *schema.ResourceData, monitor *api.SSLMonitor) {
	d.Set("display_name", monitor.DisplayName)
	d.Set("type", monitor.Type)
	d.Set("domain_name", monitor.DomainName)
	d.Set("timeout", monitor.Timeout)
	d.Set("protocol", monitor.Protocol)
	d.Set("port", monitor.Port)
	d.Set("expire_days", monitor.ExpireDays)
	d.Set("http_protocol_version", monitor.HTTPProtocolVersion)
	d.Set("ignore_domain_mismatch", monitor.IgnoreDomainMismatch)
	d.Set("ignore_trust", monitor.IgnoreTrust)
	d.Set("location_profile_id", monitor.LocationProfileID)
	d.Set("notification_profile_id", monitor.NotificationProfileID)
	d.Set("threshold_profile_id", monitor.ThresholdProfileID)
	d.Set("monitor_groups", monitor.MonitorGroups)
	d.Set("user_group_ids", monitor.UserGroupIDs)

}
