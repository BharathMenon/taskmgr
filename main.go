package main

import (
	// "flag"
	// "fmt"
	// "os"
	//api "github.com/BharathMenon/taskmgr/api"
	gin "github.com/gin-gonic/gin"
    v1 "github.com/BharathMenon/taskmgr/api/v1"
   auth "github.com/BharathMenon/taskmgr/auth"
	//task "github.com/BharathMenon/taskmgr/task"
)

	func main() {
		//api.StartServer()
		r:=gin.Default()
		//swagger stuff
		authGroup:=r.Group("/auth")
		{
			authGroup.POST("/register",auth.Register)
			authGroup.POST("/login",auth.Login)
		}
		apiV1 := r.Group("/api/v1")
		apiV1.Use(auth.AuthRequired())
		{
        apiV1.GET("/tasks", v1.ListTasks)
        apiV1.POST("/tasks", v1.NewTask)
        apiV1.GET("/tasks/:id", v1.GetTask)
        apiV1.PUT("/tasks/:id", v1.UpdateTask)
        apiV1.DELETE("/tasks/:id", v1.DeleteTask)
        apiV1.PUT("/tasks/:id/complete", v1.MarkComplete)
    }
	r.Run(":8080") 
	// 	addFlag := flag.Bool("add", false, "Add a task")
	// 	listFlag := flag.Bool("list", false, "List tasks")
	// 	updateFlag := flag.Bool("update", false, "Update a task")
	// 	deleteFlag := flag.Bool("delete", false, "Delete a task")
	// 	completeFlag := flag.Bool("complete", false, "Mark a task complete")
	
	// 	// Supporting flags
	// 	idFlag := flag.Int("id", 0, "Task ID (for update/delete/complete)")
	// 	titleFlag := flag.String("title", "", "Task title")
	// 	descFlag := flag.String("desc", "", "Task description")
	// 	statusFlag := flag.String("status", "", "Task status (pending|done)")
	
	// 	// You can add your logic here to use the flags
	// 	flag.Parse()
	// 	actions := 0

	// 	//actions := 0
	// for _, v := range []bool{*addFlag, *listFlag, *updateFlag, *deleteFlag, *completeFlag} {
	// 	if v {
	// 		actions++
	// 	}
	// }
	// if actions == 0 {
	// 	fmt.Println("No action specified. Use --help for usage.")
	// 	return
	// }
	// if actions > 1 {
	// 	fmt.Println("Only one action allowed per invocation. Use --help.")
	// 	return
	// }

	// path := task.TasksFilePath()

	// switch {
	// case *addFlag:
	// 	task, err := task.AddTask(path, *titleFlag, *descFlag)
	// 	if err != nil {
	// 		fmt.Println("Error adding task:", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println("Added task:", task.ID)

	// case *listFlag:
	// 	tasks, err := task.ListTasks(path)
	// 	if err != nil {
	// 		fmt.Println("Error listing tasks:", err)
	// 		os.Exit(1)
	// 	}
	// 	task.PrintTasks(tasks)

	// case *updateFlag:
	// 	if *idFlag <= 0 {
	// 		fmt.Println("--id is required for update")
	// 		os.Exit(1)
	// 	}
	// 	// pass nil for unchanged
	// 	var titlePtr, descPtr, statusPtr *string
	// 	if *titleFlag != "" {
	// 		titlePtr = titleFlag
	// 	}
	// 	if *descFlag != "" {
	// 		descPtr = descFlag
	// 	}
	// 	if *statusFlag != "" {
	// 		statusPtr = statusFlag
	// 	}
	// 	task, err := task.UpdateTask(path, *idFlag, titlePtr, descPtr, statusPtr)
	// 	if err != nil {
	// 		fmt.Println("Error updating task:", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println("Updated task:", task.ID)

	// case *deleteFlag:
	// 	if *idFlag <= 0 {
	// 		fmt.Println("--id is required for delete")
	// 		os.Exit(1)
	// 	}
	// 	if err := task.DeleteTask(path, *idFlag); err != nil {
	// 		fmt.Println("Error deleting task:", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println("Deleted task:", *idFlag)

	// case *completeFlag:
	// 	if *idFlag <= 0 {
	// 		fmt.Println("--id is required for complete")
	// 		os.Exit(1)
	// 	}
	// 	task, err := task.CompleteTask(path, *idFlag)
	// 	if err != nil {
	// 		fmt.Println("Error completing task:", err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println("Completed task:", task.ID)
	// }
	}