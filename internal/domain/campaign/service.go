package campaign

import (
	"emailn/internal/domain/campaign/contract"
	internalerror "emailn/internal/internalError"
	"errors"
)

type ServiceImp struct {
	Repository Repository
	SendMail   func(campaign *Campaign) error
}

type Service interface {
	Create(newCampaign contract.NewCampaign) (string, error)
	GetBy(id string) (*contract.CampaignResponse, error)
	Delete(id string) error
	Start(id string) error
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaign) (string, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)

	if err != nil {
		return "", err
	}
	err = s.Repository.Create(campaign)
	if err != nil {
		return "", internalerror.ErrInternal
	}
	return campaign.ID, nil
}

func (s *ServiceImp) GetBy(id string) (*contract.CampaignResponse, error) {

	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return nil, internalerror.ProcessErrorToReturn(err)
	}

	return &contract.CampaignResponse{
		ID:                   campaign.ID,
		Name:                 campaign.Name,
		Content:              campaign.Content,
		Status:               campaign.Status,
		AmountOfEmailsToSend: len(campaign.Contacts),
		CreatedBy:            campaign.CreatedBy,
	}, nil
}

func (s *ServiceImp) Delete(id string) error {

	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return internalerror.ProcessErrorToReturn(err)
	}

	if campaign.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	campaign.Delete()

	err = s.Repository.Delete(campaign)
	if err != nil {
		return internalerror.ErrInternal
	}

	return nil
}

func (s *ServiceImp) Start(id string) error {

	campaignSaved, err := s.Repository.GetBy(id)

	if err != nil {
		return internalerror.ProcessErrorToReturn(err)
	}

	if campaignSaved.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	err = s.SendMail(campaignSaved)
	if err != nil {
		return internalerror.ErrInternal
	}

	campaignSaved.Done()
	err = s.Repository.Update(campaignSaved)
	if err != nil {
		return internalerror.ErrInternal
	}

	return nil
}
