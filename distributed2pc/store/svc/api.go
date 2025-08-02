package svc

import (
	"errors"
	"log"

	"github.com/amrishshah/distributed2pc/io"
)

type Packet struct {
	Id         string  `json:"id"`
	FoodId     int     `json:"food_id"`
	IsReserved int     `json:"is_reserved"`
	OrderId    *string `json:"order_id"`
}

func BookPacket(OrderId string, foodID int) (*Packet, error) {
	txn, _ := io.DB.Begin()

	row := txn.QueryRow(`select id, food_id,is_reserved, order_id 
	from packets
	where is_reserved is true 
	and order_id is null and food_id = ?
	limit 1
	FOR UPDATE`, foodID)
	if row.Err() != nil {
		txn.Rollback()
		return nil, row.Err()
	}
	var packet Packet
	err := row.Scan(&packet.Id, &packet.FoodId, &packet.IsReserved, &packet.OrderId)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("no delivery agent available")
	}

	_, err = txn.Exec("UPDATE packets SET is_reserved = false, order_id = ? where id = ?", OrderId, packet.Id)

	if err != nil {
		txn.Rollback()
		return nil, err
	}

	err = txn.Commit()
	if err != nil {
		return nil, err
	}
	packet.OrderId = &OrderId
	packet.IsReserved = 0
	return &packet, nil
}

func ReversePacket(foodID int) (*Packet, error) {
	log.Print("sd")
	txn, _ := io.DB.Begin()
	log.Print("sd 1")
	log.Println(foodID)
	row := txn.QueryRow(`select id, food_id,is_reserved, order_id 
	from packets
	where is_reserved = 0 
	and order_id is null
	and food_id = ?
	limit 1
	FOR UPDATE`, foodID)
	if row.Err() != nil {
		log.Print("sd 2")
		txn.Rollback()
		return nil, row.Err()
	}
	var packet Packet
	err := row.Scan(&packet.Id, &packet.FoodId, &packet.IsReserved, &packet.OrderId)

	if err != nil {
		log.Print("sd 3")
		log.Print(err.Error())
		txn.Rollback()
		return nil, errors.New("no food packet available")
	}

	_, err = txn.Exec("UPDATE packets SET is_reserved = true where id = ?", packet.Id)

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
	return &packet, nil
}
