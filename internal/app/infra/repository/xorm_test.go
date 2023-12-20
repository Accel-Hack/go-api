package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Accel-Hack/go-api/internal/domain/sample/model"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"xorm.io/xorm"
)

// init_sample.sql inserts following records into `SAMPLE` table.
// ('00000000-0000-0000-0000-000000000000', 'test-japanese',         '1994-09-14', true,  '2003-06-14', '2003-06-14', false, null),
// ('00000000-0000-0000-0000-000000000001', 'test-deleted-japanese', '1994-10-12', true,  '2003-06-14', '2004-10-12', true, '2004-10-12'),
// ('00000000-0000-0000-0000-000000000002', 'test-deleted-foreiner', '1994-11-08', false, '2003-06-14', '2004-11-08', true, '2004-11-08'),
// ('00000000-0000-0000-0000-000000000003', 'test-foreiner',         '1994-11-08', false, '2003-06-14', '2004-06-14', false, null),
// ('00000000-0000-0000-0000-000000000004', 'test-ninja',            '1994-12-12', true,  '2003-06-14', '2004-06-14', false, null);
const (
	INIT_SCRIPT  = "./testdata/mysql/init_sample.sql"
	SAMPLE_TABLE = "SAMPLE"
)

const columns = []string{"ID", "NAME", "BIRTHDAY", "IS_JAPANESE", "CREATED_AT", "UPDATED_AT", "IS_DELETED", "DELETED_AT"}

func TestSampleXorm_FindByID(t *testing.T) {
	const query = "SELECT ID, NAME, BIRTHDAY, IS_JAPANESE FROM SAMPLE WHERE IS_DELETED = FALSE AND ID = ?"
	ctx := context.Background()
	e := setupEngine(ctx, t)
	repo := NewSampleXorm(e, SAMPLE_TABLE)
	tests := map[string]struct {
		id      string
        row *sqlmock.Rows
		want    *model.Sample
		wantErr error
	}{
		"return sample when id exists": {
			id: "00000000-0000-0000-0000-000000000000",
            row: sqlmock.NewRows("ID", "NAME", "BIRTHDAY"),
			want: &model.Sample{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				Name:       "test-japanese",
				Birthday:   time.Date(1994, 9, 14, 0, 0, 0, 0, time.Local),
				IsJapanese: true,
			},
		},
		"return err when id is deleted user": {
			id:      "00000000-0000-0000-0000-000000000001",
			wantErr: ErrNotFound,
		},
		"return err when id does not exist": {
			id:      "00000000-0000-0000-0000-000000000009",
			wantErr: ErrNotFound,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
    mock.ExpectQuery(query).WithArgs(tt.id).WillReturnRows(rows ...*sqlmock.Rows)

			got, err := repo.FindByID(context.Background(), uuid.MustParse(tt.id))
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSampleXorm_FindByNameLike(t *testing.T) {
	ctx := context.Background()
	e := setupEngine(ctx, t)
	repo := NewSampleXorm(e, SAMPLE_TABLE)
	tests := map[string]struct {
		name    string
		offset  int
		limit   int
		want    *model.PagedSamples
		wantErr error
	}{
		"return non-deleted samples when sample found": {
			name:   "test-",
			offset: 0,
			limit:  10,
			want: &model.PagedSamples{
				Total: 3,
				Samples: []model.Sample{
					{
						ID:         uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						Name:       "test-japanese",
						Birthday:   time.Date(1994, 9, 14, 0, 0, 0, 0, time.Local),
						IsJapanese: true,
					},
					{
						ID:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Name:       "test-foreiner",
						Birthday:   time.Date(1994, 11, 8, 0, 0, 0, 0, time.Local),
						IsJapanese: false,
					},
					{
						ID:         uuid.MustParse("00000000-0000-0000-0000-000000000004"),
						Name:       "test-ninja",
						Birthday:   time.Date(1994, 12, 12, 0, 0, 0, 0, time.Local),
						IsJapanese: true,
					},
				},
			},
		},
		"return samples according to limit": {
			name:   "test-",
			offset: 0,
			limit:  1,
			want: &model.PagedSamples{
				Total: 3,
				Samples: []model.Sample{
					{
						ID:         uuid.MustParse("00000000-0000-0000-0000-000000000000"),
						Name:       "test-japanese",
						Birthday:   time.Date(1994, 9, 14, 0, 0, 0, 0, time.Local),
						IsJapanese: true,
					},
				},
			},
		},
		"return samples according to offset": {
			name:   "test-",
			offset: 1,
			limit:  2,
			want: &model.PagedSamples{
				Total: 3,
				Samples: []model.Sample{
					{
						ID:         uuid.MustParse("00000000-0000-0000-0000-000000000003"),
						Name:       "test-foreiner",
						Birthday:   time.Date(1994, 11, 8, 0, 0, 0, 0, time.Local),
						IsJapanese: false,
					},
					{
						ID:         uuid.MustParse("00000000-0000-0000-0000-000000000004"),
						Name:       "test-ninja",
						Birthday:   time.Date(1994, 12, 12, 0, 0, 0, 0, time.Local),
						IsJapanese: true,
					},
				},
			},
		},
		"return empty when only deleted samples found": {
			name:   "-d-",
			offset: 0,
			limit:  10,
			want: &model.PagedSamples{
				Total:   0,
				Samples: []model.Sample{},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := repo.FindByNameLike(context.Background(), tt.name, tt.offset, tt.limit)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMain(m *testing.M) {
	if s, err := os.ReadFile(INIT_SCRIPT); err == nil {
		log.Printf("Init MySQL Container: \n%v\n", string(s))
	} else {
		log.Printf("Failed to read init script %s: %v\n", INIT_SCRIPT, err)
	}
	os.Exit(m.Run())
}

func startMySQLContainer(ctx context.Context, tb testing.TB) (*mysql.MySQLContainer, error) {
	return mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.0"),
		mysql.WithDatabase("test"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("testpass"),
		mysql.WithScripts(INIT_SCRIPT),
		withLogger(testcontainers.TestLogger(tb)),
	)
}

func withLogger(logger testcontainers.Logging) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Logger = logger
	}
}

func setupEngine(ctx context.Context, t *testing.T) *xorm.Engine {
	container, err := startMySQLContainer(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})
	connString, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}
	e, err := xorm.NewEngine("mysql", connString)
	if err != nil {
		t.Fatal(err)
	}
	return e
}
