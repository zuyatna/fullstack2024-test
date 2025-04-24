package usecase

import (
	"fullstack2024-test/model"
	"fullstack2024-test/repository"
)

type ClientUseCase struct {
	clientRepository repository.IClientRepository
}

func NewClientUseCase(clientRepository repository.IClientRepository) *ClientUseCase {
	return &ClientUseCase{clientRepository: clientRepository}
}

func (clientUseCase *ClientUseCase) CreateClient(client *model.Client) (*model.Client, error) {
	return clientUseCase.clientRepository.CreateClient(client)
}

func (clientUseCase *ClientUseCase) GetClientByID(id int) (*model.Client, error) {
	return clientUseCase.clientRepository.GetClientByID(id)
}

func (clientUseCase *ClientUseCase) UpdateClientByID(id int, client *model.Client) (*model.Client, error) {
	return clientUseCase.clientRepository.UpdateClientByID(id, client)
}

func (clientUseCase *ClientUseCase) DeleteClientByID(id int) error {
	return clientUseCase.clientRepository.DeleteClientByID(id)
}
