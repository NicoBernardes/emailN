package campaign

import "testing"

func TestNewCampaign(t *testing.T) {
	name := "Campaign X"
	content := "Body"
	contacts := []string{"email1@e.com", "email2@e.com"}

	campaign := NewCampaign(name, content, contacts)

	if campaign.ID != "1" {
		t.Errorf("Expected ID 1")
	} else if campaign.Name != name {
		t.Errorf("Expected correct name")
	} else if campaign.Content != content {
		t.Errorf("Expected correct content")
	} else if len(campaign.Contacts) != len(contacts) {
		t.Errorf("Expected correct contacts")
	}
}
