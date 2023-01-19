package entutil_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lht102/workflow-playground/approval-service/api"
	"github.com/lht102/workflow-playground/approval-service/ent"
	"github.com/lht102/workflow-playground/approval-service/ent/payment"
	"github.com/lht102/workflow-playground/approval-service/ent/review"
	"github.com/lht102/workflow-playground/approval-service/entutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestEntityPaymentToAPIPayment(t *testing.T) {
	tests := []struct {
		input    *ent.Payment
		expected api.Payment
	}{
		{
			input: &ent.Payment{
				ID:         1,
				RequestID:  uuid.MustParse("d7cf5ded-9232-4779-930e-0348b29edec1"),
				Status:     payment.StatusREJECTED,
				Remark:     null.StringFrom("This is a remark").Ptr(),
				CreateTime: time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 1, 2, 0, 0, 0, time.UTC),
				Edges: ent.PaymentEdges{
					Reviews: []*ent.Review{
						{
							ID:         1,
							Event:      review.EventREJECT,
							ReviewerID: "alice@abc.com",
							Comment:    null.StringFrom("This a comment").Ptr(),
							CreateTime: time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC),
							UpdateTime: time.Date(2023, 1, 1, 2, 0, 0, 0, time.UTC),
						},
					},
				},
			},
			expected: api.Payment{
				Id:         1,
				RequestId:  uuid.MustParse("d7cf5ded-9232-4779-930e-0348b29edec1"),
				Status:     api.REJECTED,
				Remark:     null.StringFrom("This is a remark").Ptr(),
				CreateTime: time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC),
				UpdateTime: time.Date(2023, 1, 1, 2, 0, 0, 0, time.UTC),
				Reviews: []api.Review{
					{
						Id:         1,
						Event:      api.REJECT,
						ReviewerId: "alice@abc.com",
						Comment:    null.StringFrom("This a comment").Ptr(),
						CreateTime: time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC),
						UpdateTime: time.Date(2023, 1, 1, 2, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			actual := entutil.PaymentToAPIPayment(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
