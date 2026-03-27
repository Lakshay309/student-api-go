package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	myTypes "github.com/Lakshay309/student-api-go/internal/types"
	"github.com/Lakshay309/student-api-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		slog.Info("creating a student");
		var student myTypes.Student;
		err:=json.NewDecoder(r.Body).Decode(&student); 
		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("Empty Body")))
			return
		}
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		// w.Write([]byte("Welcome to student api"))

		// request validation
		if err:=validator.New().Struct(student); err!=nil{
			validateErrs:=err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError(validateErrs))
			return 
		}


		response.WriteJson(w,http.StatusCreated,map[string]string{"sucess":"ok"})
	}
}

