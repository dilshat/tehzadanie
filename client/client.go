package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/dilshat/bank/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("SERVER_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewBankServerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Add clients
	client1, err := c.AddClient(ctx, &pb.AddClientRequest{Fio: "Ivanov Ivan Ivanovich", Phone: "996555555555"})
	if err != nil {
		log.Fatalf("could not add client: %v", err)
	}
	log.Printf("Added client with id: %d", client1.ClientId)

	client2, err := c.AddClient(ctx, &pb.AddClientRequest{Fio: "Sidorov Sidor Sidorovich", Phone: "996123456789"})
	if err != nil {
		log.Fatalf("could not add client: %v", err)
	}
	log.Printf("Added client with id: %d", client2.ClientId)

	// Deposit some money
	reply, err := c.DepositMoney(ctx, &pb.DepositMoneyRequest{
		ClientId: client1.ClientId,
		Amount:   100,
	})
	if err != nil {
		log.Fatalf("could not deposit client balance: %v", err)
	}
	log.Printf("Deposited %d to client %d, new balance is %d", 100, client1.ClientId, reply.Balance)

	reply, err = c.DepositMoney(ctx, &pb.DepositMoneyRequest{
		ClientId: client2.ClientId,
		Amount:   200,
	})
	if err != nil {
		log.Fatalf("could not deposit client balance: %v", err)
	}
	log.Printf("Deposited %d to client %d, new balance is %d", 200, client2.ClientId, reply.Balance)

	// Withdraw some money
	reply2, err2 := c.WithdrawMoney(ctx, &pb.WithdrawMoneyRequest{
		ClientId: client1.ClientId,
		Amount:   50,
	})
	if err2 != nil {
		log.Fatalf("could not withdraw from client: %v", err2)
	}
	log.Printf("Withdrawed %d from client %d, new balance is %d", 50, client1.ClientId, reply2.Balance)

	reply2, err2 = c.WithdrawMoney(ctx, &pb.WithdrawMoneyRequest{
		ClientId: client2.ClientId,
		Amount:   250,
	})
	if err2 != nil {
		log.Fatalf("could not withdraw from client: %v", err2)
	}
	log.Printf("Withdrawed %d from client %d, new balance is %d", 250, client2.ClientId, reply2.Balance)

}
