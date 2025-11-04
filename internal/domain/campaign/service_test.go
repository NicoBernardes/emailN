package campaign_test

import (
	"emailn/internal/domain/campaign"
	internalerror "emailn/internal/internalError"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaign = campaign.NewCampaignRequest{
		Name:      "Teste X",
		Content:   "Body Hi!",
		Emails:    []string{"teste1@teste.com"},
		CreatedBy: "teste@teste.com.br",
	}
	campaignPending *campaign.Campaign
	campaignStarted *campaign.Campaign
	repositoryMock  *internalmock.CampaignRepositoryMock
	service         = campaign.ServiceImp{}
)

func setUp() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	campaignPending, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	campaignStarted = &campaign.Campaign{ID: "1", Status: campaign.Started}
}

func setUpGetByIdRepositoryBy(campaign *campaign.Campaign) {
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
}

func setUpSendEmail(err error) {
	sendMail := func(campaign *campaign.Campaign) error {
		return err
	}
	service.SendMail = sendMail
}

func Test_Create_Campaign(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)

	assert.NotNil(t, id)
	assert.Nil(t, err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	setUp()

	_, err := service.Create(campaign.NewCampaignRequest{})

	assert.False(t, errors.Is(err, internalerror.ErrInternal))

}

func Test_Create_SaveCampaign(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)

}

func Test_Create_ValidadeRepositorySave(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(t, errors.Is(err, internalerror.ErrInternal))

}

func Test_GetByIdReturnCampaign(t *testing.T) {
	setUp()

	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaignPending.ID
	})).Return(campaignPending, nil)

	campaignReturned, _ := service.GetBy(campaignPending.ID)

	assert.Equal(t, campaignPending.ID, campaignReturned.ID)
	assert.Equal(t, campaignPending.Name, campaignReturned.Name)
	assert.Equal(t, campaignPending.Content, campaignReturned.Content)
	assert.Equal(t, campaignPending.Status, campaignReturned.Status)
	assert.Equal(t, campaignPending.CreatedBy, campaignReturned.CreatedBy)
}

func Test_GetByIdReturnErrorWhenSomethingWrong(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))
	_, err := service.GetBy("invalid_campaign")

	assert.Equal(t, internalerror.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnRecordNotFound_when_campaign_does_not_exists(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid_campaign")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_ReturnStatusInvalid_when_campaign_has_status_not_pending(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Delete(campaignStarted.ID)

	assert.Equal(t, "Campaign status invalid", err.Error())
}

func Test_Delete_ReturnInternalError_when_delete_has_problem(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignPending.ID)

	assert.Equal(t, internalerror.ErrInternal.Error(), err.Error())
}

func Test_Delete_ReturnNil_when_delete_success(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPending == campaign
	})).Return(nil)

	err := service.Delete(campaignPending.ID)

	assert.Nil(t, err)
}

func Test_Start_ReturnRecordNotFound_when_campaign_does_not_exists(t *testing.T) {
	setUp()
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("invalid_campaign")

	assert.Equal(t, err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_ReturnStatusInvalid_when_campaign_has_status_not_pending(t *testing.T) {
	setUp()
	setUpGetByIdRepositoryBy(campaignStarted)
	repositoryMock.On("Update", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	err := service.Start(campaignStarted.ID)
	assert.Nil(t, err)
}

func Test_Start_CampaignWasUpdated_StatusIsStarted(t *testing.T) {
	setUp()
	setUpSendEmail(nil)
	setUpGetByIdRepositoryBy(campaignPending)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Started
	})).Return(nil)
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignPending.ID)

	assert.Equal(t, campaign.Started, campaignPending.Status)
}

func Test_SendEmailAndUpdateStatus_WhenFail_StatusIsFail(t *testing.T) {
	setUp()
	setUpSendEmail(errors.New("error to send email"))
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Fail
	})).Return(nil)

	service.SendEmailAndUpdateStatus(campaignPending)

	repositoryMock.AssertExpectations(t)

}

func Test_SendEmailAndUpdateStatus_WhenFail_StatusIsDone(t *testing.T) {
	setUp()
	setUpSendEmail(nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)

	service.SendEmailAndUpdateStatus(campaignPending)

	repositoryMock.AssertExpectations(t)

}
