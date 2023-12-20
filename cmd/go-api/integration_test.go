package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// init_sample.sql inserts following rows into `SAMPLE` table.
// ('00000000-0000-0000-0000-000000000000', 'test-japanese',         '1994-09-14', true,  '2003-06-14', '2003-06-14', false, null),
// ('00000000-0000-0000-0000-000000000001', 'test-deleted-japanese', '1994-10-12', true,  '2003-06-14', '2004-10-12', true, '2004-10-12'),
// ('00000000-0000-0000-0000-000000000002', 'test-deleted-foreiner', '1994-11-08', false, '2003-06-14', '2004-11-08', true, '2004-11-08'),
// ('00000000-0000-0000-0000-000000000003', 'test-foreiner',         '1994-11-08', false, '2003-06-14', '2004-06-14', false, null),
// ('00000000-0000-0000-0000-000000000004', 'test-ninja',            '1994-12-12', true,  '2003-06-14', '2004-06-14', false, null);
const (
	INIT_SCRIPT  = "./testdata/mysql/init_sample.sql"
	SAMPLE_TABLE = "SAMPLE"
)

const (
	SampleJSON_0 = `{"ID":"00000000-0000-0000-0000-000000000000","Name":"test-japanese","Birthday":"1994-09-14T00:00:00+09:00","IsJapanese":true}`
	SampleJSON_1 = `{"ID":"00000000-0000-0000-0000-000000000001","Name":"test-deleted-japanese","Birthday":"1994-10-12T00:00:00+09:00","IsJapanese":true}`
	SampleJSON_2 = `{"ID":"00000000-0000-0000-0000-000000000002","Name":"test-deleted-foreiner","Birthday":"1994-11-08T00:00:00+09:00","IsJapanese":false}`
	SampleJSON_3 = `{"ID":"00000000-0000-0000-0000-000000000003","Name":"test-foreiner","Birthday":"1994-11-08T00:00:00+09:00","IsJapanese":false}`
	SampleJSON_4 = `{"ID":"00000000-0000-0000-0000-000000000004","Name":"test-ninja","Birthday":"1994-12-12T00:00:00+09:00","IsJapanese":true}`
)

var httpClient = &http.Client{}

func TestGoAPIOption_Run_GET_Sample(t *testing.T) {
	appCtx := context.Background()
	setup(context.Background(), appCtx, t)
	type testcase struct {
		url      string
		want     string
		wantCode int
	}
	tests := map[string]testcase{
		"existing sample id": {
			url:      "http://localhost:8080/sample?id=00000000-0000-0000-0000-000000000000",
			want:     SampleJSON_0 + "\n",
			wantCode: http.StatusOK,
		},
		"deleted sample id": {
			url:      "http://localhost:8080/sample?id=00000000-0000-0000-0000-000000000001",
			wantCode: http.StatusInternalServerError,
		},
		"non-existing sample id": {
			url:      "http://localhost:8080/sample?id=00000000-0000-0000-0000-000000000010",
			wantCode: http.StatusInternalServerError,
		},
		"no id": {
			url:      "http://localhost:8080/sample",
			wantCode: http.StatusBadRequest,
		},
		"invalid id": {
			url:      "http://localhost:8080/sample?id=invalid-id",
			wantCode: http.StatusBadRequest,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			assert.NoError(t, err)
			doWithAssert(req, tt.want, tt.wantCode, t)
		})
	}
}

func TestGoAPIOption_Run_GET_Samples(t *testing.T) {
	setup(context.Background(), context.Background(), t)
	type testcase struct {
		url       string
		want      []string
		wantTotal int
	}
	tests := map[string]testcase{
		"no query": {
			url:       "http://localhost:8080/samples",
			want:      []string{SampleJSON_0, SampleJSON_3, SampleJSON_4},
			wantTotal: 3,
		},
		"limit=1": {
			url:       "http://localhost:8080/samples?limit=1",
			want:      []string{SampleJSON_0},
			wantTotal: 3,
		},
		"offset=2": {
			url:       "http://localhost:8080/samples?offset=2",
			want:      []string{SampleJSON_4},
			wantTotal: 3,
		},
		"name=foreiner": {
			url:       "http://localhost:8080/samples?name=foreiner",
			want:      []string{SampleJSON_3},
			wantTotal: 1,
		},
		"name=deleted": {
			url:       "http://localhost:8080/samples?name=deleted",
			want:      []string{},
			wantTotal: 0,
		},
		"name=does-not-exist": {
			url:       "http://localhost:8080/samples?name=does-not-exist",
			want:      []string{},
			wantTotal: 0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			assert.NoError(t, err)
			want := fmt.Sprintf(`{"Total":%d,"Samples":[%s]}`+"\n", tt.wantTotal, strings.Join(tt.want, ","))
			doWithAssert(req, want, http.StatusOK, t)
		})
	}
}

