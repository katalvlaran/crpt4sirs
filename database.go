package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "user=sir password=crpt-pass dbname=crpt4sirs sslmode=disable" // Use a configuration management system for this
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the PostgreSQL database.")
}

//func saveSignToDatabase(db *sql.DB, sign Signal) (int, error) {
//	//t := time.Now()
//	//stmt, err := db.Prepare(`INSERT INTO signal (position_side, leverage, chat, symbol, time)
//	//								VALUES ($1, $2, $3, $4, $5)`)
//	//if err != nil {
//	//	return err
//	//}
//	//defer stmt.Close()
//	//
//	//if sign.Symbol == "" || sign.PositionSide == "" {
//	//	return fmt.Errorf("empty symbol or position side")
//	//}
//	//
//	//_, err = stmt.Exec(sign.PositionSide, sign.Leverage, sign.Chat, sign.Symbol, t)
//	//if err != nil {
//	//	return err
//	//}
//	////SELECT currval('your_table_id_seq');
//	//
//	//fmt.Println("Message saved to the database.")
//	//return nil
//	t := time.Now()
//	stmt, err := db.Prepare(`INSERT INTO signal (position_side, leverage, chat, symbol, time)
//                            VALUES ($1, $2, $3, $4, $5) RETURNING id`)
//	if err != nil {
//		return 0, err
//	}
//	defer stmt.Close()
//
//	if sign.Symbol == "" || sign.PositionSide == "" {
//		return 0, fmt.Errorf("empty symbol or position side")
//	}
//
//	var id int
//	err = stmt.QueryRow(sign.PositionSide /*sign.Leverage,*/, sign.Chat, sign.Symbol, t).Scan(&id)
//	if err != nil {
//		return id, err
//	}
//
//	fmt.Println("Message saved to the database with ID:", id)
//	//return id
//	return id, err
//}

func saveToTable(db *sql.DB, data interface{} /*price string,*/, tableName string) (int, error) {
	t := time.Now()
	var stmt *sql.Stmt
	var err error

	switch tableName {
	case "signal":
		stmt, err = db.Prepare(`INSERT INTO signal (position_side, chat, symbol, time) 
                                VALUES ($1, $2, $3, $4) RETURNING id`)
	case "orders":
		stmt, err = db.Prepare(`INSERT INTO orders (symbol, side, otype, price, quantity, signal_id, time)
	                           VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`)
	case "testklines":
		stmt, err = db.Prepare(`INSERT INTO testklines (open_time, close_time, open_price, high_price, low_price, close_price,
								volume, quote_asset_volume, number_of_trades, taker_buy_base_asset_volume, 
                       			taker_buy_quote_asset_volume)
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`)
	default:
		return 0, fmt.Errorf("unsupported table name")
	}

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int

	switch tableName {
	case "signal":
		signal, ok := data.(Signal)
		if !ok {
			return 0, fmt.Errorf("invalid data type for signals")
		}
		if signal.Symbol == "" || signal.PositionSide == "" {
			return 0, fmt.Errorf("empty symbol or position side")
		}
		err = stmt.QueryRow(signal.PositionSide, signal.Chat, signal.Symbol, t).Scan(&id)

	case "orders":
		order, ok := data.(Order)
		if !ok {
			return 0, fmt.Errorf("invalid data type for orders")
		}

		err = stmt.QueryRow(order.Symbol, order.Side, order.OType, order.Price, order.Quantity, order.SignalId, t).Scan(&id)

	case "testklines":
		btcklines, ok := data.(Klines)
		if !ok {
			return 0, fmt.Errorf("invalid data type for btcklines")
		}

		err = stmt.QueryRow(btcklines.openTime, btcklines.closeTime, btcklines.openPrice, btcklines.highPrice,
			btcklines.lowPrice, btcklines.closePrice, btcklines.volume, btcklines.quoteAssetVolume,
			btcklines.numberOfTrades, btcklines.takerBuyBaseAssetVolume, btcklines.takerBuyQuoteAssetVolume).Scan(&id)

	}

	if err != nil {
		return id, err
	}

	fmt.Printf("Data saved to the '%s' table with ID: %d\n", tableName, id)
	return id, nil
}

func getMaxSignalId(db *sql.DB) (int, error) {
	var id int
	query := "SELECT max(id) FROM signal"
	err := db.QueryRow(query).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
