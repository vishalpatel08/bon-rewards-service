package dbconnection

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vishalpatel08/bon-rewards-service/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	repo := &Repository{db: db}
	if err := repo.createTables(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) createTables() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);`

	billsTable := `
	CREATE TABLE IF NOT EXISTS bills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		due_date DATETIME NOT NULL,
		payment_date DATETIME,
		status TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	rewardsTable := `
	CREATE TABLE IF NOT EXISTS rewards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		description TEXT NOT NULL,
		issued_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	for _, stmt := range []string{usersTable, billsTable, rewardsTable} {
		if _, err := r.db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetBillByID(id int64) (*models.Bill, error) {
	query := "SELECT id, user_id, amount, due_date, payment_date, status FROM bills WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var bill models.Bill
	err := row.Scan(&bill.ID, &bill.UserID, &bill.Amount, &bill.DueDate, &bill.PaymentDate, &bill.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &bill, nil
}

func (r *Repository) UpdateBill(bill *models.Bill) error {
	query := "UPDATE bills SET payment_date = ?, status = ? WHERE id = ?"
	_, err := r.db.Exec(query, bill.PaymentDate, bill.Status, bill.ID)
	return err
}

func (r *Repository) GetLastPaidBillsByUser(userID int64, limit int) ([]models.Bill, error) {
	query := `
		SELECT id, user_id, amount, due_date, payment_date, status
		FROM bills
		WHERE user_id = ? AND status != 'UNPAID'
		ORDER BY payment_date DESC
		LIMIT ?`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []models.Bill
	for rows.Next() {
		var bill models.Bill
		if err := rows.Scan(&bill.ID, &bill.UserID, &bill.Amount, &bill.DueDate, &bill.PaymentDate, &bill.Status); err != nil {
			return nil, err
		}
		bills = append(bills, bill)
	}

	return bills, nil
}

func (r *Repository) CreateReward(reward *models.Reward) error {
	query := "INSERT INTO rewards (user_id, description, issued_at) VALUES (?, ?, ?)"
	res, err := r.db.Exec(query, reward.UserID, reward.Description, reward.IssuedAt)
	if err != nil {
		return err
	}
	reward.ID, _ = res.LastInsertId()
	return nil
}

func (r *Repository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (name, created_at) VALUES (?, ?)"
	res, err := r.db.Exec(query, user.Name, user.CreatedAt)
	if err != nil {
		return err
	}
	user.ID, _ = res.LastInsertId()
	return nil
}

func (r *Repository) CreateBill(bill *models.Bill) error {
	query := "INSERT INTO bills (user_id, amount, due_date, status) VALUES (?, ?, ?, ?)"
	res, err := r.db.Exec(query, bill.UserID, bill.Amount, bill.DueDate, bill.Status)
	if err != nil {
		return err
	}
	bill.ID, _ = res.LastInsertId()
	return nil
}
