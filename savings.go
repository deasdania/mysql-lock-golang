package main



import (



"database/sql"

    "errors"

    "fmt" 

    "strconv" 



_ "github.com/go-sql-driver/mysql"



) 



var (



    enableTableLocking bool = true



)  



func getDB()(*sql.DB, error) {



    dbUser     := "root"

    dbPassword := ""

    dbName     := "sample_db"



    db, err := sql.Open("mysql", dbUser + ":" + dbPassword + "@tcp(127.0.0.1:3306)/" + dbName) 



    if err != nil {

        return nil, err

    }



    return db, nil



}



func addEntry(p map[string]interface{}) (map[string]interface{}, error){



    accountId, err := strconv.ParseInt(fmt.Sprint(p["account_id"]), 10, 64)



    if err != nil {

        return nil, err

    }



    transType := p["trans_type"].(string)

    amount    := p["amount"].(float64)



    credit := 0.0

    debit  := 0.0



    if transType == "D" {

        credit = amount

        debit  = 0.00

    } else {

        credit = 0.00

        debit  = amount

    }



    db, err := getDB()



    if err != nil {

        return nil, err

    }



    defer db.Close()





    if enableTableLocking == true {

        lockTables(db)

    }





    resp, err := getBalance(db, accountId) 



    accountBalance := resp["account_balance"].(float64)          



    if amount > accountBalance && transType == "W" {



        if enableTableLocking  == true {

            unlockTables(db)

        }     



       return nil, errors.New("Insufficient balance. " + fmt.Sprint(accountBalance))



    }



    queryString := "insert into savings (account_id, trans_type, debit, credit) values (?, ?, ?, ?)"



    stmt, err   := db.Prepare(queryString) 



    if err != nil {

        return nil, err       

    }



    defer stmt.Close()     



    res, err := stmt.Exec(accountId, transType, debit, credit)  



    if err != nil {

        return nil, err

    }



    refId, err := res.LastInsertId()



    if err != nil {

        return nil, err

    }



    resp, err = getBalance(db, accountId) 



    accountBalance = resp["account_balance"].(float64) 





    if enableTableLocking {

        unlockTables(db)

    }





    response := map[string]interface{}{

                    "ref_id" :    refId,

                    "account_id": accountId,

                    "amount":     amount,

                    "balance":    accountBalance,

                } 



    return response, nil  

}



func getBalance(db *sql.DB, accountId int64) (map[string]interface{}, error) {    



    queryString := "select ifnull(sum(credit - debit), 0) as account_balance from savings where account_id = ?"



    stmt, err := db.Prepare(queryString) 



    if err != nil {

        return nil, err       

    }



    accountBalance := 0.00



    err = stmt.QueryRow(accountId).Scan(&accountBalance)



    if err != nil {

       return nil, err

    }



    response := map[string]interface{}{

                    "account_balance" : accountBalance,

                }



    return response, nil

 }



func lockTables(db *sql.DB) error {



    queryString := "lock tables savings write"



    _, err := db.Exec(queryString) 



    if err != nil {

        return err       

    }



    return nil

}



func unlockTables(db *sql.DB) error {



    queryString := "unlock tables"



    _, err := db.Exec(queryString) 



    if err != nil {

        return err       

    }



    return nil

}