func TestGoAPIOption_Run_PUT_Sample(t *testing.T) {
	type testcase struct {
		url         string
		wantPutCode int
		assertAfter assertByIDFunc
	}
	tests := map[string]testcase{
		"valid parameters": {
			url:         "http://localhost:8080/sample?name=put-test&birthday=2000-01-01&is_japanese=true",
			wantPutCode: http.StatusOK,
			assertAfter: func(id string, t *testing.T) {
				want := sampleJSON(id, "put-test", "2000-01-01T00:00:00+09:00", true) + "\n"
				assertGetByID(id, want, http.StatusOK, t)
			},
		},
		"empty name": {
			url:         "http://localhost:8080/sample?name=&birthday=2000-01-01&is_japanese=true",
			wantPutCode: http.StatusOK,
			assertAfter: func(id string, t *testing.T) {
				want := sampleJSON(id, "", "2000-01-01T00:00:00+09:00", true) + "\n"
				assertGetByID(id, want, http.StatusOK, t)
			},
		},
		"nil name": {
			url:         "http://localhost:8080/sample?birthday=2000-01-01&is_japanese=true",
			wantPutCode: http.StatusBadRequest,
		},
		"empty birthday": {
			url:         "http://localhost:8080/sample?name=put-test&birthday=&is_japanese=true",
			wantPutCode: http.StatusBadRequest,
		},
		"nil birthday": {
			url:         "http://localhost:8080/sample?name=put-test&is_japanese=true",
			wantPutCode: http.StatusBadRequest,
		},
		"empty is_japanese": {
			url:         "http://localhost:8080/sample?name=put-test&birthday=2000-01-01&is_japanese=",
			wantPutCode: http.StatusBadRequest,
		},
		"nil is_japanese": {
			url:         "http://localhost:8080/sample?name=put-test&birthday=2000-01-01",
			wantPutCode: http.StatusBadRequest,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			setup(context.Background(), context.Background(), t)
			// run
			req, err := http.NewRequest(http.MethodPut, tt.url, nil)
			assert.NoError(t, err)
			putResp, err := httpClient.Do(req)
			assert.NoError(t, err)
			t.Cleanup(func() { putResp.Body.Close() })
			assert.Equal(t, tt.wantPutCode, putResp.StatusCode)
			// assert after state
			if is2xx(tt.wantPutCode) {
				var v struct {
					ID string `json:"ID"`
				}
				err = json.NewDecoder(putResp.Body).Decode(&v)
				assert.NoError(t, err)
				tt.assertAfter(v.ID, t)
			}
		})
	}
}

func TestGoAPIOption_Run_POST_Sample(t *testing.T) {
	type testcase struct {
		id             string
		url            string
		wantCode       int
		assertBefore   assertByIDFunc
		wantBefore     string
		wantBeforeCode int
		wantAfter      string
		assertAfter    assertByIDFunc
		wantAfterCode  int
	}
	tests := map[string]testcase{
		"no id": {
			id:           "",
			url:          "http://localhost:8080/sample?name=post-test&birthday=2000-01-01&is_japanese=true",
			wantCode:     http.StatusBadRequest,
			assertBefore: nopAssertByID(),
			assertAfter:  nopAssertByID(),
		},
		"existing sample id with name, birthday, is_japanese": {
			id:           "00000000-0000-0000-0000-000000000004",
			url:          "http://localhost:8080/sample?id=00000000-0000-0000-0000-000000000004&name=post-test&birthday=2000-01-01&is_japanese=true",
			wantCode:     http.StatusOK,
			assertBefore: getAndAssertWith(SampleJSON_4+"\n", http.StatusOK),
			assertAfter:  getAndAssertWith(`{"ID":"00000000-0000-0000-0000-000000000004","Name":"post-test","Birthday":"2000-01-01T00:00:00+09:00","IsJapanese":true}`+"\n", http.StatusOK),
		},
		"existing sample id with no change": {
			id:           "00000000-0000-0000-0000-000000000004",
			url:          "http://localhost:8080/sample?id=00000000-0000-0000-0000-000000000004",
			wantCode:     http.StatusOK,
			assertBefore: getAndAssertWith(SampleJSON_4+"\n", http.StatusOK),
			assertAfter:  getAndAssertWith(SampleJSON_4+"\n", http.StatusOK),
		},
		"non-existing sample id": {
			id:           "00000000-0000-0000-0000-000000000010",
			url:          "http://localhost:8080/sample?id=00000000-0000-0000-0000-000000000010&name=post-test",
			wantCode:     http.StatusOK,
			assertBefore: getAndAssertWith("", http.StatusInternalServerError),
			assertAfter:  getAndAssertWith("", http.StatusInternalServerError),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			setup(context.Background(), context.Background(), t)
			// check before state
			tt.assertBefore(tt.id, t)
			// run
			req, err := http.NewRequest(http.MethodPost, tt.url, nil)
			assert.NoError(t, err)
			want := fmt.Sprintf(`{"id":%q}`+"\n", tt.id)
			doWithAssert(req, want, tt.wantCode, t)
			// check after state
			tt.assertAfter(tt.id, t)
		})
	}
}

