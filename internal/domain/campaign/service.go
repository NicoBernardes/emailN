package campaign

import "emailn/internal/domain/campaign/contract"

type Service struct {
	Repository Repository
}

func (s *Service) Create(newCampaign contract.NewCampaign) error {

	return nil
}
