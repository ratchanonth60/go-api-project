package main

import (
	"flag"
	"log"
	"project-api/internal/infra/config"
	"project-api/internal/task"
)

func main() {

	// Parse flags
	configType := flag.String("config", "yaml", "Configuration type (env or yaml)")
	flag.Parse()

	// โหลด configuration
	var err error
	if *configType == "env" {
		config.IsYaml = false
		err = config.LoadConfig("")
	} else {
		err = config.LoadConfig("conf/app.yaml")
	}
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	server, err := task.NewMachineryServer()
	if err != nil {
		log.Fatalf("Failed to start Machinery server: %v", err)
	}
	err = server.RegisterTasks(map[string]interface{}{
		"send_confirmation_email": func(toEmail, token, name string, host string) error {
			return task.TaskSendConfirmationEmail(toEmail, token, name, host)
		},
	})
	if err != nil {
		log.Fatalf("Failed to register tasks: %v", err)
	}

	// เริ่ม worker
	worker := server.NewWorker("email_worker", 10) // 10 concurrent workers
	if err := worker.Launch(); err != nil {
		log.Fatalf("Failed to launch worker: %v", err)
	}
}
