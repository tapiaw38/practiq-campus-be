package quiz_attempt

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/tapiaw38/practiq-campus-be/internal/domain"
	"github.com/tapiaw38/practiq-campus-be/internal/platform/tenantcontext"
)

// Opt-in real PostgreSQL tests; each run owns a separate schema.
func TestPostgresLifecycle(t *testing.T) {
	dsn := os.Getenv("CAMPUS_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set CAMPUS_TEST_DATABASE_URL to a disposable PostgreSQL database")
	}
	root, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	schema := fmt.Sprintf("campus_quiz_test_%d", time.Now().UnixNano())
	if _, err = root.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	defer root.Exec("DROP SCHEMA " + schema + " CASCADE")
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	values := u.Query()
	values.Set("search_path", schema)
	u.RawQuery = values.Encode()
	db, err := sql.Open("postgres", u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.Exec(query, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`CREATE TABLE courses(id uuid PRIMARY KEY, tenant_id text);
 CREATE TABLE quizzes(id uuid PRIMARY KEY, course_id uuid REFERENCES courses);
 CREATE TABLE quiz_questions(id uuid PRIMARY KEY, quiz_id uuid REFERENCES quizzes, type text, statement text, correct_answer text, points int, position int);
 CREATE TABLE quiz_attempts(id uuid PRIMARY KEY DEFAULT gen_random_uuid(), quiz_id uuid REFERENCES quizzes, user_id text, attempt_number int, started_at timestamptz DEFAULT NOW(), expires_at timestamptz, submitted_at timestamptz, score int DEFAULT 0, max_score int DEFAULT 0, UNIQUE(quiz_id,user_id,attempt_number));
 CREATE TABLE quiz_answers(id uuid PRIMARY KEY DEFAULT gen_random_uuid(), attempt_id uuid REFERENCES quiz_attempts, question_id uuid REFERENCES quiz_questions, answer_text text, is_correct boolean, UNIQUE(attempt_id, question_id));
 INSERT INTO courses VALUES ('00000000-0000-0000-0000-000000000001','school');
 INSERT INTO quizzes VALUES ('00000000-0000-0000-0000-000000000002','00000000-0000-0000-0000-000000000001');
 INSERT INTO quiz_questions VALUES ('00000000-0000-0000-0000-000000000003','00000000-0000-0000-0000-000000000002','true_false','True?', 'true', 5, 0);
 INSERT INTO quizzes VALUES ('00000000-0000-0000-0000-000000000004','00000000-0000-0000-0000-000000000001');
 INSERT INTO quiz_questions VALUES ('00000000-0000-0000-0000-000000000005','00000000-0000-0000-0000-000000000004','true_false','Other quiz?', 'true', 5, 0);`)
	ctx := context.WithValue(context.Background(), tenantcontext.Key, "school")
	r := NewRepository(db)
	const quizID = "00000000-0000-0000-0000-000000000002"
	const questionID = "00000000-0000-0000-0000-000000000003"
	const otherQuizQuestionID = "00000000-0000-0000-0000-000000000005"
	params := StartParams{QuizID: quizID, UserID: "student", MaxAttempts: 1, WindowOpen: true}
	t.Run("concurrent starts resume one attempt", func(t *testing.T) {
		var wg sync.WaitGroup
		ids := make(chan string, 8)
		// A gate so the requests genuinely overlap. Without it the first
		// goroutine tends to commit before the others have read anything, and
		// they resume rather than collide — which makes the test pass even
		// with the serialisation removed.
		gate := make(chan struct{})
		var ready sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			ready.Add(1)
			go func() {
				defer wg.Done()
				ready.Done()
				<-gate
				result, err := r.StartOrResume(ctx, params)
				if err != nil {
					t.Error(err)
					return
				}
				ids <- result.Attempt.ID
			}()
		}
		ready.Wait()
		close(gate)
		wg.Wait()
		close(ids)
		var first string
		for id := range ids {
			if first == "" {
				first = id
			}
			if id != first {
				t.Error("created multiple attempts")
			}
		}
	})
	t.Run("start serialises on the advisory lock", func(t *testing.T) {
		// Racing eight goroutines does not actually prove the serialisation:
		// on a local database the first transaction commits before the others
		// have read anything, so they resume and the test passes even with the
		// lock removed. This asserts the lock is taken instead of hoping the
		// window opens — hold it from outside and require the start to wait.
		const probe = "lock-probe"
		holder, err := db.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer holder.Rollback()
		if _, err = holder.Exec(`SELECT pg_advisory_xact_lock(hashtext($1 || ':' || $2))`, quizID, probe); err != nil {
			t.Fatal(err)
		}

		done := make(chan error, 1)
		go func() {
			p := params
			p.UserID = probe
			_, e := r.StartOrResume(ctx, p)
			done <- e
		}()

		deadline := time.Now().Add(5 * time.Second)
		waiting := false
		for time.Now().Before(deadline) {
			if err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_stat_activity WHERE pid<>pg_backend_pid() AND wait_event_type='Lock' AND query LIKE '%pg_advisory_xact_lock%')`).Scan(&waiting); err != nil {
				t.Fatal(err)
			}
			if waiting {
				break
			}
			select {
			case e := <-done:
				t.Fatalf("start finished without waiting for the lock: %v", e)
			default:
			}
			time.Sleep(10 * time.Millisecond)
		}
		if !waiting {
			t.Fatal("start did not take the advisory lock")
		}
		if err = holder.Rollback(); err != nil {
			t.Fatal(err)
		}
		if err = <-done; err != nil {
			t.Fatal(err)
		}
	})
	result, err := r.StartOrResume(ctx, params)
	if err != nil {
		t.Fatal(err)
	}
	id := result.Attempt.ID
	t.Run("invalid batch is atomic", func(t *testing.T) {
		// A question that exists but belongs to a different quiz: the foreign
		// key is satisfied, so only the ownership check can reject it.
		err := r.SaveDraft(ctx, id, "student", []domain.QuizAnswer{{QuestionID: questionID, AnswerText: "true"}, {QuestionID: otherQuizQuestionID}})
		if err != ErrForeignQuestion {
			t.Fatalf("got %v", err)
		}
		rows, err := r.ListAnswers(ctx, id)
		if err != nil || len(rows) != 0 {
			t.Fatalf("partial write: %v %v", rows, err)
		}
	})
	if err = r.SaveDraft(ctx, id, "student", []domain.QuizAnswer{{QuestionID: questionID, AnswerText: "false"}}); err != nil {
		t.Fatal(err)
	}
	exec("UPDATE quiz_attempts SET expires_at=NOW()-interval '1 second' WHERE id=$1", id)
	t.Run("finalize waits for latest saved answer", func(t *testing.T) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		if _, err = tx.Exec("SELECT id FROM quiz_attempts WHERE id=$1 FOR UPDATE", id); err != nil {
			t.Fatal(err)
		}
		done := make(chan error, 1)
		go func() {
			settled, e := r.FinalizeExpired(ctx, id)
			if e == nil && (settled == nil || settled.Score != 5) {
				e = fmt.Errorf("stale grade: %+v", settled)
			}
			done <- e
		}()
		// Wait until the finalizer actually blocks on the writer's lock.
		deadline := time.Now().Add(5 * time.Second)
		blocked := false
		for time.Now().Before(deadline) {
			if err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_stat_activity WHERE pid<>pg_backend_pid() AND wait_event_type='Lock' AND query LIKE '%FROM quiz_attempts WHERE id = $1 FOR UPDATE%')`).Scan(&blocked); err != nil {
				t.Fatal(err)
			}
			if blocked {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if !blocked {
			t.Fatal("finalizer did not wait for lock")
		}
		if _, err = tx.Exec("UPDATE quiz_answers SET answer_text='true' WHERE attempt_id=$1", id); err != nil {
			t.Fatal(err)
		}
		if err = tx.Commit(); err != nil {
			t.Fatal(err)
		}
		if err = <-done; err != nil {
			t.Fatal(err)
		}
		again, err := r.FinalizeExpired(ctx, id)
		if err != nil || again.Score != 5 {
			t.Fatalf("not idempotent: %+v %v", again, err)
		}
	})
	t.Run("tenant isolation", func(t *testing.T) {
		other := context.WithValue(context.Background(), tenantcontext.Key, "other")
		if _, err := r.FinalizeExpired(other, id); err != ErrQuizNotInTenant {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("concurrent submissions close once", func(t *testing.T) {
		p := params
		p.UserID = "second-student"
		started, err := r.StartOrResume(ctx, p)
		if err != nil {
			t.Fatal(err)
		}
		live, err := r.FinalizeExpired(ctx, started.Attempt.ID)
		if err != nil || live.SubmittedAt != nil {
			t.Fatalf("closed live attempt: %+v %v", live, err)
		}
		outcomes := make(chan error, 2)
		for i := 0; i < 2; i++ {
			go func() {
				_, err := r.SubmitGraded(ctx, started.Attempt.ID, p.UserID, []domain.QuizAnswer{{QuestionID: questionID, AnswerText: "true", IsCorrect: true}}, 5, 5)
				outcomes <- err
			}()
		}
		successes := 0
		for i := 0; i < 2; i++ {
			err := <-outcomes
			if err == nil {
				successes++
			} else if err != ErrAttemptAlreadySubmitted {
				t.Fatal(err)
			}
		}
		if successes != 1 {
			t.Fatalf("got %d successful submissions", successes)
		}
		if err := r.SaveDraft(ctx, started.Attempt.ID, p.UserID, []domain.QuizAnswer{{QuestionID: questionID, AnswerText: "false"}}); err != ErrAttemptAlreadySubmitted {
			t.Fatalf("draft changed submitted attempt: %v", err)
		}
	})
}
