package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type GolfType struct {
	GolfID    int     `json:"id"`
	Golfname  string  `json:"coursename"`
	Price     float64 `json:"price"`
	Totalhole string  `json:"totalhole"`
}

const coursePath = "golfcourses"

var golfList []GolfType

var Db *sql.DB

const basePath = "/api"

func getCourse(courseid int) (*GolfType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := Db.QueryRowContext(ctx, `
		SELECT
			courseid,
			golfcoursename,
			price,
			totalhole
		FROM golfcourse
		WHERE courseid = ?`, courseid)

	course := &GolfType{}
	err := row.Scan(
		&course.GolfID,
		&course.Golfname,
		&course.Price,
		&course.Totalhole,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil
}

func removeCourse(courseID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := Db.ExecContext(ctx, `DELETE FROM golfcourse WHERE courseid = ?`, courseID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func getCourseList() ([]GolfType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := Db.QueryContext(ctx, `
		SELECT
			courseid,
			golfcoursename,
			price,
			totalhole
		FROM golfcourse`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	golfcourses := make([]GolfType, 0)
	for results.Next() {
		var course GolfType
		err := results.Scan(&course.GolfID, &course.Golfname, &course.Price, &course.Totalhole)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		golfcourses = append(golfcourses, course)
	}
	return golfcourses, nil
}

func insertGolfCourse(course GolfType) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := Db.ExecContext(ctx, `
		INSERT INTO golfcourse (courseid, golfcoursename, price, totalhole)
		VALUES (?, ?, ?, ?)`,
		course.GolfID,
		course.Golfname,
		course.Price,
		course.Totalhole,
	)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func updateGolfCourse(courseID int, course GolfType) (*GolfType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := Db.ExecContext(ctx, `
		UPDATE golfcourse
		SET golfcoursename = ?, price = ?, totalhole = ?
		WHERE courseid = ?`,
		course.Golfname,
		course.Price,
		course.Totalhole,
		courseID,
	)
	if err != nil {
		log.Println("Update failed:", err)
		return nil, err
	}

	return &course, nil
}

func handleCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseList, err := getCourseList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(courseList)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(j)

	case http.MethodPost:
		var course GolfType
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		CourseID, err := insertGolfCourse(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"courseid":%d}`, CourseID)))

	case http.MethodOptions:
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleCourse(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/golfcourses/")
	courseID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print("Invalid course ID:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:

		course, err := getCourse(courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write(j)

	case http.MethodDelete:
		err := removeCourse(courseID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		var course GolfType
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedCourse, err := updateGolfCourse(courseID, course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// ส่งผลลัพธ์การอัพเดตกลับ
		j, err := json.Marshal(updatedCourse)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(j)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		handler.ServeHTTP(w, r)
	})
}

func SetupDB() {
	var err error
	Db, err = sql.Open("mysql", "root:p@ssw0rd_11@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	log.Println("Connected to the MySQL database successfully!")

	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}

func SetupRoutes(apiBasePath string) {
	courseHandler := http.HandlerFunc(handleCourse)
	coursesHandler := http.HandlerFunc(handleCourses)

	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler)) // for /api/golfcourses/2
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler)) // for /api/golfcourses
}

func main() {
	SetupDB()
	SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
