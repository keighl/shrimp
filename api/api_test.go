package api

import (
  "shrimp/utils"
  "shrimp/models"
  "testing"
  "net/http/httptest"
  "github.com/go-martini/martini"
  _ "github.com/jrallison/go-workers"
  "github.com/modocache/gory"
  "reflect"
  "github.com/keighl/mandrill"
)

var (
  alreadySetup bool
  err error
)

// TODO find a library that does setup/tear down instead of this
func setup(t *testing.T) {
  if (alreadySetup) { return }

  Config    = utils.ConfigForFile("../conf/test.json")
  DB        = utils.DBForConfig(Config)
  models.DB = DB
  // workers.Configure(map[string]string{
  //   "server": config.WorkerServer,
  //   "database": config.WorkerDatabase,
  //   "pool": config.WorkerPool,
  //   "process": config.WorkerProcess,
  // })

  DB.Exec("TRUNCATE TABLE users")
  DB.Exec("TRUNCATE TABLE todos")
  DB.Exec("TRUNCATE TABLE password_resets")
  models.DefineFactories()
  alreadySetup = true
}

func testTools(t *testing.T) (*martini.ClassicMartini, *httptest.ResponseRecorder) {
  setup(t)
  return utils.MartiniServer(false), httptest.NewRecorder()
}

func expect(t *testing.T, a interface{}, b interface{}) {
  if a != b {
    t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

func refute(t *testing.T, a interface{}, b interface{}) {
  if a == b {
    t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
  }
}

var (
  uzer *models.User // memoized user
)

func Uzer(t *testing.T) (*models.User) {
  if (uzer != nil) { return uzer }

  uzer = gory.Build("user").(*models.User)
  err = DB.Create(uzer).Error
  if (err != nil) { t.Error(err) }

  return uzer
}

func MockMailerTrue(c martini.Context) {
  c.Map(SendEmail(func (message *mandrill.Message)(bool) { return true }))
}

func MockMailerFalse(c martini.Context) {
  c.Map(SendEmail(func (message *mandrill.Message)(bool) { return false }))
}


//////////////////////////////
// API ENVELOPE //////////////

func Test_Api500Envelope(t *testing.T) {
  setup(t)
  data := new(ApiData)
  data.ApiError = &ApiError{"There was an unexpected error!", []string{}}
  targetByHand := ApiEnvelope{data}
  targetByMethod := Api500Envelope()
  expect(t, reflect.DeepEqual(targetByHand, targetByMethod), true)
}