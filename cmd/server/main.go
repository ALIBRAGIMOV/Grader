package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"grader/pkg/queue"
	loggerModel "grader/pkg/server/logger"
	"grader/pkg/server/middleware"
	"grader/pkg/server/session"
	solutionDelivery "grader/pkg/server/solution/delivery"
	solutionRepository "grader/pkg/server/solution/repo"
	solutionService "grader/pkg/server/solution/service"
	taskDelivery "grader/pkg/server/task/delivery"
	taskRepository "grader/pkg/server/task/repo"
	taskService "grader/pkg/server/task/service"
	userDelivery "grader/pkg/server/user/delivery"
	userRepository "grader/pkg/server/user/repo"
	userService "grader/pkg/server/user/service"
	"grader/pkg/utils"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	rabbitConn *amqp.Connection
	rabbitChan *amqp.Channel
	pgxDSN     = "postgresql://al1:A199625a11qazxc@localhost:5432/grader?sslmode=disable"
)

type config struct {
	JwtSecret string
}

func getPostgres() *sql.DB {
	db, err := sql.Open("postgres", pgxDSN)
	if err != nil {
		log.Fatalln("can't parse config", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
  			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			password VARCHAR(255) NOT NULL,
		    admin BOOLEAN NOT NULL DEFAULT false
		);
	`)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			admins INTEGER[],
			created_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS solutions (
			id SERIAL PRIMARY KEY,
			user_data JSONB NOT NULL,
			task_id INTEGER NOT NULL,
			file JSONB,
			result JSONB,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			FOREIGN KEY (task_id) REFERENCES tasks (id)
		);
	`)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("users, tasks, and solutions tables created")

	return db
}

func getRedisClient() *redis.Client {
	clientAddr := "localhost:6379"
	client := redis.NewClient(&redis.Options{
		Addr:     clientAddr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("redis connected")
	log.Println(clientAddr)

	return client
}

func main() {
	flag.Parse()
	var err error

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	cfg := config{}
	cfg.JwtSecret = os.Getenv("JWT_SECRET")

	rabbitConn, err = amqp.Dial(*queue.RabbitAddr)
	utils.FatalOnError("cant connect to rabbit", err)

	rabbitChan, err = rabbitConn.Channel()
	utils.FatalOnError("cant open chan", err)
	defer rabbitChan.Close()

	_, err = rabbitChan.QueueDeclare(
		queue.SolutionQueueName, // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	utils.FatalOnError("cant init queue", err)

	_, err = rabbitChan.QueueDeclare(
		queue.ResultQueueName, // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	utils.FatalOnError("cant declare result queue", err)

	port := 3000
	addr := ":3000"
	hostname, _ := os.Hostname()
	r := chi.NewRouter()

	jwt := cfg.JwtSecret

	templates := template.Must(template.ParseGlob("../../templates/*"))
	pgxDB := getPostgres()
	redisClient := getRedisClient()

	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()

	zapLogger = zapLogger.With(
		zap.String("hostname", hostname),
		zap.String("build", "a02ff0d0"),
	)

	zapLogger.Info("starting server",
		zap.String("logger", "ZAP"),
		zap.String("addr", addr),
		zap.Int("port", port),
	)

	l := &loggerModel.Logger{Zap: zapLogger, Level: 1}

	deflog := zapLogger.With(
		zap.String("logger", "defaultLogger"),
	).WithOptions(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.DebugLevel),
	).Sugar()

	utils.Init(deflog)

	sessionJWT := session.NewSessionJWT(jwt, redisClient)

	usersRepoPQ := userRepository.NewPgxRepo(pgxDB)
	userService := userService.NewUserService(usersRepoPQ, sessionJWT)
	userHandler := &userDelivery.UserHandler{
		Tmpl:        templates,
		UserService: userService,
	}

	solutionRepoPQ := solutionRepository.NewPgxRepo(pgxDB)
	solutionService := solutionService.NewSolutionService(solutionRepoPQ)
	solutionHandler := &solutionDelivery.SolutionHandler{
		SolutionService: solutionService,
		RabbitChan:      rabbitChan,
	}

	tasksRepoPQ := taskRepository.NewPgxRepo(pgxDB)
	taskService := taskService.NewTaskService(tasksRepoPQ)
	taskHandler := &taskDelivery.TaskHandler{
		Tmpl:            templates,
		TaskService:     taskService,
		SolutionService: solutionService,
		UserService:     userService,
	}

	//====== Pages
	//Auth
	r.Get("/login", userHandler.Login)
	r.Get("/signup", userHandler.SignUp)

	//User
	r.Get("/tasks", userHandler.Tasks)
	r.Get("/tasks/{id}", taskHandler.TaskByID)
	r.Get("/tasks/{id}/solutions/{solutionID}", taskHandler.TaskByID)
	r.Get("/tasks/user/{user}", taskHandler.TasksByUser)

	//Admin
	r.Get("/tasks/admin/task/all", taskHandler.TaskList)
	r.Get("/tasks/admin/task/create", taskHandler.TaskCreate)
	r.Get("/tasks/admin/task/{id}/edit", taskHandler.TaskEdit)
	r.Get("/tasks/admin/task/{id}/solutions", taskHandler.TaskSolutions)
	//======

	//====== API
	r.Post("/api/v1/user/register", userHandler.Register)
	r.Post("/api/v1/user/login", userHandler.Auth)
	r.Post("/api/v1/user/logout", userHandler.Logout)
	r.Post("/api/v1/solution/upload", solutionHandler.UploadSolution)
	r.Post("/api/v1/task/create", taskHandler.TaskAdd)
	r.Post("/api/v1/task/update", taskHandler.TaskUpdate)
	//======

	//Webhook
	r.Post("/webhook/solution/result", solutionHandler.SolutionResult)
	//======

	auth := middleware.Auth(sessionJWT, r)
	siteMux := middleware.AccessLog(auth)
	siteMux = middleware.Logger(l, siteMux)
	siteMux = middleware.ReqID(siteMux)
	siteMux = middleware.Panic(siteMux)

	log.Printf("server start localhost%s", addr)

	log.Fatal(http.ListenAndServe(addr, siteMux))
}
