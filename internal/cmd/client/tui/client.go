package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

var nilPetClient = logMsg{Source: "client", JSON: `{"msg": "PetServiceClient is nil"}`}

type Client struct {
	petCl gostreamv1.PetServiceClient
	ctx   context.Context
}

func (c Client) DeletePet(id uint64) tea.Cmd {
	return func() tea.Msg {
		if c.petCl == nil {
			return nilPetClient
		}

		_, err := c.petCl.DeletePet(c.ctx, &gostreamv1.DeletePetRequest{
			PetId: id,
		})
		if err != nil {
			return logMsg{Source: "client", JSON: fmt.Sprintf(`{"msg": "Failed to delete pet %d"}`, id)}
		}

		return nil
	}

}

func (c Client) CreatePet(name, kind string, age uint32) tea.Cmd {
	return func() tea.Msg {
		if c.petCl == nil {
			return nilPetClient
		}

		_, err := c.petCl.CreatePet(c.ctx, &gostreamv1.CreatePetRequest{
			Pet: &gostreamv1.Pet{
				Age:  age,
				Name: name,
				Kind: kind,
			},
		})

		if err != nil {
			return logMsg{Source: "client", JSON: fmt.Sprintf(`{"msg": "Failed to create new pet: %v"}`, name)}
		}

		return nil
	}
}
