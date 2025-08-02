package svc

import (
	"errors"
	"log"

	"github.com/amrishshah/distributed2pc/io"
)

type Agent struct {
	Id         string  `json:"id"`
	IsReserved int     `json:"is_reserved"`
	OrderId    *string `json:"order_id"`
}

func BookAgent(OrderId string) (*Agent, error) {
	txn, _ := io.DB.Begin()

	row := txn.QueryRow(`select id, is_reserved, order_id 
	from agents
	where is_reserved is true 
	and order_id is null 
	limit 1
	FOR UPDATE`)
	if row.Err() != nil {
		txn.Rollback()
		return nil, row.Err()
	}
	var agent Agent
	err := row.Scan(&agent.Id, &agent.IsReserved, &agent.OrderId)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("no delivery agent available")
	}

	_, err = txn.Exec("UPDATE agents SET is_reserved = false, order_id = ? where id = ?", OrderId, agent.Id)

	if err != nil {
		txn.Rollback()
		return nil, err
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}
	agent.OrderId = &OrderId
	agent.IsReserved = 0
	return &agent, nil

}

func ReverseAgent() (*Agent, error) {
	log.Print("sd")
	txn, _ := io.DB.Begin()
	log.Print("sd 1")
	row := txn.QueryRow(`select id, is_reserved, order_id 
	from agents
	where is_reserved = 0 
	and order_id is null 
	limit 1
	FOR UPDATE`)
	if row.Err() != nil {
		log.Print("sd 2")
		txn.Rollback()
		return nil, row.Err()
	}
	var agent Agent
	err := row.Scan(&agent.Id, &agent.IsReserved, &agent.OrderId)

	if err != nil {
		log.Print("sd 3")
		txn.Rollback()
		return nil, errors.New("no delivery agent available")
	}

	_, err = txn.Exec("UPDATE agents SET is_reserved = true where id = ?", agent.Id)

	if err != nil {
		log.Print("sd 4")
		txn.Rollback()
		return nil, err
	}

	err = txn.Commit()
	if err != nil {
		log.Print("sd 5")
		return nil, err
	}
	return &agent, nil
}
