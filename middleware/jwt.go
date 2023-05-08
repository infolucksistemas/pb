package middleware

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"

	"github.com/infolucksistemas/pb"
)

const address = "localhost:50051"

// Crie um novo middleware de validação de token JWT
func JWTMiddleware() fiber.Handler {

	// Retorne o middleware de validação de token JWT
	return func(ctx *fiber.Ctx) error {

		// Crie uma conexão gRPC para o serviço de validação de token JWT
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// Crie um cliente gRPC para o serviço de validação de token JWT
		client := pb.NewTokenServiceClient(conn)

		// Obtenha o token JWT da solicitação
		token := ctx.Get("Authorization")
		if token == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization token not found",
			})
		}

		// Chame o serviço de validação de token JWT para verificar se o token é válido
		resp, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("%v", err),
			})

		}

		// Se o token não for válido, retorne uma resposta de erro
		if !resp.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid authorization token",
			})
		}

		// Se o token for válido, continue com o manipulador da rota
		return ctx.Next()
	}

}
