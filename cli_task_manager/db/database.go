package database

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")

var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

func Init(db_path string) error {
	var err error
	db, err = bolt.Open(db_path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	calllback_fn := func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	}

	return db.Update(calllback_fn)
}

func CreateTask(task string) (int, error) {

	var id int

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		id_u64, _ := bucket.NextSequence()
		id = int(id_u64)
		key := intToByte(id)
		return bucket.Put(key, []byte(task))
	})

	if err != nil {
		return -1, err
	}

	return id, nil
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			tasks = append(tasks, Task{
				Key:   byteToInt(key),
				Value: string(value),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		return bucket.Delete(intToByte(key))
	})
}

func intToByte(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func byteToInt(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
