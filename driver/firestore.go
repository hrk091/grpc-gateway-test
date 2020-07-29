package driver

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"log"
	"reflect"
	"strings"

	"grpc-gateway-test/env"
)

type FirestoreDriver struct {
	Doc1stKey string
	Col2ndKey string
}

func (d *FirestoreDriver) getDocRef(client *firestore.Client, category, key string) *firestore.DocumentRef {
	if d.Doc1stKey != "" && d.Col2ndKey != "" {
		return client.Collection(category).Doc(d.Doc1stKey).Collection(d.Col2ndKey).Doc(key)
	} else {
		return client.Collection(category).Doc(key)
	}
}

func (d *FirestoreDriver) getColRef(client *firestore.Client, category string) *firestore.CollectionRef {
	if d.Doc1stKey != "" && d.Col2ndKey != "" {
		return client.Collection(category).Doc(d.Doc1stKey).Collection(d.Col2ndKey)
	} else {
		return client.Collection(category)
	}
}

func (d *FirestoreDriver) getDocIter(client *firestore.Client, ctx context.Context, category string, name string, ) *firestore.DocumentIterator {
	var col *firestore.CollectionRef
	if d.Doc1stKey != "" && d.Col2ndKey != "" {
		col = client.Collection(category).Doc(d.Doc1stKey).Collection(d.Col2ndKey)
	} else {
		col = client.Collection(category)
		//return client.Collection(category).Documents(ctx)
	}
	if name != "" {
		return col.Where("name", "==", name).Documents(ctx)
	} else {
		return col.Documents(ctx)
	}

}

func (d *FirestoreDriver) Get(ctx context.Context, category, key string) (map[string]interface{}, error) {
	client := GetFireStoreClient(ctx)
	defer CloseFireStoreClient(client)

	docRef := d.getDocRef(client, category, key)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	data := doc.Data()
	data["Id"] = docRef.ID
	return data, nil
}

func (d *FirestoreDriver) GetAll(ctx context.Context, category string, queryName string) ([]map[string]interface{}, error) {
	client := GetFireStoreClient(ctx)
	defer CloseFireStoreClient(client)
	iter := d.getDocIter(client, ctx, category, queryName)
	rep := make([]map[string]interface{}, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		data["Id"] = doc.Ref.ID
		rep = append(rep, data)
	}
	return rep, nil
}

func (d *FirestoreDriver) Create(ctx context.Context, category, key string, doc interface{}) error {
	client := GetFireStoreClient(ctx)
	defer CloseFireStoreClient(client)
	varMap := StructToMap(doc)
	delete(varMap, "id")
	delete(varMap, "category")

	if key != "" {
		docRef := d.getDocRef(client, category, key)
		curr, _ := docRef.Get(ctx)
		if curr != nil {
			err := errors.New("Error: Doc already exists: " + key)
			return err
		}

		_, err := docRef.Set(ctx, varMap)
		if err != nil {
			return err
		}
	} else {
		colRef := d.getColRef(client, category)
		_, _, err := colRef.Add(ctx, varMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *FirestoreDriver) Update(ctx context.Context, category, key string, doc interface{}) (map[string]interface{}, error) {
	client := GetFireStoreClient(ctx)
	defer CloseFireStoreClient(client)

	varMap := StructToMap(doc)
	delete(varMap, "id")
	delete(varMap, "category")

	docRef := d.getDocRef(client, category, key)
	_, err := docRef.Set(ctx, varMap, firestore.MergeAll)
	if err != nil {
		return nil, err
	}

	// get updated value
	res, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	data := res.Data()
	data["id"] = res.Ref.ID
	return data, nil
}

func (d *FirestoreDriver) Delete(ctx context.Context, category, key string) error {
	client := GetFireStoreClient(ctx)
	defer CloseFireStoreClient(client)

	docRef := d.getDocRef(client, category, key)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetFireStoreClient(ctx context.Context) *firestore.Client {

	client, err := firestore.NewClient(ctx, env.ProjectID)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}
	return client
}

func CloseFireStoreClient(client *firestore.Client) {
	err := client.Close()
	if err != nil {
		log.Panicf("Failed to close client: %v", err)
	}
}

func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()

	for i := 0; i < elem.NumField(); i++ {
		name := elem.Type().Field(i).Name
		field := strings.ToLower(name[0:1]) + name[1:]
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}
