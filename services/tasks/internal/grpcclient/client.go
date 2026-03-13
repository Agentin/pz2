package grpcclient

import (
	"context"

	authv1 "github.com/student/tech-ip-sem2/pkg/api/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthClient struct {
	client authv1.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(addr string) (*AuthClient, error) {
	// Устанавливаем соединение (без таймаута — соединение может быть долгим)
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) // для учебных целей без TLS
	if err != nil {
		return nil, err
	}
	return &AuthClient{
		client: authv1.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// Verify проверяет токен с заданным контекстом (с deadline).
func (c *AuthClient) Verify(ctx context.Context, token string) (bool, string, error) {
	resp, err := c.client.Verify(ctx, &authv1.VerifyRequest{Token: token})
	if err != nil {
		// Преобразуем gRPC ошибку
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return false, "", nil // токен невалиден, это не ошибка клиента, а статус
			case codes.DeadlineExceeded:
				return false, "", err // таймаут
			default:
				return false, "", err // другие ошибки (недоступность и т.п.)
			}
		}
		return false, "", err
	}
	return resp.Valid, resp.Subject, nil
}
