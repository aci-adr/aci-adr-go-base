package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ForexData struct {
	ID                           primitive.ObjectID `bson:"_id"`
	Tier                         string             `bson:"tier"`
	DirectIndirectFlag           string             `bson:"directIndirectFlag"`
	Multiplier                   float64            `bson:"multiplier"`
	BuyRate                      float64            `bson:"buyRate"`
	SellRate                     float64            `bson:"sellRate"`
	TolerancePercentage          int                `bson:"tolerancePercentage"`
	EffectiveDate                *time.Time         `bson:"effectiveDate"`
	ExpirationDate               *time.Time         `bson:"expirationDate"`
	ContractRequirementThreshold string             `bson:"contractRequirementThreshold"`
	TenantID                     int                `bson:"tenantId"`
	BankID                       int                `bson:"bankId"`
	BaseCurrency                 string             `bson:"baseCurrency"`
	TargetCurrency               string             `bson:"targetCurrency"`
	CreatedDate                  time.Time          `bson:"createdDate"`
	DocVersion                   int                `bson:"docVersion"`
	UpdatedDate                  time.Time          `bson:"updatedDate"`
}
