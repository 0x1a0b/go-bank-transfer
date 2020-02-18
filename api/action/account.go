package action

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gsabadini/go-stone/domain"
	"github.com/gsabadini/go-stone/infrastructure/database"
	"github.com/gsabadini/go-stone/repository"
	"github.com/gsabadini/go-stone/usecase"
)

type Account struct {
	dbHandler database.NoSQLDBHandler
}

func NewAccount(dbHandler database.NoSQLDBHandler) Account {
	return Account{dbHandler: dbHandler}
}

func (a Account) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var account domain.Account
	if err = json.Unmarshal(reqBody, &account); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var accountRepository = repository.NewAccount(a.dbHandler)
	err = usecase.Create(accountRepository, account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Success"))
}
