package validation

import (
	"encoding/json"
	"net/http"

	"github.com/JubaerHossain/rootx/pkg/utils"
	"github.com/go-playground/validator/v10"
)

func BodyParse(s interface{}, w http.ResponseWriter, r *http.Request, isValidation bool) error {
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return err
	}

	if isValidation {
		validate := validator.New()
		validateErr := validate.Struct(s)
		if validateErr != nil {
			utils.ValidationResponse(w, http.StatusBadRequest, validateErr.(validator.ValidationErrors), s)
			return validateErr
		}
	}
	return nil
}
