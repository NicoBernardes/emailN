package campaign

import (
	internalerror "emailn/internal/internalError"
	"time"

	"github.com/rs/xid"
)

type Contact struct {
	ID         string `gorm:"size:50"`
	Email      string `validate:"email" gorm:"size:100"`
	CampaignId string `gorm:"size:50"`
}

type Campaign struct {
	ID        string    `validate:"required" gorm:"size:50;not null"`
	Name      string    `validate:"min=5,max=24" gorm:"size:100;not null"`
	CreatedOn time.Time `validate:"required" gorm:"not null"`
	UpdatedOn time.Time
	Content   string    `validate:"min=5,max=1024" gorm:"size:1024;not null"`
	Contacts  []Contact `validate:"min=1,dive"` //dive to another struct
	Status    string    `gorm:"size:20;not null"`
	CreatedBy string    `validate:"email" gorm:"size:50;not null"`
}

const (
	Pending  = "Pending"
	Canceled = "Canceled"
	Deleted  = "Deleted"
	Started  = "Started"
	Fail     = "Fail"
	Done     = "Done"
)

func (c *Campaign) Done() {
	c.Status = Done
	c.UpdatedOn = time.Now()
}

func (c *Campaign) Cancel() {
	c.Status = Canceled
	c.UpdatedOn = time.Now()
}

func (c *Campaign) Delete() {
	c.Status = Deleted
	c.UpdatedOn = time.Now()
}

func (c *Campaign) Fail() {
	c.Status = Fail
	c.UpdatedOn = time.Now()
}

func (c *Campaign) Started() {
	c.Status = Started
	c.UpdatedOn = time.Now()
}

func NewCampaign(name, content string, emails []string, createdBy string) (*Campaign, error) {

	contacts := make([]Contact, len(emails))

	for index, email := range emails {
		contacts[index].Email = email
		contacts[index].ID = xid.New().String()
	}

	campaign := &Campaign{
		ID:        xid.New().String(),
		Name:      name,
		Content:   content,
		CreatedOn: time.Now(),
		Contacts:  contacts,
		Status:    Pending,
		CreatedBy: createdBy,
	}
	err := internalerror.ValidateStruct(campaign)
	if err == nil {
		return campaign, nil
	}
	return nil, err
}
