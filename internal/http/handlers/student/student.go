package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Lakshay309/student-api-go/internal/storage"
	myTypes "github.com/Lakshay309/student-api-go/internal/types"
	"github.com/Lakshay309/student-api-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc{
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

		lastId,err:=storage.CreateStudent(student.Name,student.Email,student.Age)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
		}
		slog.Info("user created successfully",slog.String("userId",fmt.Sprint(lastId)))

		response.WriteJson(w,http.StatusCreated,map[string]int64{"id":lastId})
	}
}

func GetById(storage storage.Storage)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		id:=r.PathValue("id")
		slog.Info("getting a student",slog.String("id",id))
		intId,err:=strconv.ParseInt(id,10,64);
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		student,err:=storage.GetStudentById(intId)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		response.WriteJson(w,http.StatusFound,student)
	}
}