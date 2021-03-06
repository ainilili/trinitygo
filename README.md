# trinitygo

[![Build Status](https://api.travis-ci.org/PolarPanda611/trinitygo.svg)](https://travis-ci.org/PolarPanda611/trinitygo)
[![Go Report Card](https://goreportcard.com/badge/github.com/PolarPanda611/trinitygo)](https://goreportcard.com/report/github.com/PolarPanda611/trinitygo)
[![GoDoc](https://godoc.org/github.com/PolarPanda611/trinitygo?status.svg)](https://godoc.org/github.com/PolarPanda611/trinitygo)
[![Release](https://img.shields.io/github/release/PolarPanda611/trinitygo.svg?style=flat-square)](https://github.com/PolarPanda611/trinitygo/releases)

golang restframework plugin with gin+gorm, fast and high scalable    

## Installation

```bash
$ go get github.com/PolarPanda611/trinitygo/trinitygo
$ trinitygo NewHttp [Your Project Name]
$ cd [Your Project Name] 
$ trinitygo NewCrud [Your Model Name]
$ swag init

// start your journey in Trinity
// you can check the demo under example folder 
```

done.

## Feature

* integrate gorm
* integrate gin
* fast register router
* customize middleware
* customize runtime (Tracing Analysis, user authentication , event bus ...)
* support atomic request 
* support customize validator ( API permission , data validation ...)
* support URL query analyze (search , filter , order by , preload ...)
* full dependency injection support
* cmd to init Project folder
* cmd to init Project CRUD Code && unit test template
* cmd to init Swagger docs

## Docs

## Http Server 

* start server 

```new HTTP server 
// can see the example in example/http

trinitygo.SetConfigPath(configPath)  // put config path here
// by default the health check is disabled 
//trinitygo.EnableHealthCheckURL("/v1/ping")  ==> // the health check path will be /baseurl/v1/ping
//trinitygo.EnableHealthCheckURL()            ==> // if the path not set , the check path will be default /baseurl/ping
//trinitygo.SetHealthCheckDefaultHandler(handler) == > also can use this func to change the default health check handler
t := trinitygo.DefaultHTTP()
t.InitHTTP()
t.ServeHTTP()

```

* integrate gorm

```install gorm
// by default , trinitygo will install the gorm according the config 
// config.toml
[database]
db_type = "postgres" #mysql  postgres
server = "host=127.0.0.1 port=60901 user=trinity password= dbname=trinity sslmode=disable" #mysql option =  charset=utf8&parseTime=True&loc=Local
table_prefix =  "trinity_"
max_idle_conn =  10
max_open_conn =  100

```

```customize install gorm 

f := func() *gorm.DB {
    return db.DefaultInstallGORM(
        app.config.GetDebug(),
        true,
        app.config.GetDBType(),
        app.config.GetDBTablePrefix(),
        app.config.GetDBServer(),
        app.config.GetDbMaxIdleConn(),
        app.config.GetDbMaxOpenConn(),
    )
}

t := trinitygo.DefaultHTTP()
t.InstallDB(f)
t.InitHTTP()
t.ServeHTTP()

```

* fast register router

```RegisterController 

// example/http/domain/controller/http
// @RegisterController to trinitygo 
// when init trinity  and the trinity 
// will auto register as root controller to 
// the gin router 
// @application.NewRequestMapping 
// the new request mapping will add new router 
// to trinity router 
// @this example 
// the instance has to be struct 
// add the tag name as you want , 
// and you can also leave it as blank 
// router : 
// GET --->  /users/:id   ==> GET func as the handler  
// GET --->  /users       ==> Getsssss func as the handler  
func init() {
	trinitygo.RegisterController("/users",userControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GET", PermissionValidator([]string{"manager"}), gValidator, g1Validator),
		application.NewRequestMapping(httputil.GET, "", "Getsssss"),
	)
}


```

```RegisterInstance 
// You can bind your service and repository layer to instance 
// The instance will auto dependency injection to your 
// controller 
// the instance has to be struct 
// add the tag name as you want , 
// and you can also leave it as blank 
func init() {
	trinitygo.RegisterInstance(userServiceImpl{} , "xxxx")
}

```


* customize middleware
``` Customize Middleware 

// by default , trinity will register 
// the following middleware 
// runtime middleware 
// logger middleware 
// recovery middleware
// you can also add your middleware 
// e.g : 
// authentication middleware ...
...
t.UseMiddleware(mlogger.New(app))
t.UseMiddleware(httprecovery.New(app))
t.UseMiddleware(mruntime.New(app))
...

```

* customize runtime 

``` Customize Runtime 
// register your runtime key before trinity serve the service 
//@param "trace_id"  --> key of the runtime 
//@param false       --> if required , when the runtime is not existed in request
//@param func        --> to create default value 
//@param islog       --> if this value will logging  
// usage : 
// Tracing Analysis : the trace_id will be passing in the whole lifecycle 
// in your request  , log middleware ,db logger ... 
// DB callback : register in db callback , see callback in db/install (auto manager the updatetime , update user etc ..)
t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", false, func() string { return uuid.New().String() }, true))
t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", false, func() string { return "" }, false))
t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", false, func() string { return "" }, true))

```

* support automic request 
```Automic request 
// set true for atomic_request to 
// open the automic_request , 
// your request will auto be wrapped 
// in one transaction 
// if you response err , 
// the transaction will be auto rollbacked 
// if response normal , 
// the transaction will be auto commit 
[app]
...
atomic_request = true
...



// if the controller you add the tag transaction 
// this tag will replace the automic request 
// if transaction is true , will get the db with tx 
// if transaction is false , will get the db without tx
// if tag set autowired , system will auto inject this field 
// with registered instance 
// if want to point to the specfic instance 
// add resource tag 
type userControllerImpl struct {
	UserSrv service.UserService `autowired:"true" resource="xxxxx"`
	Tctx    application.Context `autowired:"true" transaction:"false"`
}


```

* support customize validator 
``` Customize validator 
// see the example/http/domain/controller
// the validatior will run sorted 
// PermissionValidator -> gValidator -> g1Validator
application.NewRequestMapping(httputil.GET, "/:id", "GET", PermissionValidator([]string{"manager"}), gValidator, g1Validator),


var gValidator = func(tctx application.Context) error {
	id, _ := strconv.Atoi(tctx.GinCtx().Param("id"))
	if id < 3 {
		return errors.New("gValidator no permission")
	}
	return nil 
}

var g1Validator = func(tctx application.Context) {
	id, _ := strconv.Atoi(tctx.GinCtx().Param("id"))
	if id > 3 {
		return errors.New("g1Validator no permission")
	}
	return nil 
}

// PermissionValidator example validator
func PermissionValidator(requiredP []string) func(application.Context) {
	return func(c application.Context) {
		// c.GinCtx().Set("permission", []string{"employee"}) // no permission
		c.GinCtx().Set("permission", []string{"employee", "manager"}) // ok
		in := util.SliceInSlice(requiredP, c.GinCtx().GetStringSlice("permission"))
		if !in {
			c.httpResponseUnauthorizedErr(errors.New("np permission"))
		}
	}
}


```

* support URL query analyze 

   * URL Query config 
```
	_userConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
        // TablePrefix : table prefix defined in config.toml
        // the model User will be treated as "trinitygo_user"
		TablePrefix:  "trinitygo_",
        // DbBackend: handle the authorization 
        // no effect with the url change 
        // the base filter
		DbBackend:    nil,
		PageSize:     20,
		FilterList:   []string{"user_name", "user_name__ilike"},
		OrderByList:  []string{"id"},
		SearchByList: []string{"user_name", "email"},
		FilterCustomizeFunc: map[string]interface{}{
			"test": func(db *gorm.DB, queryValue string) *gorm.DB {
				fmt.Println("Where xxxxx = ?", queryValue)
				return db.Where("xxxxx = ?", queryValue)
			},
		},
		IsDebug: true,
	}
)


* full dependency injection support
    * path_param get value of path in your router 
	example : 
	case 1 : router: /api/user/:id          //path : id 
			 realize: /api/user/1234
			 using tag : ID int64 `path_param:"id"` to get ID with path value :1234
```
	//Controller 
	func (c *userControllerImpl) GetUserByID(args struct {
		ID int64 `path_param:"id"`
	}) {
		res, err := c.UserSrv.GetUserByID(args.ID)
		c.Tctx.HTTPResponseOk(res, err)
		return
	}

```

	* query_param get value of query in your router 
	example : 
	case 1. router: /api/user    
			 realize: /api/user/:id?name=xxx&phone=123
			 using tag : Query string `query_param:""` to get full query with query value : Query =  "name=xxx&phone=123"

	case 2. router: /api/user     
			 realize: /api/user/:id?name=xxx&phone=123
			 using tag : Name string `query_param:"name"` to get  query with query value of name : Query =  "xxx"
			 using tag : Phone string `query_param:"phone"` to get  query with query value of phone : Query =  "123"
```
	//Controller 
	func (c *userControllerImpl) GetUserList(args struct {
		Query  string `query_param:""`
		Name   string `query_param:"name"`
		Phone  string `query_param:"phone"`
	}) {
		res, err := c.UserSrv.GetUserList(args.Query)
		c.Tctx.HTTPResponseOk(res, err)
		return
	}

```


	* body_param get value of query in your router 
	example : 
	case 1. router: /api/user/:id     
			 realize: /api/user/123  body :{"username":"123"}
			 using tag : ID int64 `path_param:"id"`  to get full query with query value : ID =  "123"
			 using tag : User model.User `body_param:""`  to get full body convert to struct with request body :User =  User{Username:"123"}
			 using tag : Username string `body_param:"username"`  to get username with request body by body param key  : Username = "123"

```
	//Controller 
	func (c *userControllerImpl) CreateUser(args struct {
		User model.User `body_param:""`
	}) {
		res, err := c.UserSrv.CreateUser(&args.User)
		c.Tctx.httpResponseCreated(res, err)
		return
	}

```


## GRPC Server 
```new GRPC server 
// can see the example in example/server

trinitygo.SetConfigPath(configPath) // put config path here
t := trinitygo.DefaultGRPC()
t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", true, func() string { return "" }, true))
t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", true, func() string { return "" }, true))
t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", true, func() string { return "" }, true))
t.InitGRPC()
{
    helloworldpb.RegisterGreeterServer(t.GetGRPCServer(), &grpc.Server{})  // register your grpc server here
}
t.ServeGRPC()

```


// More detail see the example 