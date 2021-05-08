package tasks

import (
    //"fmt"
)


type TaskManager struct {
    RunnerManager *TaskRunManager
    TaskStore *TaskStore
    AbookStore *AbookStore
}


func (tm *TaskManager) Process() error {
    tasks, err := tm.TaskStore.GetTasks()
    
    if err != nil {
        return err
    }

    for _, task := range tasks {
        err := tm.ProcessTask(&task)
        if err != nil {
            return err
        }
    }


    return nil
}

func (tm *TaskManager) ProcessTask(task *SyncTask) error {
    updated, books := tm.RunnerManager.RunTask(*task)

    err := tm.TaskStore.SaveTask(updated)
    if err != nil {
       return err 
    }


    for _, book := range books {
        //fmt.Println(book.Id, book.Title, book.Date)
        tm.AbookStore.Upsert(book)
    }

    return nil
}
