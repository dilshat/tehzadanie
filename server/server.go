package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/dilshat/bank/api"
	db "github.com/dilshat/bank/db/gen"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBankServerServer
	ds   *db.Queries
	conn *pgxpool.Pool
}

func (s *server) AddClient(ctx context.Context, in *pb.AddClientRequest) (*pb.AddClientReply, error) {

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint
	dsTx := s.ds.WithTx(tx)

	client, err := dsTx.CreateClient(ctx, db.CreateClientParams{Fio: in.Fio, Phone: in.Phone})
	if err != nil {
		return nil, err
	}

	_, err = dsTx.CreateClientBalance(ctx, db.CreateClientBalanceParams{ClientID: client.ID, Balance: 0})
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &pb.AddClientReply{ClientId: int32(client.ID)}, nil
}

func (s *server) DepositMoney(ctx context.Context, in *pb.DepositMoneyRequest) (*pb.DepositMoneyReply, error) {
	newBalance, err := s.ds.Deposit(ctx, db.DepositParams{Balance: in.Amount, ClientID: int64(in.ClientId)})
	if err != nil {
		return nil, err
	}

	return &pb.DepositMoneyReply{Balance: newBalance}, nil
}

func (s *server) WithdrawMoney(ctx context.Context, in *pb.WithdrawMoneyRequest) (*pb.WithdrawMoneyReply, error) {
	rowsAffected, err := s.ds.Withdraw(ctx, db.WithdrawParams{Balance: in.Amount, ClientID: int64(in.ClientId)})
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("balance is not enough")
	}

	newBalance, err := s.ds.GetClientBalance(ctx, int64(in.ClientId))
	if err != nil {
		return nil, err
	}

	return &pb.WithdrawMoneyReply{Balance: newBalance}, nil
}

func newServer(conn *pgxpool.Pool, ds *db.Queries) *server {
	return &server{conn: conn, ds: ds}
}

func main() {

	conn, err := pgxpool.Connect(
		context.Background(), fmt.Sprintf("postgresql://%s:%s@%s/%s",
			os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB_ADDR"), os.Getenv("POSTGRES_DB")))
	if err != nil {
		log.Fatalf("failed to open connection to database: %v", err)
	}

	server := newServer(conn, db.New(conn))

	s := grpc.NewServer()
	pb.RegisterBankServerServer(s, server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
