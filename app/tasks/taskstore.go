package tasks

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TaskStore struct {
	Collection *mongo.Collection
}

func toBson(task *SyncTask) map[string]interface{} {
	res := make(map[string]interface{})

	res["key"] = task.Key
	res["type"] = task.Type
	res["params"] = task.Params
	res["lastRun"] = task.LastRun
	res["errorMsg"] = task.ErrorMsg
	res["state"] = task.State
	res["enabled"] = task.Enabled

	return res
}

func fromBson(model bson.M) *SyncTask {
	res := SyncTask{}

	res.Id = model["_id"].(primitive.ObjectID).Hex()

	if val, e := model["key"]; e {
		res.Key = val.(string)
	}

	if val, e := model["type"]; e {
		res.Type = val.(string)
	}

	if val, e := model["params"]; e && val != nil {

		res.Params = make(map[string]string)
		for k, v := range val.(primitive.M) {
			res.Params[k] = v.(string)
		}
	}

	if val, e := model["lastRun"]; e && val != nil {
		res.LastRun = val.(primitive.DateTime).Time()
	}

	if val, e := model["errorMsg"]; e {
		res.ErrorMsg = val.(string)
	}

	if val, e := model["state"]; e {
		res.State = val.(string)
	}

	if val, e := model["enabled"]; e {
		res.Enabled = val.(bool)
	}

	return &res
}

func (ts *TaskStore) InsertTask(task SyncTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := ts.Collection.InsertOne(ctx, toBson(&task))

	return err
}

func (ts *TaskStore) SaveTask(task SyncTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(task.Id)

	if err != nil {
		return fmt.Errorf("Invalid task id '%v'", task.Id)
	}

	res, err := ts.Collection.UpdateByID(ctx, id, bson.M{"$set": toBson(&task)})

	if res.MatchedCount != 1 {
		return fmt.Errorf("Can't update task with id '%v'", task.Id)
	}

	if err != nil {
		return fmt.Errorf("error on update %v: %v", task, err)
	}

	return nil
}

func (ts *TaskStore) GetTasks() ([]SyncTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := ts.Collection.Find(ctx, bson.D{{Key: "enabled", Value: true}})

	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var models []bson.M
	err = cursor.All(ctx, &models)

	if err != nil {
		return nil, err
	}

	results := make([]SyncTask, len(models))

	for i, model := range models {
		results[i] = *fromBson(model)
	}

	return results, nil
}
