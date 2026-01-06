package tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type task struct {
	ID			int 		`json:"id"`
	Description	string		`json:"description"`
	Status		string		`json:"status"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type tasks struct {
	NextID	int				`json:"next_id"`
	Db		map[int]*task	`json:"tasks"`
}

func (t task) String() string {
	createdAt := t.CreatedAt.Format("Jan Mon 10:00")
	updatedAt := t.UpdatedAt.Format("Jan Mon 10:00")

	return fmt.Sprintf("ID: %d\tDescription: %s\tStatus: %s\tCreated at: %s\tUpdated at: %s",
		 t.ID, t.Description, t.Status, createdAt, updatedAt)
}

func (t *tasks) initializeTasks() error {
	file, err := os.OpenFile("db.json", os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	t.Db = map[int]*task{}

	err = json.NewDecoder(file).Decode(t)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// encode() reopens the same file but with a trunc flag,
// this makes sure at each write to the file a new state of json
// file is encoded 
func (t *tasks) encode() error {
	file, err := os.OpenFile("db.json", os.O_TRUNC | os.O_CREATE | os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if err := json.NewEncoder(file).Encode(&t); err != nil {
		return err
	}
	return nil
}

// Parses the user input and routes the commands to their appropriate handler
func Track() {
	t := tasks{}
	if err := t.initializeTasks(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "add":
		t.add()
	case "update":
		t.update()
	case "delete":
		t.delete()
	case "mark-in-progress", "mark-done":
		t.updateStatus()
	case "list":
		t.list()
	default:
		fmt.Println("Unrecognized program argument", os.Args[1])
	}
}

// add() Creates the json file if it does not already exist and
// adds a task to it.
func (t *tasks) add() {
	if len(os.Args) != 3 {
		fmt.Println("Error: add needs a description, no less no more.")
		return
	}

	t.NextID++
	id := t.NextID
	newTask := task{
		ID:				id,
		Description:	os.Args[2],
		Status:			"todo",
		CreatedAt:		time.Now(),
		UpdatedAt:		time.Time{},
	}

	// Add task to tasks
	t.Db[id] = &newTask

	// Encode tasks and write to json file (db)
	if err := t.encode(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	fmt.Printf("Task added, and its identity: %d\n", id)
}

func (t *tasks) update() {
	if len(os.Args) != 4 {
		fmt.Println("Error: update requires task ID and description")
		return
	}

	// Parse task id
	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// Fetch a task
	task, ok := t.Db[id]
	if !ok {
		fmt.Fprintf(os.Stderr, "Task with id: %d does not exist\n", id)
		return
	}

	// Udpate the task
	task.Description = os.Args[3]
	task.UpdatedAt = time.Now()
	t.Db[id] = task

	//Encode tasks to json and write to json file (db)
	if err := t.encode(); err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Printf("Task %d updated\n", id)
}

func (t *tasks) delete() {
	if len(os.Args) != 3 {
		fmt.Println("Error: delete requires task ID")
		return
	}

	// Parse task id
	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// Delete a task
	delete(t.Db, id)

	//Encode tasks to json and write to json file (db)
	if err := t.encode(); err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Printf("Task %d deleted\n", id)
}

// Update the status of the tasks
func (t *tasks) updateStatus() {
	if len(os.Args) != 3 {
		fmt.Println("Error: updating task status requires the task's ID")
		return
	}

	// Parse task id
	id, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	// Fetch a task
	task, ok := t.Db[id]
	if !ok {
		fmt.Fprintf(os.Stderr, "Task with id: %d does not exist\n", id)
		return
	}

	// Udpate the status of the task
	if strings.ToLower(os.Args[1]) == "mark-done" {
		task.Status = "done"
	} else {
		task.Status = "in-progress"
	}

	task.UpdatedAt = time.Now()
	t.Db[id] = task

	//Encode tasks to json and write to json file (db)
	if err := t.encode(); err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Printf("Task %d status is updated\n", id)
}

// list() lists the tasks in ascending order of tasks id
// if a status is specified, it prints only tasks that match the id.
func (t *tasks) list() {
	if len(os.Args) > 3 {
		fmt.Println("Error: list takes one or no argument")
		return
	}

	keys := make([]int, len(t.Db))

	// Get the keys and sort them, to enable us list the tasks orderly.
	i := 0
	for k, _ := range t.Db {
		keys[i] = k
		i++
	}

	sort.Ints(keys)

	// List the tasks
	if len(os.Args) > 2 {
		for _, k := range keys {
			if os.Args[2] == t.Db[k].Status {
				fmt.Println(t.Db[k].String())
			}
		}
		return
	}

	for _, k := range keys {
		fmt.Printf("%s\n", t.Db[k].String())
	}
}