package services

import (
    "errors"
    "insurance-claims-api/internal/models"
    "insurance-claims-api/internal/repositories"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type ClaimService interface {
    CreateClaim(userID primitive.ObjectID, req models.CreateClaimRequest) (*models.Claim, error)
    GetMyClaims(userID primitive.ObjectID, page, limit int) ([]models.Claim, int64, error)
    GetAllClaims(role string, page, limit int) ([]models.Claim, int64, error)
    GetClaimByID(userID primitive.ObjectID, role string, claimID primitive.ObjectID) (*models.Claim, error)
    UpdateClaim(userID primitive.ObjectID, claimID primitive.ObjectID, req models.UpdateClaimRequest) error
    DeleteClaim(userID primitive.ObjectID, claimID primitive.ObjectID) error
    SubmitClaim(userID primitive.ObjectID, claimID primitive.ObjectID) error
    ReviewClaim(verifierID primitive.ObjectID, claimID primitive.ObjectID, note string) error
    ApproveClaim(approverID primitive.ObjectID, claimID primitive.ObjectID) error
    RejectClaim(approverID primitive.ObjectID, claimID primitive.ObjectID, reason string) error
}

type claimService struct {
    claimRepo repositories.ClaimRepository
    userRepo  repositories.UserRepository
}

func NewClaimService(claimRepo repositories.ClaimRepository, userRepo repositories.UserRepository) ClaimService {
    return &claimService{claimRepo, userRepo}
}

func (s *claimService) CreateClaim(userID primitive.ObjectID, req models.CreateClaimRequest) (*models.Claim, error) {
    claim := &models.Claim{
        ID:           primitive.NewObjectID(),
        UserID:       userID,
        PolicyNumber: req.PolicyNumber,
        ClaimAmount:  req.ClaimAmount,
        Description:  req.Description,
        Documents:    req.Documents,
        Status:       models.Draft,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
        History: []models.ClaimHistory{
            {
                Status:    models.Draft,
                ChangedBy: userID,
                ChangedAt: time.Now(),
            },
        },
    }

    err := s.claimRepo.Create(claim)
    if err != nil {
        return nil, err
    }
    return claim, nil
}

func (s *claimService) GetMyClaims(userID primitive.ObjectID, page, limit int) ([]models.Claim, int64, error) {
    return s.claimRepo.FindByUserID(userID, page, limit)
}

func (s *claimService) GetAllClaims(role string, page, limit int) ([]models.Claim, int64, error) {
    if role == "user" {
        return nil, 0, errors.New("forbidden")
    }
    filter := bson.M{}
    if role == "verifier" {
        filter = bson.M{"status": bson.M{"$in": []models.ClaimStatus{models.Submitted, models.Reviewed}}}
    }
    if role == "approver" {
        filter = bson.M{"status": bson.M{"$in": []models.ClaimStatus{models.Reviewed, models.Approved, models.Rejected}}}
    }
    return s.claimRepo.FindAll(filter, page, limit)
}

func (s *claimService) GetClaimByID(userID primitive.ObjectID, role string, claimID primitive.ObjectID) (*models.Claim, error) {
    claim, err := s.claimRepo.FindByID(claimID)
    if err != nil {
        return nil, err
    }
    if role != "verifier" && role != "approver" && claim.UserID != userID {
        return nil, errors.New("forbidden")
    }
    return claim, nil
}

func (s *claimService) UpdateClaim(userID primitive.ObjectID, claimID primitive.ObjectID, req models.UpdateClaimRequest) error {
    claim, err := s.claimRepo.FindByID(claimID)
    if err != nil {
        return err
    }
    if claim.UserID != userID || claim.Status != models.Draft {
        return errors.New("cannot edit claim")
    }
    if req.PolicyNumber != "" {
        claim.PolicyNumber = req.PolicyNumber
    }
    if req.ClaimAmount > 0 {
        claim.ClaimAmount = req.ClaimAmount
    }
    if req.Description != "" {
        claim.Description = req.Description
    }
    if req.Documents != nil {
        claim.Documents = req.Documents
    }
    return s.claimRepo.Update(claim)
}

func (s *claimService) DeleteClaim(userID primitive.ObjectID, claimID primitive.ObjectID) error {
    claim, err := s.claimRepo.FindByID(claimID)
    if err != nil {
        return err
    }
    if claim.UserID != userID || claim.Status != models.Draft {
        return errors.New("cannot delete claim")
    }
    return s.claimRepo.Delete(claimID)
}

func (s *claimService) SubmitClaim(userID primitive.ObjectID, claimID primitive.ObjectID) error {
    // claim, err := s.claimRepo.FindByID(claimID)
    // if err != nil || claim.UserID != userID || claim.Status != models.Draft {
    //     return errors.New("invalid operation")
    // }
    // claim.Status = models.Submitted
    // history := models.ClaimHistory{
    //     Status:    models.Submitted,
    //     ChangedBy: userID,
    //     ChangedAt: time.Now(),
    // }
    // s.claimRepo.AddHistory(claimID, history)
    // return s.claimRepo.Update(claim)

    claim, err := s.claimRepo.FindByID(claimID)
    if err != nil || claim.Status != models.Draft {
        return errors.New("invalid operation")
    }

    now := time.Now()
    history := models.ClaimHistory{
        Status:    models.Submitted,
        ChangedBy: userID,
        ChangedAt: now,
    }

    update := bson.M{
        "$set": bson.M{
            "status":     models.Submitted,
            "updated_at": now,
        },
        "$push": bson.M{"history": history},
    }

    return s.claimRepo.UpdateWithPush(claimID, update)
}

// Verifier only
func (s *claimService) ReviewClaim(verifierID primitive.ObjectID, claimID primitive.ObjectID, note string) error {
    // claim, err := s.claimRepo.FindByID(claimID)
    // if err != nil || claim.Status != models.Submitted {
    //     return errors.New("invalid operation")
    // }
    // claim.Status = models.Reviewed
    // history := models.ClaimHistory{
    //     Status:    models.Reviewed,
    //     ChangedBy: verifierID,
    //     ChangedAt: time.Now(),
    //     Note:      note,
    // }
    // s.claimRepo.AddHistory(claimID, history)
    // return s.claimRepo.Update(claim)

    claim, err := s.claimRepo.FindByID(claimID)
    if err != nil || claim.Status != models.Submitted {
        return errors.New("invalid operation")
    }

    now := time.Now()
    history := models.ClaimHistory{
        Status:    models.Reviewed,
        ChangedBy: verifierID,
        ChangedAt: now,
        Note:      note,
    }

    update := bson.M{
        "$set": bson.M{
            "status":     models.Reviewed,
            "updated_at": now,
        },
        "$push": bson.M{"history": history},
    }

    return s.claimRepo.UpdateWithPush(claimID, update)
}

// Approver only
func (s *claimService) ApproveClaim(approverID primitive.ObjectID, claimID primitive.ObjectID) error {
    // claim, err := s.claimRepo.FindByID(claimID)
    // if err != nil || claim.Status != models.Reviewed {
    //     return errors.New("invalid operation")
    // }
    // claim.Status = models.Approved
    // history := models.ClaimHistory{
    //     Status:    models.Approved,
    //     ChangedBy: approverID,
    //     ChangedAt: time.Now(),
    // }
    // s.claimRepo.AddHistory(claimID, history)
    // return s.claimRepo.Update(claim)

    claim, err := s.claimRepo.FindByID(claimID)
    if err != nil || claim.Status != models.Reviewed {
        return errors.New("invalid operation")
    }

    now := time.Now()
    update := bson.M{
        "$set": bson.M{
            "status":     models.Approved,
            "updated_at": now,
        },
        "$push": bson.M{
            "history": models.ClaimHistory{
                Status:    models.Approved,
                ChangedBy: approverID,
                ChangedAt: now,
            },
        },
    }

    return s.claimRepo.UpdateWithPush(claimID, update)
}

func (s *claimService) RejectClaim(approverID primitive.ObjectID, claimID primitive.ObjectID, reason string) error {
    // claim, err := s.claimRepo.FindByID(claimID)
    // if err != nil || claim.Status != models.Reviewed {
    //     return errors.New("invalid operation")
    // }
    // claim.Status = models.Rejected
    // history := models.ClaimHistory{
    //     Status:    models.Rejected,
    //     ChangedBy: approverID,
    //     ChangedAt: time.Now(),
    //     Note:      reason,
    // }
    // s.claimRepo.AddHistory(claimID, history)
    // return s.claimRepo.Update(claim)

    claim, err := s.claimRepo.FindByID(claimID)
    if err != nil || claim.Status != models.Reviewed {
        return errors.New("invalid operation")
    }

    now := time.Now()
    update := bson.M{
        "$set": bson.M{
            "status":     models.Rejected,
            "updated_at": now,
        },
        "$push": bson.M{
            "history": models.ClaimHistory{
                Status:    models.Rejected,
                ChangedBy: approverID,
                ChangedAt: now,
                Note:      reason,
            },
        },
    }

    return s.claimRepo.UpdateWithPush(claimID, update)
}