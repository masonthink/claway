package internal_test

import (
	"context"
	"math"
	"strings"
	"testing"

	"github.com/claway/server/internal/model"
	"github.com/claway/server/internal/service"
	"github.com/claway/server/internal/testutil"
)

func TestFullWorkflow(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	defer pool.Close()
	testutil.CleanupDB(t, pool)

	ctx := context.Background()
	svc := testutil.CreateTestService(t, pool)

	// Shared state across subtests
	var (
		initiatorID   int64
		contributorID int64
		stdIdeaID     int64
		stdTasks      []*model.Task
		prdID         int64
		prdPrice      float64
		totalCostUSD  float64
	)

	// ---------------------------------------------------------------
	// Setup: Create 2 test users
	// ---------------------------------------------------------------
	t.Run("Setup", func(t *testing.T) {
		initiatorID = testutil.CreateTestUser(t, pool, "openclaw-initiator-001", "initiator_alice")
		contributorID = testutil.CreateTestUser(t, pool, "openclaw-contributor-001", "contributor_bob")
		if initiatorID == 0 || contributorID == 0 {
			t.Fatal("user IDs should be non-zero")
		}
	})

	// ---------------------------------------------------------------
	// 1. CreateIdea: 4 tasks (D1-D4) + 4 documents
	// ---------------------------------------------------------------
	t.Run("CreateIdea", func(t *testing.T) {
		idea, err := svc.CreateIdea(ctx, initiatorID, service.CreateIdeaRequest{
			Title:               "Test Idea",
			Description:         "An idea for integration testing",
			TargetUserHint:      "developers",
			ProblemDefinition:   "need better tools",
			InitiatorCutPercent: 20,
		})
		if err != nil {
			t.Fatalf("CreateIdea failed: %v", err)
		}
		stdIdeaID = idea.ID

		if idea.Status != model.IdeaStatusActive {
			t.Fatalf("expected status active, got %s", idea.Status)
		}

		tasks, err := svc.ListTasks(ctx, stdIdeaID)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
		if len(tasks) != 4 {
			t.Fatalf("expected 4 tasks, got %d", len(tasks))
		}
		stdTasks = tasks

		// Verify each task has a document
		for _, task := range tasks {
			doc, err := svc.GetDocument(ctx, task.ID)
			if err != nil {
				t.Fatalf("GetDocument for task %s failed: %v", task.Type, err)
			}
			if doc.TaskID != task.ID {
				t.Fatalf("document task_id mismatch: expected %d, got %d", task.ID, doc.TaskID)
			}
		}
	})

	// ---------------------------------------------------------------
	// 2. CreateIdea_Validation: Test validation errors
	// ---------------------------------------------------------------
	t.Run("CreateIdea_Validation", func(t *testing.T) {
		// Empty title
		_, err := svc.CreateIdea(ctx, initiatorID, service.CreateIdeaRequest{
			Title:               "",
			Description:         "desc",
			InitiatorCutPercent: 20,
		})
		if err == nil {
			t.Fatal("expected error for empty title")
		}

		// Empty description
		_, err = svc.CreateIdea(ctx, initiatorID, service.CreateIdeaRequest{
			Title:               "title",
			Description:         "",
			InitiatorCutPercent: 20,
		})
		if err == nil {
			t.Fatal("expected error for empty description")
		}
	})

	// ---------------------------------------------------------------
	// 3. ListIdeas: Verify pagination and count
	// ---------------------------------------------------------------
	t.Run("ListIdeas", func(t *testing.T) {
		resp, err := svc.ListIdeas(ctx, "active", 10, 0)
		if err != nil {
			t.Fatalf("ListIdeas failed: %v", err)
		}
		if resp.Total < 1 {
			t.Fatalf("expected at least 1 idea, got %d", resp.Total)
		}

		// Test pagination: limit=1
		resp2, err := svc.ListIdeas(ctx, "active", 1, 0)
		if err != nil {
			t.Fatalf("ListIdeas (paginated) failed: %v", err)
		}
		if len(resp2.Ideas) != 1 {
			t.Fatalf("expected 1 idea with limit=1, got %d", len(resp2.Ideas))
		}
		if resp2.Total != resp.Total {
			t.Fatalf("total count should be same regardless of limit: expected %d, got %d", resp.Total, resp2.Total)
		}
	})

	// ---------------------------------------------------------------
	// Helper: find task by type in stdTasks
	// ---------------------------------------------------------------
	findTask := func(taskType model.TaskType) *model.Task {
		for _, task := range stdTasks {
			if task.Type == taskType {
				return task
			}
		}
		t.Fatalf("task %s not found", taskType)
		return nil
	}

	d1Task := findTask(model.TaskTypeD1)
	d2Task := findTask(model.TaskTypeD2)
	d3Task := findTask(model.TaskTypeD3)
	d4Task := findTask(model.TaskTypeD4)

	// ---------------------------------------------------------------
	// 4. ClaimTask_D1: No dependencies, should succeed
	// ---------------------------------------------------------------
	t.Run("ClaimTask_D1", func(t *testing.T) {
		err := svc.ClaimTask(ctx, d1Task.ID, contributorID)
		if err != nil {
			t.Fatalf("ClaimTask D1 failed: %v", err)
		}

		task, err := svc.GetTask(ctx, d1Task.ID)
		if err != nil {
			t.Fatalf("GetTask D1 failed: %v", err)
		}
		if task.Status != model.TaskStatusClaimed {
			t.Fatalf("expected D1 status claimed, got %s", task.Status)
		}
		if !task.ClaimedBy.Valid || task.ClaimedBy.Int64 != contributorID {
			t.Fatalf("expected D1 claimed_by %d, got %v", contributorID, task.ClaimedBy)
		}
	})

	// ---------------------------------------------------------------
	// 5. ClaimTask_D3_Blocked: D1,D2 not approved yet
	// ---------------------------------------------------------------
	t.Run("ClaimTask_D3_Blocked", func(t *testing.T) {
		err := svc.ClaimTask(ctx, d3Task.ID, contributorID)
		if err == nil {
			t.Fatal("expected error when claiming D3 with unapproved dependencies")
		}
		if !strings.Contains(err.Error(), "not yet approved") {
			t.Fatalf("expected dependency error, got: %v", err)
		}
	})

	// ---------------------------------------------------------------
	// 6. ClaimTask_D4_Blocked: D3 not approved yet
	// ---------------------------------------------------------------
	t.Run("ClaimTask_D4_Blocked", func(t *testing.T) {
		err := svc.ClaimTask(ctx, d4Task.ID, contributorID)
		if err == nil {
			t.Fatal("expected error when claiming D4 with unapproved dependencies")
		}
		if !strings.Contains(err.Error(), "not yet approved") {
			t.Fatalf("expected dependency error, got: %v", err)
		}
	})

	// ---------------------------------------------------------------
	// 7. UpdateDocument: Contributor updates D1 document
	// ---------------------------------------------------------------
	t.Run("UpdateDocument_D1", func(t *testing.T) {
		err := svc.UpdateDocument(ctx, d1Task.ID, contributorID, service.UpdateDocumentRequest{
			Content: "D1 competitive analysis draft content",
		})
		if err != nil {
			t.Fatalf("UpdateDocument D1 failed: %v", err)
		}

		doc, err := svc.GetDocument(ctx, d1Task.ID)
		if err != nil {
			t.Fatalf("GetDocument D1 failed: %v", err)
		}
		if doc.Content != "D1 competitive analysis draft content" {
			t.Fatalf("expected updated content, got: %s", doc.Content)
		}
		if doc.CurrentVersion != 2 {
			t.Fatalf("expected version 2, got %d", doc.CurrentVersion)
		}

		// Verify version record was created
		versions, err := svc.ListDocumentVersions(ctx, d1Task.ID)
		if err != nil {
			t.Fatalf("ListDocumentVersions failed: %v", err)
		}
		if len(versions) < 1 {
			t.Fatal("expected at least 1 document version")
		}
	})

	// ---------------------------------------------------------------
	// 8. SubmitTask_D1_WithTokenUsage: Submit with self-reported token usage
	// ---------------------------------------------------------------
	t.Run("SubmitTask_D1_WithTokenUsage", func(t *testing.T) {
		err := svc.SubmitTask(ctx, d1Task.ID, contributorID, service.SubmitTaskRequest{
			Content: "D1 final competitive analysis content",
			Note:    "Completed analysis",
			TokenUsage: &service.TokenUsageReport{
				Model:     "claude-sonnet-4-20250514",
				TokensIn:  45000,
				TokensOut: 12000,
				CostUSD:   0.05,
			},
		})
		if err != nil {
			t.Fatalf("SubmitTask D1 failed: %v", err)
		}

		task, err := svc.GetTask(ctx, d1Task.ID)
		if err != nil {
			t.Fatalf("GetTask D1 failed: %v", err)
		}
		if task.Status != model.TaskStatusSubmitted {
			t.Fatalf("expected D1 status submitted, got %s", task.Status)
		}
		// Token usage should have been accumulated
		if task.CostUSDAccumulated < 0.049 {
			t.Fatalf("expected cost_usd_accumulated >= 0.05, got %f", task.CostUSDAccumulated)
		}
	})

	// ---------------------------------------------------------------
	// 9. ReviewTask_Revision_D1: Initiator requests revision
	// ---------------------------------------------------------------
	t.Run("ReviewTask_Revision_D1", func(t *testing.T) {
		err := svc.ReviewTask(ctx, d1Task.ID, initiatorID, service.ReviewTaskRequest{
			Action:   "revision",
			Feedback: "Please add pricing comparison table for top 3 competitors.",
		})
		if err != nil {
			t.Fatalf("ReviewTask revision D1 failed: %v", err)
		}

		task, err := svc.GetTask(ctx, d1Task.ID)
		if err != nil {
			t.Fatalf("GetTask D1 failed: %v", err)
		}
		if task.Status != model.TaskStatusRevision {
			t.Fatalf("expected D1 status revision, got %s", task.Status)
		}
		if !task.ReviewFeedback.Valid || !strings.Contains(task.ReviewFeedback.String, "pricing comparison") {
			t.Fatalf("expected review_feedback to contain feedback, got: %v", task.ReviewFeedback)
		}
	})

	// ---------------------------------------------------------------
	// 10. UpdateDocument_During_Revision: Contributor can save drafts in revision status
	// ---------------------------------------------------------------
	t.Run("UpdateDocument_During_Revision", func(t *testing.T) {
		err := svc.UpdateDocument(ctx, d1Task.ID, contributorID, service.UpdateDocumentRequest{
			Content: "D1 revised draft with pricing table",
		})
		if err != nil {
			t.Fatalf("UpdateDocument during revision failed: %v", err)
		}
	})

	// ---------------------------------------------------------------
	// 11. Resubmit_D1: Contributor resubmits after revision
	// ---------------------------------------------------------------
	t.Run("Resubmit_D1", func(t *testing.T) {
		err := svc.SubmitTask(ctx, d1Task.ID, contributorID, service.SubmitTaskRequest{
			Content: "D1 final content with pricing comparison table added",
			Note:    "Added pricing table as requested",
			TokenUsage: &service.TokenUsageReport{
				Model:     "claude-sonnet-4-20250514",
				TokensIn:  10000,
				TokensOut: 5000,
				CostUSD:   0.02,
			},
		})
		if err != nil {
			t.Fatalf("Resubmit D1 failed: %v", err)
		}

		task, err := svc.GetTask(ctx, d1Task.ID)
		if err != nil {
			t.Fatalf("GetTask D1 failed: %v", err)
		}
		if task.Status != model.TaskStatusSubmitted {
			t.Fatalf("expected D1 status submitted after resubmit, got %s", task.Status)
		}
		// Cost should have accumulated (0.05 + 0.02)
		if task.CostUSDAccumulated < 0.069 {
			t.Fatalf("expected accumulated cost >= 0.07, got %f", task.CostUSDAccumulated)
		}
	})

	// ---------------------------------------------------------------
	// 12. ReviewTask_Approve_D1: Initiator approves D1
	// ---------------------------------------------------------------
	t.Run("ReviewTask_Approve_D1", func(t *testing.T) {
		contributorBalanceBefore := testutil.GetUserBalance(t, pool, contributorID)

		err := svc.ReviewTask(ctx, d1Task.ID, initiatorID, service.ReviewTaskRequest{
			Action:       "approve",
			QualityScore: 1.2,
		})
		if err != nil {
			t.Fatalf("ReviewTask approve D1 failed: %v", err)
		}

		// Verify task status
		task, err := svc.GetTask(ctx, d1Task.ID)
		if err != nil {
			t.Fatalf("GetTask D1 after approval failed: %v", err)
		}
		if task.Status != model.TaskStatusApproved {
			t.Fatalf("expected D1 status approved, got %s", task.Status)
		}

		// Verify contribution record
		contribs, err := svc.GetMyContributions(ctx, contributorID)
		if err != nil {
			t.Fatalf("GetMyContributions failed: %v", err)
		}
		found := false
		for _, c := range contribs {
			if c.TaskID == d1Task.ID {
				found = true
				if c.QualityScore != 1.2 {
					t.Fatalf("expected quality_score 1.2, got %f", c.QualityScore)
				}
				break
			}
		}
		if !found {
			t.Fatal("contribution record for D1 not found")
		}

		// Verify credits awarded: cost * quality_score * 1000
		expectedCredits := task.CostUSDAccumulated * 1.2 * 1000
		contributorBalanceAfter := testutil.GetUserBalance(t, pool, contributorID)
		actualDelta := contributorBalanceAfter - contributorBalanceBefore
		if math.Abs(actualDelta-expectedCredits) > 0.1 {
			t.Fatalf("expected credits delta ~%.2f, got %.2f", expectedCredits, actualDelta)
		}

		// Verify credit transaction recorded
		creditsResp, err := svc.GetMyCredits(ctx, contributorID, 10, 0)
		if err != nil {
			t.Fatalf("GetMyCredits failed: %v", err)
		}
		foundTx := false
		for _, tx := range creditsResp.Transactions {
			if tx.ReferenceType == "task" && tx.ReferenceID == d1Task.ID && tx.Type == "earn_contribute" {
				foundTx = true
				break
			}
		}
		if !foundTx {
			t.Fatal("credit transaction for D1 approval not found")
		}
	})

	// ---------------------------------------------------------------
	// 13. ReviewTask_NotInitiator: Contributor tries to review
	// ---------------------------------------------------------------
	t.Run("ReviewTask_NotInitiator", func(t *testing.T) {
		// D2 has no dependencies, claim and submit
		err := svc.ClaimTask(ctx, d2Task.ID, contributorID)
		if err != nil {
			t.Fatalf("ClaimTask D2 failed: %v", err)
		}
		err = svc.SubmitTask(ctx, d2Task.ID, contributorID, service.SubmitTaskRequest{
			Content: "D2 user persona content",
			Note:    "Done",
			TokenUsage: &service.TokenUsageReport{
				Model:     "claude-sonnet-4-20250514",
				TokensIn:  30000,
				TokensOut: 10000,
				CostUSD:   0.03,
			},
		})
		if err != nil {
			t.Fatalf("SubmitTask D2 failed: %v", err)
		}

		err = svc.ReviewTask(ctx, d2Task.ID, contributorID, service.ReviewTaskRequest{
			Action:       "approve",
			QualityScore: 1.0,
		})
		if err == nil {
			t.Fatal("expected error when non-initiator tries to review")
		}
		if !strings.Contains(err.Error(), "initiator") {
			t.Fatalf("expected initiator error, got: %v", err)
		}
	})

	// ---------------------------------------------------------------
	// 14. GetIdeaContext: D1 approved content included, others not
	// ---------------------------------------------------------------
	t.Run("GetIdeaContext", func(t *testing.T) {
		ctxResp, err := svc.GetIdeaContext(ctx, stdIdeaID)
		if err != nil {
			t.Fatalf("GetIdeaContext failed: %v", err)
		}
		if ctxResp.IdeaID != stdIdeaID {
			t.Fatalf("expected idea_id %d, got %d", stdIdeaID, ctxResp.IdeaID)
		}

		for _, entry := range ctxResp.Entries {
			if entry.TaskType == model.TaskTypeD1 {
				// D1 is approved, should have content
				if entry.Content == "" {
					t.Fatal("expected D1 approved content in context, got empty")
				}
			} else {
				// Other tasks are not approved, should not have content
				if entry.Content != "" {
					t.Fatalf("expected no content for non-approved task %s, got: %s", entry.TaskType, entry.Content)
				}
			}
		}
	})

	// ---------------------------------------------------------------
	// 15. Approve D2, then claim+submit+approve D3 and D4
	// ---------------------------------------------------------------
	t.Run("SubmitAndApprove_D2", func(t *testing.T) {
		// D2 is already submitted from the NotInitiator test above.
		err := svc.ReviewTask(ctx, d2Task.ID, initiatorID, service.ReviewTaskRequest{
			Action:       "approve",
			QualityScore: 1.0,
		})
		if err != nil {
			t.Fatalf("ReviewTask approve D2 failed: %v", err)
		}

		task, err := svc.GetTask(ctx, d2Task.ID)
		if err != nil {
			t.Fatalf("GetTask D2 failed: %v", err)
		}
		if task.Status != model.TaskStatusApproved {
			t.Fatalf("expected D2 status approved, got %s", task.Status)
		}
	})

	// ---------------------------------------------------------------
	// 16. ClaimTask_D3_Now: D1 and D2 are approved, D3 should work
	// ---------------------------------------------------------------
	t.Run("ClaimTask_D3_Now", func(t *testing.T) {
		err := svc.ClaimTask(ctx, d3Task.ID, contributorID)
		if err != nil {
			t.Fatalf("ClaimTask D3 failed (should succeed now): %v", err)
		}

		task, err := svc.GetTask(ctx, d3Task.ID)
		if err != nil {
			t.Fatalf("GetTask D3 failed: %v", err)
		}
		if task.Status != model.TaskStatusClaimed {
			t.Fatalf("expected D3 status claimed, got %s", task.Status)
		}
	})

	// ---------------------------------------------------------------
	// 17. SubmitAndApprove_D3_D4: Complete remaining tasks
	// ---------------------------------------------------------------
	t.Run("SubmitAndApprove_D3_D4", func(t *testing.T) {
		remainingTypes := []model.TaskType{model.TaskTypeD3, model.TaskTypeD4}

		for _, tt := range remainingTypes {
			task := findTask(tt)

			// D3 is already claimed above; D4 needs to be claimed
			if tt != model.TaskTypeD3 {
				err := svc.ClaimTask(ctx, task.ID, contributorID)
				if err != nil {
					t.Fatalf("ClaimTask %s failed: %v", tt, err)
				}
			}

			content := string(tt) + " deliverable content for integration test"
			cost := 0.02 + float64(task.ID%10)*0.005
			err := svc.SubmitTask(ctx, task.ID, contributorID, service.SubmitTaskRequest{
				Content: content,
				Note:    string(tt) + " done",
				TokenUsage: &service.TokenUsageReport{
					Model:     "claude-sonnet-4-20250514",
					TokensIn:  40000,
					TokensOut: 15000,
					CostUSD:   cost,
				},
			})
			if err != nil {
				t.Fatalf("SubmitTask %s failed: %v", tt, err)
			}

			err = svc.ReviewTask(ctx, task.ID, initiatorID, service.ReviewTaskRequest{
				Action:       "approve",
				QualityScore: 1.0,
			})
			if err != nil {
				t.Fatalf("ReviewTask approve %s failed: %v", tt, err)
			}
		}

		// Verify all 4 tasks are approved
		tasks, err := svc.ListTasks(ctx, stdIdeaID)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
		for _, task := range tasks {
			if task.Status != model.TaskStatusApproved {
				t.Fatalf("expected task %s to be approved, got %s", task.Type, task.Status)
			}
		}
	})

	// ---------------------------------------------------------------
	// 18. PublishPRD: Initiator publishes PRD
	// ---------------------------------------------------------------
	t.Run("PublishPRD", func(t *testing.T) {
		prd, err := svc.PublishPRD(ctx, stdIdeaID, initiatorID)
		if err != nil {
			t.Fatalf("PublishPRD failed: %v", err)
		}
		prdID = prd.ID

		// Verify PRD has merged content
		if prd.Content == "" {
			t.Fatal("expected PRD content to be non-empty")
		}
		if !strings.Contains(prd.Content, "Test Idea") {
			t.Fatal("expected PRD content to contain idea title")
		}
		// Should contain content from each deliverable
		for _, tt := range []string{"doc1", "doc2", "doc3", "doc4"} {
			if !strings.Contains(prd.Content, tt) {
				t.Fatalf("expected PRD content to reference %s", tt)
			}
		}

		// Calculate expected price: total_cost * 2 * 1000
		tasks, err := svc.ListTasks(ctx, stdIdeaID)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}
		totalCostUSD = 0
		for _, task := range tasks {
			totalCostUSD += task.CostUSDAccumulated
		}
		expectedPrice := totalCostUSD * 2 * 1000
		prdPrice = prd.PriceCredits

		if math.Abs(prdPrice-expectedPrice) > 0.1 {
			t.Fatalf("expected PRD price ~%.2f, got %.2f", expectedPrice, prdPrice)
		}

		// Verify idea status changed to completed
		idea, err := svc.GetIdea(ctx, stdIdeaID)
		if err != nil {
			t.Fatalf("GetIdea failed: %v", err)
		}
		if idea.Status != model.IdeaStatusCompleted {
			t.Fatalf("expected idea status completed, got %s", idea.Status)
		}
	})

	// ---------------------------------------------------------------
	// 19. PurchasePRD: Contributor purchases the PRD
	// ---------------------------------------------------------------
	t.Run("PurchasePRD", func(t *testing.T) {
		// Give contributor enough credits
		testutil.GiveCredits(t, pool, contributorID, prdPrice+1000)

		contributorBefore := testutil.GetUserBalance(t, pool, contributorID)
		initiatorBefore := testutil.GetUserBalance(t, pool, initiatorID)

		err := svc.PurchasePRD(ctx, contributorID, prdID)
		if err != nil {
			t.Fatalf("PurchasePRD failed: %v", err)
		}

		// Verify net change for contributor (who is both buyer and sole contributor):
		// Net = -prdPrice + contributorsPool (70% of prdPrice)
		// Net = -prdPrice * 0.30 (they lose 30%: 10% platform + 20% initiator)
		contributorAfter := testutil.GetUserBalance(t, pool, contributorID)
		netDeducted := contributorBefore - contributorAfter
		expectedNetDeduction := prdPrice * 0.30 // buyer loses 30% (platform 10% + initiator 20%)
		if math.Abs(netDeducted-expectedNetDeduction) > 0.1 {
			t.Fatalf("expected net deduction ~%.2f, got %.2f", expectedNetDeduction, netDeducted)
		}

		// Verify initiator received their cut (20%)
		expectedInitiatorShare := prdPrice * 0.20
		initiatorAfter := testutil.GetUserBalance(t, pool, initiatorID)
		initiatorDelta := initiatorAfter - initiatorBefore
		if math.Abs(initiatorDelta-expectedInitiatorShare) > 0.1 {
			t.Fatalf("expected initiator share ~%.2f, got %.2f", expectedInitiatorShare, initiatorDelta)
		}

		// Verify contributor received pool via transactions
		creditsResp, err := svc.GetMyCredits(ctx, contributorID, 50, 0)
		if err != nil {
			t.Fatalf("GetMyCredits failed: %v", err)
		}

		var earnShareTotal float64
		for _, tx := range creditsResp.Transactions {
			if tx.Type == "earn_read_share" && tx.ReferenceType == "prd" && tx.ReferenceID == prdID {
				earnShareTotal += tx.Amount
			}
		}

		expectedContributorsPool := prdPrice * 0.70
		if math.Abs(earnShareTotal-expectedContributorsPool) > 0.1 {
			t.Fatalf("expected contributor earn_read_share total ~%.2f, got %.2f", expectedContributorsPool, earnShareTotal)
		}
	})

	// ---------------------------------------------------------------
	// 20. PurchasePRD_AlreadyPurchased
	// ---------------------------------------------------------------
	t.Run("PurchasePRD_AlreadyPurchased", func(t *testing.T) {
		err := svc.PurchasePRD(ctx, contributorID, prdID)
		if err == nil {
			t.Fatal("expected error for duplicate purchase")
		}
		if !strings.Contains(err.Error(), "already purchased") {
			t.Fatalf("expected 'already purchased' error, got: %v", err)
		}
	})

	// ---------------------------------------------------------------
	// 21. PurchasePRD_InsufficientCredits
	// ---------------------------------------------------------------
	t.Run("PurchasePRD_InsufficientCredits", func(t *testing.T) {
		brokeUserID := testutil.CreateTestUser(t, pool, "openclaw-broke-001", "broke_charlie")

		err := svc.PurchasePRD(ctx, brokeUserID, prdID)
		if err == nil {
			t.Fatal("expected error for insufficient credits")
		}
		if !strings.Contains(err.Error(), "insufficient credits") {
			t.Fatalf("expected 'insufficient credits' error, got: %v", err)
		}
	})

	// ---------------------------------------------------------------
	// 22. ComputeStats: Verify compute aggregation endpoints
	// ---------------------------------------------------------------
	t.Run("ComputeStats", func(t *testing.T) {
		// User compute total
		userCompute, err := svc.GetMyCompute(ctx, contributorID)
		if err != nil {
			t.Fatalf("GetMyCompute failed: %v", err)
		}
		if userCompute.TotalCostUSD <= 0 {
			t.Fatal("expected non-zero compute total for contributor")
		}

		// User idea compute
		userIdeaCompute, err := svc.GetMyIdeaCompute(ctx, contributorID, stdIdeaID)
		if err != nil {
			t.Fatalf("GetMyIdeaCompute failed: %v", err)
		}
		if userIdeaCompute.TotalCostUSD <= 0 {
			t.Fatal("expected non-zero idea compute total")
		}
		if userIdeaCompute.IdeaID != stdIdeaID {
			t.Fatalf("expected idea_id %d, got %d", stdIdeaID, userIdeaCompute.IdeaID)
		}

		// Idea compute breakdown (per contributor)
		ideaCompute, err := svc.GetIdeaCompute(ctx, stdIdeaID)
		if err != nil {
			t.Fatalf("GetIdeaCompute failed: %v", err)
		}
		if len(ideaCompute.Breakdown) == 0 {
			t.Fatal("expected non-empty idea compute breakdown")
		}
		// Only contributor_bob used tokens
		if ideaCompute.Breakdown[0].UserID != contributorID {
			t.Fatalf("expected contributor user_id %d, got %d", contributorID, ideaCompute.Breakdown[0].UserID)
		}
		if ideaCompute.Breakdown[0].TotalCost <= 0 {
			t.Fatal("expected non-zero total cost in breakdown")
		}
		// D1 submitted twice (initial + resubmit) + D2 + D3 + D4 = 5 calls
		if ideaCompute.Breakdown[0].CallCount != 5 {
			t.Fatalf("expected 5 API calls (D1 x2 + D2 + D3 + D4), got %d", ideaCompute.Breakdown[0].CallCount)
		}

		// Task compute breakdown
		taskCompute, err := svc.GetTaskCompute(ctx, d1Task.ID)
		if err != nil {
			t.Fatalf("GetTaskCompute failed: %v", err)
		}
		if len(taskCompute.Breakdown) == 0 {
			t.Fatal("expected non-empty task compute breakdown")
		}
		// D1 had 2 submissions (initial + resubmit after revision)
		if taskCompute.Breakdown[0].CallCount != 2 {
			t.Fatalf("expected 2 API calls for D1 (initial + resubmit), got %d", taskCompute.Breakdown[0].CallCount)
		}

		// Platform compute total
		platformCompute, err := svc.GetPlatformCompute(ctx)
		if err != nil {
			t.Fatalf("GetPlatformCompute failed: %v", err)
		}
		if platformCompute.TotalCostUSD <= 0 {
			t.Fatal("expected non-zero platform compute total")
		}
	})

	// ---------------------------------------------------------------
	// 23. Revision_Validation: Revision without feedback should fail
	// ---------------------------------------------------------------
	t.Run("Revision_Validation", func(t *testing.T) {
		// Create another idea to test revision validation
		idea2, err := svc.CreateIdea(ctx, initiatorID, service.CreateIdeaRequest{
			Title:               "Revision Validation Idea",
			Description:         "Testing revision validation",
			InitiatorCutPercent: 20,
		})
		if err != nil {
			t.Fatalf("CreateIdea failed: %v", err)
		}

		tasks2, err := svc.ListTasks(ctx, idea2.ID)
		if err != nil {
			t.Fatalf("ListTasks failed: %v", err)
		}

		// Claim and submit D1
		d1 := tasks2[0]
		if err := svc.ClaimTask(ctx, d1.ID, contributorID); err != nil {
			t.Fatalf("ClaimTask failed: %v", err)
		}
		if err := svc.SubmitTask(ctx, d1.ID, contributorID, service.SubmitTaskRequest{
			Content: "test content",
			Note:    "test",
		}); err != nil {
			t.Fatalf("SubmitTask failed: %v", err)
		}

		// Try revision without feedback
		err = svc.ReviewTask(ctx, d1.ID, initiatorID, service.ReviewTaskRequest{
			Action:   "revision",
			Feedback: "",
		})
		if err == nil {
			t.Fatal("expected error for revision without feedback")
		}
		if !strings.Contains(err.Error(), "feedback is required") {
			t.Fatalf("expected 'feedback is required' error, got: %v", err)
		}
	})
}
