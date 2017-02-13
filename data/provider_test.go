package data

import (
	"database/sql"
	"testing"
)

func Test_ProviderCreate(t *testing.T) {
	setup()
	if err := providers[0].Create(); err != nil {
		t.Error(err, "Cannot create provider.")
	}
	if providers[0].Id == 0 {
		t.Errorf("No id provider.")
	}
	p, err := ProviderByID(providers[0].Id)
	if err != nil {
		t.Error(err, "Provider not created.")
	}
	if providers[0].Name != p.Name {
		t.Errorf("Provider retrieved is not the same as the one created.")
	}
}

func Test_ProviderDelete(t *testing.T) {
	setup()
	if err := providers[0].Create(); err != nil {
		t.Error(err, "Cannot create provider.")
	}
	if err := providers[0].Delete(); err != nil {
		t.Error(err, "Cannot delete provider")
	}
	_, err := ProviderByID(providers[0].Id)
	if err != sql.ErrNoRows {
		t.Error(err, "Provider not deleted.")
	}
}

func Test_ProviderUpdate(t *testing.T) {
	setup()
	if err := providers[0].Create(); err != nil {
		t.Error(err, "Cannot create provider.")
	}
	providers[0].Name = "Random Provider"
	if err := providers[0].Update(); err != nil {
		t.Error(err, "Cannot update provider.")
	}
	p, err := ProviderByID(providers[0].Id)
	if err != nil {
		t.Error(err, "Cannot get provider.")
	}
	if p.Name != "Random Provider" {
		t.Error(err, "Provider not updated.")
	}
}

func Test_Providers(t *testing.T) {
	setup()
	for _, p := range providers {
		if err := p.Create(); err != nil {
			t.Error(err, "Cannot create provider.")
		}
	}
	ps, err := Providers()
	if err != nil {
		t.Error(err, "Cannot retrieve providers.")
	}
	if len(ps) != 2 {
		t.Error(err, "Wrong number of provider retrieved")
	}
	if ps[0].Name != providers[0].Name {
		t.Error(ps[0], providers[0], "Wrong provider retrieved")
	}
}
