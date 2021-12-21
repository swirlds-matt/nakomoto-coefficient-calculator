package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

var conn *pgx.Conn

type JsonResponse struct {
	Chain_name       string `json:"chain_name"`
	Chain_token      string `json:"chain_token"`
	Naka_co_curr_val int    `json:"naka_co_curr_val"`
	Naka_co_prev_val int    `json:"naka_co_prev_val"`
	Change           int    `json:"naka_co_change_val"`
}

func main() {
	var err error
	// Note: Uncomment the following line & set the database_url secretly, please.
	os.Setenv("DATABASE_URL", "postgres://xenowits:xenowits@localhost:5432/postgres")
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	r := gin.Default()
	r.GET("/nakamoto-coefficients", func(c *gin.Context) {
		coefficients := getListOfCoefficients()
		fmt.Println("fdfdfd", coefficients)
		c.JSON(200, gin.H{
			"coefficients": coefficients,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getListOfCoefficients() []JsonResponse {
	var naka_coefficients []JsonResponse
	var rows pgx.Rows
	var err error

	queryStmt := `SELECT chain_name, chain_token, naka_co_prev_val, naka_co_curr_val from naka_coefficients`
	if rows, err = conn.Query(context.Background(), queryStmt); err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var chain_name, chain_token string
		var nc_prev_val, nc_curr_val int
		err = rows.Scan(&chain_name, &chain_token, &nc_prev_val, &nc_curr_val)
		fmt.Println("values", chain_name, chain_token, nc_prev_val, nc_curr_val)
		if err != nil {
			log.Fatalln(err)
		}
		naka_coefficients = append(naka_coefficients, JsonResponse{chain_name, chain_token, nc_curr_val, nc_prev_val, nc_curr_val - nc_prev_val})
	}
	fmt.Println("ncc", naka_coefficients)
	return naka_coefficients
}
