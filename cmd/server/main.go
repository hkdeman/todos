package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/stackus/todos/internal/assets"
	"github.com/stackus/todos/internal/domain"
	"github.com/stackus/todos/internal/features/home"
	"github.com/stackus/todos/internal/features/todos"
)

type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Environment     string
}

func main() {
	// Load configuration
	cfg := loadConfig()

	// Initialize logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create router with middleware
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(cfg.ReadTimeout))

	// CORS configuration
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsMiddleware.Handler)

	// Initialize domain
	list := domain.NewTodos()

	// Add some sample todos if in development
	if cfg.Environment == "development" {
		addSampleTodos(list)
	}

	// Initialize services
	todoService := todos.NewService(list)
	homeService := home.NewService(list)

	// Mount routes
	home.Mount(router, home.NewHandler(homeService))
	todos.Mount(router, todos.NewHandler(todoService))
	assets.Mount(router)

	// Create server
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, cfg.ShutdownTimeout)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	logger.Printf("Server is running on http://localhost%s", cfg.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func loadConfig() Config {
	var cfg Config

	flag.StringVar(&cfg.Port, "port", ":3000", "port to listen on")
	flag.StringVar(&cfg.Environment, "env", "development", "environment (development/production)")
	flag.DurationVar(&cfg.ReadTimeout, "read-timeout", 30*time.Second, "read timeout")
	flag.DurationVar(&cfg.WriteTimeout, "write-timeout", 30*time.Second, "write timeout")
	flag.DurationVar(&cfg.ShutdownTimeout, "shutdown-timeout", 30*time.Second, "shutdown timeout")
	flag.Parse()

	return cfg
}

func addSampleTodos(list *domain.Todos) {
	// Add some sample todos with the new features
	todo1 := list.Add("Bake a cake")
	todo1.DueDate = ptr(time.Now().Add(24 * time.Hour))
	todo1.Priority = domain.PriorityHigh
	todo1.Category = "Cooking"
	todo1.Tags = []string{"baking", "dessert"}

	todo2 := list.Add("Feed the cat")
	todo2.DueDate = ptr(time.Now().Add(12 * time.Hour))
	todo2.Priority = domain.PriorityMedium
	todo2.Category = "Pets"
	todo2.Tags = []string{"pet care", "daily"}

	todo3 := list.Add("Take out the trash")
	todo3.DueDate = ptr(time.Now().Add(6 * time.Hour))
	todo3.Priority = domain.PriorityLow
	todo3.Category = "Household"
	todo3.Tags = []string{"chores", "daily"}

	// Add a recurring todo
	todo4 := list.Add("Weekly team meeting")
	todo4.SetRecurring("weekly", ptr(time.Now().AddDate(0, 1, 0)))
	todo4.Category = "Work"
	todo4.Tags = []string{"meeting", "team"}

	// Add a todo with subtasks
	todo5 := list.Add("Plan vacation")
	todo5.Category = "Personal"
	todo5.Tags = []string{"travel", "planning"}

	subtask1 := list.Add("Book flights")
	subtask2 := list.Add("Reserve hotel")
	subtask3 := list.Add("Create itinerary")

	todo5.AddSubtask(subtask1)
	todo5.AddSubtask(subtask2)
	todo5.AddSubtask(subtask3)
}

// Helper function to create a pointer to a time.Time
func ptr(t time.Time) *time.Time {
	return &t
}
