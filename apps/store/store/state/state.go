package state

import (
	"context"
	"encoding/json"

	stategrpc "github.com/CoreumFoundation/CoreDEX-API/domain/state"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
	store "github.com/CoreumFoundation/CoreDEX-API/utils/mysqlstore"
)

type Application struct {
	client store.StoreBase
}

func NewApplication(client *store.StoreBase) *Application {
	app := &Application{
		client: *client,
	}
	app.initDB()
	return app
}

// Initialize tables and indexes
func (a *Application) initDB() {
	_, err := a.client.Client.Exec(`CREATE TABLE IF NOT EXISTS State (
		StateType INT, 
		Content TEXT, 
		MetaData JSON, 
		Network INT AS (JSON_UNQUOTE(JSON_EXTRACT(MetaData, '$.Network'))),
		UNIQUE KEY (Network,StateType)
	)`)
	if err != nil {
		logger.Fatalf("Error creating State table: %v", err)
	}
}

func (a *Application) Get(ctx context.Context, in *stategrpc.StateQuery) (*stategrpc.State, error) {
	// Use the mysql client to query for the provided data in the table state:
	b, err := a.client.Client.Query("SELECT MetaData, StateType, Content FROM State WHERE Network=? AND StateType=?",
		in.Network, in.StateType)
	if err != nil {
		return nil, err
	}
	defer b.Close()
	// Map the result into the stategrpc.State struct:
	state := &stategrpc.State{}
	md := make([]byte, 0)
	// We are querying be unique key so only get a single result
	for b.Next() {
		err = b.Scan(&md, &state.StateType, &state.Content)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(md, &state.MetaData)
	}
	return state, nil
}

func (a *Application) Upsert(ctx context.Context, in *stategrpc.State) error {
	md, err := json.Marshal(in.MetaData)
	if err != nil {
		logger.Errorf("Error marshalling metadata: %v", err)
		return err
	}
	// Use the mysql client to insert the provided data into the table state:
	_, err = a.client.Client.Exec(`INSERT INTO State (StateType, Content, MetaData) VALUES (?, ?, ?) 
	 		ON DUPLICATE KEY UPDATE Content=?, MetaData=?`,
		in.StateType, in.Content, md, in.Content, md)
	if err != nil {
		logger.Errorf("Error upserting state: %v", err)
		return err
	}
	return nil
}
