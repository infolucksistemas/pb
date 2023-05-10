package pb

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/infolucksistemas/pb"
	"google.golang.org/grpc"
)

const address = "localhost:9000"

// Crie um novo middleware de validação de token JWT
func JWTMiddleware() fiber.Handler {

	// Retorne o middleware de validação de token JWT
	return func(c *fiber.Ctx) error {

		// Crie uma conexão gRPC para o serviço de validação de token JWT
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// Crie um cliente gRPC para o serviço de validação de token JWT
		client := pb.NewTokenServiceClient(conn)

		// Obtenha o token JWT da solicitação
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization token not found",
			})
		}

		// Chame o serviço de validação de token JWT para verificar se o token é válido
		resp, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{Token: token})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": fmt.Sprintf("%v", err),
			})
		}

		// Se o token não for válido, retorne uma resposta de erro
		if !resp.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid authorization token",
			})
		} else {
			c.Locals("db", resp.Dados)
		}

		// Se o token for válido, continue com o manipulador da rota
		return c.Next()
	}

}
