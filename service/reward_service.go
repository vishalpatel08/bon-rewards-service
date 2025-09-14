package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vishalpatel08/bon-rewards-service/models"
)

const (
	requiredOnTimePayments = 3
	rewardDescription      = "$10 Amazon Gift Card"
)

type Repository interface {
	GetBillByID(id int64) (*models.Bill, error)
	UpdateBill(bill *models.Bill) error
	GetLastPaidBillsByUser(userID int64, limit int) ([]models.Bill, error)
	CreateReward(reward *models.Reward) error
	CreateUser(user *models.User) error
	CreateBill(bill *models.Bill) error
}

type RewardService struct {
	repo Repository
}

func NewRewardService(r Repository) *RewardService {
	return &RewardService{repo: r}
}

func (s *RewardService) PayBill(ctx context.Context, billID int64) (*models.Bill, string, error) {
	bill, err := s.repo.GetBillByID(billID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get bill: %w", err)
	}
	if bill == nil {
		return nil, "", errors.New("bill not found")
	}
	if bill.Status != models.StatusUnpaid {
		return nil, "", errors.New("bill has already been paid")
	}

	paymentTime := time.Now()
	bill.PaymentDate = &paymentTime

	if paymentTime.After(bill.DueDate) {
		bill.Status = models.StatusPaidLate
	} else {
		bill.Status = models.StatusPaidOnTime
	}

	if err := s.repo.UpdateBill(bill); err != nil {
		return nil, "", fmt.Errorf("failed to update bill: %w", err)
	}

	rewardMessage, err := s.checkForReward(bill.UserID)
	if err != nil {
		log.Printf("ERROR checking for reward for user %d: %v", bill.UserID, err)
	}

	return bill, rewardMessage, nil
}

func (s *RewardService) checkForReward(userID int64) (string, error) {
	lastBills, err := s.repo.GetLastPaidBillsByUser(userID, requiredOnTimePayments)
	if err != nil {
		return "", fmt.Errorf("could not get last paid bills: %w", err)
	}

	if len(lastBills) < requiredOnTimePayments {
		log.Printf("User %d has only %d paid bills, not eligible for reward yet.", userID, len(lastBills))
		return "", nil
	}

	for _, b := range lastBills {
		if b.Status != models.StatusPaidOnTime {
			log.Printf("User %d is not eligible for a reward. Bill ID %d was paid late.", userID, b.ID)
			return "", nil
		}
	}
	log.Printf("SUCCESS: User %d has earned a reward!", userID)

	newReward := &models.Reward{
		UserID:      userID,
		Description: rewardDescription,
		IssuedAt:    time.Now(),
	}

	if err := s.repo.CreateReward(newReward); err != nil {
		return "", fmt.Errorf("failed to create reward: %w", err)
	}
	return fmt.Sprintf("Congratulations! You've earned a %s.", rewardDescription), nil
}

func (s *RewardService) CreateUser(ctx context.Context, name string) (*models.User, error) {
	user := &models.User{
		Name:      name,
		CreatedAt: time.Now(),
	}
	err := s.repo.CreateUser(user)
	return user, err
}

func (s *RewardService) CreateBill(ctx context.Context, userID int64, amount int64, dueDate time.Time) (*models.Bill, error) {
	bill := &models.Bill{
		UserID:  userID,
		Amount:  amount,
		DueDate: dueDate,
		Status:  models.StatusUnpaid,
	}
	err := s.repo.CreateBill(bill)
	return bill, err
}
