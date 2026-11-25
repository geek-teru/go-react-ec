//go:build e2e

package helper

// テーブルができたら、ここにデータ投入のファクトリを足す。形の例:
//
// type UserOpt func(*userRow)
//
// func WithEmail(s string) UserOpt { return func(u *userRow) { u.Email = s } }
//
// func CreateUser(t *testing.T, db *sql.DB, opts ...UserOpt) int64 {
// 	t.Helper()
// 	u := userRow{Email: fmt.Sprintf("user-%d@example.test", time.Now().UnixNano()), Name: "テスト"}
// 	for _, o := range opts { o(&u) }
// 	var id int64
// 	if err := db.QueryRowContext(context.Background(),
// 		`INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id`, u.Email, u.Name).Scan(&id); err != nil {
// 		t.Fatalf("create user: %v", err)
// 	}
// 	return id
// }
