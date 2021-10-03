package site24x7

import (
	"testing"

	"github.com/Site24x7/terraform-provider-site24x7/api"
	apierrors "github.com/Site24x7/terraform-provider-site24x7/api/errors"
	"github.com/Site24x7/terraform-provider-site24x7/fake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsiteMonitorCreate(t *testing.T) {
	d := websiteMonitorTestResourceData(t)

	c := fake.NewClient()

	a := &api.WebsiteMonitor{
		DisplayName:           "foo",
		Type:                  "URL",
		Website:               "www.test.tld",
		CheckFrequency:        "5",
		HTTPMethod:            "G",
		Timeout:               10,
		LocationProfileID:     "456",
		NotificationProfileID: "789",
		ThresholdProfileID:    "012",
		UseNameServer:         true,
		MatchCase:             true,
		UserAgent:             "firefox",
		MonitorGroups:         []string{"234", "567"},
		UserGroupIDs:          []string{"123", "456"},
		AuthUser:              "username",
		AuthPass:              "password",
		CustomHeaders: []api.Header{
			{
				Name:  "Accept",
				Value: "application/json",
			},
			{
				Name:  "Cache-Control",
				Value: "nocache",
			},
		},
		ActionIDs: []api.ActionRef{
			{
				ActionID:  "123action",
				AlertType: 1,
			},
			{
				ActionID:  "234action",
				AlertType: 5,
			},
		},
		// UnmatchingKeyword: &api.ValueAndSeverity{
		// 	Value:    "foo",
		// 	Severity: 2,
		// },
		// MatchingKeyword: &api.ValueAndSeverity{
		// 	Value:    "bar",
		// 	Severity: 2,
		// },
		// MatchRegex: &api.ValueAndSeverity{
		// 	Value:    ".*",
		// 	Severity: 2,
		// },
	}

	c.FakeWebsiteMonitors.On("Create", a).Return(a, nil).Once()

	require.NoError(t, websiteMonitorCreate(d, c))

	c.FakeWebsiteMonitors.On("Create", a).Return(a, apierrors.NewStatusError(500, "error")).Once()

	err := websiteMonitorCreate(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestWebsiteMonitorUpdate(t *testing.T) {
	d := websiteMonitorTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	a := &api.WebsiteMonitor{
		MonitorID:             "123",
		DisplayName:           "foo",
		Type:                  "URL",
		Website:               "www.test.tld",
		CheckFrequency:        "5",
		HTTPMethod:            "G",
		Timeout:               10,
		LocationProfileID:     "456",
		NotificationProfileID: "789",
		ThresholdProfileID:    "012",
		UseNameServer:         true,
		MatchCase:             true,
		UserAgent:             "firefox",
		MonitorGroups:         []string{"234", "567"},
		UserGroupIDs:          []string{"123", "456"},
		AuthUser:              "username",
		AuthPass:              "password",
		CustomHeaders: []api.Header{
			{
				Name:  "Accept",
				Value: "application/json",
			},
			{
				Name:  "Cache-Control",
				Value: "nocache",
			},
		},
		ActionIDs: []api.ActionRef{
			{
				ActionID:  "123action",
				AlertType: 1,
			},
			{
				ActionID:  "234action",
				AlertType: 5,
			},
		},
		// UnmatchingKeyword: &api.ValueAndSeverity{
		// 	Value:    "foo",
		// 	Severity: 2,
		// },
		// MatchingKeyword: &api.ValueAndSeverity{
		// 	Value:    "bar",
		// 	Severity: 2,
		// },
		// MatchRegex: &api.ValueAndSeverity{
		// 	Value:    ".*",
		// 	Severity: 2,
		// },
	}

	c.FakeWebsiteMonitors.On("Update", a).Return(a, nil).Once()

	require.NoError(t, websiteMonitorUpdate(d, c))

	c.FakeWebsiteMonitors.On("Update", a).Return(a, apierrors.NewStatusError(500, "error")).Once()

	err := websiteMonitorUpdate(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestWebsiteMonitorRead(t *testing.T) {
	d := websiteMonitorTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeWebsiteMonitors.On("Get", "123").Return(&api.WebsiteMonitor{}, nil).Once()

	require.NoError(t, websiteMonitorRead(d, c))

	c.FakeWebsiteMonitors.On("Get", "123").Return(nil, apierrors.NewStatusError(500, "error")).Once()

	err := websiteMonitorRead(d, c)

	assert.Equal(t, apierrors.NewStatusError(500, "error"), err)
}

func TestWebsiteMonitorDelete(t *testing.T) {
	d := websiteMonitorTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeWebsiteMonitors.On("Delete", "123").Return(nil).Once()

	require.NoError(t, websiteMonitorDelete(d, c))

	c.FakeWebsiteMonitors.On("Delete", "123").Return(apierrors.NewStatusError(404, "not found")).Once()

	require.NoError(t, websiteMonitorDelete(d, c))
}

func TestWebsiteMonitorExists(t *testing.T) {
	d := websiteMonitorTestResourceData(t)
	d.SetId("123")

	c := fake.NewClient()

	c.FakeWebsiteMonitors.On("Get", "123").Return(&api.WebsiteMonitor{}, nil).Once()

	exists, err := websiteMonitorExists(d, c)

	require.NoError(t, err)
	assert.True(t, exists)

	c.FakeWebsiteMonitors.On("Get", "123").Return(nil, apierrors.NewStatusError(404, "not found")).Once()

	exists, err = websiteMonitorExists(d, c)

	require.NoError(t, err)
	assert.False(t, exists)

	c.FakeWebsiteMonitors.On("Get", "123").Return(nil, apierrors.NewStatusError(500, "error")).Once()

	exists, err = websiteMonitorExists(d, c)

	require.Equal(t, apierrors.NewStatusError(500, "error"), err)
	assert.False(t, exists)
}

func websiteMonitorTestResourceData(t *testing.T) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, WebsiteMonitorSchema, map[string]interface{}{
		"display_name":    "foo",
		"type":            "URL",
		"website":         "www.test.tld",
		"check_frequency": "5",
		"http_method":     "G",
		"auth_user":       "username",
		"auth_pass":       "password",
		"match_case":      true,
		"user_agent":      "firefox",
		"custom_headers": map[string]interface{}{
			"Accept":        "application/json",
			"Cache-Control": "nocache",
		},
		"timeout":                 10,
		"location_profile_id":     "456",
		"notification_profile_id": "789",
		"threshold_profile_id":    "012",
		"monitor_groups": []interface{}{
			"234",
			"567",
		},
		"user_group_ids": []interface{}{
			"123",
			"456",
		},
		"use_name_server": true,
		"actions": map[string]interface{}{
			"1": "123action",
			"5": "234action",
		},
	})
}