func TestGoAPIOption_Run_DELETE_Sample(t *testing.T) {
	type testcase struct {
		id           string
		url          string
		assertBefore assertByIDFunc
		wantCode     int
		assertAfter  assertByIDFunc
	}
	tests := map[string]testcase{
		"no id": {
			id:           "",
			url:          "http://localhost:8080/sample",
			assertBefore: getAndAssertWith("", http.StatusBadRequest),
			wantCode:     http.StatusBadRequest,
			assertAfter:  getAndAssertWith("", http.StatusBadRequest),
		},
		"existing sample id": {
			id:           "00000000-0000-0000-0000-000000000004",
			url:          "http://localhost:8080/sample?id=00000000-0000-0000-0000-000000000004",
			assertBefore: getAndAssertWith(SampleJSON_4+"\n", http.StatusOK),
			wantCode:     http.StatusOK,
			assertAfter:  getAndAssertWith("", http.StatusInternalServerError),
		},
		"non-existing sample id": {
			id:           "00000000-0000-0000-0000-000000000010",
			url:          "http://localhost:8080/sample?id=00000000-0000-0000-0000-000000000010",
			assertBefore: getAndAssertWith("", http.StatusInternalServerError),
			wantCode:     http.StatusOK,
			assertAfter:  getAndAssertWith("", http.StatusInternalServerError),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			setup(context.Background(), context.Background(), t)
			// check before state
			tt.assertBefore(tt.id, t)
			// run
			req, err := http.NewRequest(http.MethodDelete, tt.url, nil)
			assert.NoError(t, err)
			doWithAssert(req, "", tt.wantCode, t)
			// check delete
			tt.assertAfter(tt.id, t)
		})
	}
}

// doWithAssert requests using req and then assert as follow:
//   - there is no error
//   - response status code is equal to wantCode
//   - response body is equal to want only if status code is 2xx
//
// The response body will be closed in t.Cleanup.
func doWithAssert(req *http.Request, want string, wantCode int, t *testing.T) {
	resp, err := httpClient.Do(req)
	assert.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	got, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, wantCode, resp.StatusCode)
	if is2xx(resp.StatusCode) {
		assert.Equal(t, want, string(got))
	}
}

// assertGetByID makes GET request to /sample with id query then execute doWithAssert.
func assertGetByID(id, want string, wantCode int, t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/sample?id="+id, nil)
	assert.NoError(t, err)
	doWithAssert(req, want, wantCode, t)
}

type assertByIDFunc func(id string, t *testing.T)

func nopAssertByID() assertByIDFunc {
	return func(id string, t *testing.T) {}
}

// getAndAssertWith returns assertByIDFunc which checks body and status code of GET request is equal to want and wantCode each other.
func getAndAssertWith(want string, wantCode int) assertByIDFunc {
	return func(id string, t *testing.T) {
		assertGetByID(id, want, wantCode, t)
	}
}

func is2xx(code int) bool {
	return 200 <= code && code < 300
}

func sampleJSON(id string, name string, birthday string, isJapanese bool) string {
	return fmt.Sprintf(`{"ID":%q,"Name":%q,"Birthday":%q,"IsJapanese":%v}`, id, name, birthday, isJapanese)
}

func TestMain(m *testing.M) {
	if s, err := os.ReadFile(INIT_SCRIPT); err == nil {
		log.Printf("Init MySQL Container: \n%v\n", string(s))
	} else {
		log.Printf("Failed to read init script %s: %v\n", INIT_SCRIPT, err)
	}
	os.Exit(m.Run())
}

func setup(containerCtx, appCtx context.Context, t *testing.T) {
	container, err := startMySQLContainer(containerCtx, t)
	assert.NoError(t, err)
	t.Cleanup(func() {
		if err := container.Terminate(containerCtx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})
	dsn, err := container.ConnectionString(containerCtx)
	assert.NoError(t, err)
	// run go-api
	appCtx, stop := context.WithCancel(appCtx)
	t.Cleanup(stop)
	go func() {
		assert.NoError(t, (&GoAPICmd{
			MySQL:  MySQLOption{Table: SAMPLE_TABLE, DSN: dsn},
			Server: ServerOption{Port: "8080"},
			Log:    LogOption{SlogLevel{slog.LevelError}},
		}).Run(appCtx))
	}()
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
