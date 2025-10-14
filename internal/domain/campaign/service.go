package campaign

import (
	"emailn/internal/domain/campaign/contract"
	internalerror "emailn/internal/internalError"
)

type Service struct {
	Repository Repository
}

func (s *Service) Create(newCampaign contract.NewCampaign) (string, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)

	if err != nil {
		return "", err
	}
	err = s.Repository.Save(campaign)
	if err != nil {
		return "", internalerror.ErrInternal
	}
	s.Repository.Save(campaign)
	return campaign.ID, nil
}
