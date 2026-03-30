package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/Lakshay309/student-api-go/internal/config"
	myTypes "github.com/Lakshay309/student-api-go/internal/types"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config)(*Sqlite,error){
	db,err:=sql.Open("sqlite",cfg.StoragePath)
	if err!=nil{
		return nil,err;
	}

	_,err=db.Exec(`CREATE TABLE IF NOT EXISTS students(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)
	if err!=nil{
		return nil,err;
	}
	return &Sqlite{
		Db: db,
	},nil
}

func (s *Sqlite) CreateStudent(name string,email string,age int) (int64,error){
	stmt,err:=s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)");
	if err!=nil{
		return 0,err;
	}
	defer stmt.Close()
	result,err:=stmt.Exec(name,email,age);
	if err!=nil{
		return 0,err;
	}
	lastId,err:=result.LastInsertId()
	if err!=nil{
		return 0,err;
	}
	return lastId,nil;
}

func (s Sqlite) GetStudentById(id int64)( myTypes.Student,error){
	stmt,err:=s.Db.Prepare("SELECT * FROM students where id=? LIMIT 1")
	if err!=nil{
		return myTypes.Student{},err;
	}
	defer stmt.Close()

	var student myTypes.Student
	err=stmt.QueryRow(id).Scan(&student.Id,&student.Name,&student.Email,&student.Age)
	if err!=nil{
		if err== sql.ErrNoRows{
			return myTypes.Student{},fmt.Errorf("No sudent with that id: %s",fmt.Sprint(id))
		}
		return myTypes.Student{},fmt.Errorf("query Error: %w",err)
	}
	return student,nil;
}

func  (s Sqlite) GetStudentList()([]myTypes.Student,error){
	stmt,err:=s.Db.Prepare("SELECT * FROM students");
	if err!=nil{
		return []myTypes.Student{},err;
	}
	defer stmt.Close()
	var students []myTypes.Student;
	row,err:=stmt.Query();
	if err!=nil{
		return []myTypes.Student{},err;
	}
	defer row.Close()
	for row.Next(){
		var student myTypes.Student;
		if err:=row.Scan(&student.Id,&student.Name,&student.Email,&student.Age);err!=nil{
			return []myTypes.Student{},err;
		}
		students = append(students, student)
	}
	return students,nil;
}