# GOIT - A Golang integration test kit
Integration test is very important to make sure repository layer works smoothly with database and do its business jobs correctly. However, integration test usually comes with high cose of development and slow running if you do not have a smart fixture strategy and take advantages of Go parallel test feature. A good integration test kit will encourage
your team to write more quality and maintainable tests with a low development cost.

In a nutshell, what benefit this library would offer to you:

- Reusable fixture strategy with nested builder solution.
- Parallelism support with file based SQLite and transaction based SQL.
- Well structured tests in suites with setup and teardown, thanks to [Testify Suite](https://pkg.go.dev/github.com/stretchr/).
- The kit can be extended easily in your code.

## Reusable fixtures with nested builder
Suppose that you are working with the below entities in a service, in which a `product` depends on `category` and `warehouse` then a `warehouse` also has `city` as a dependency in the chain.

    product   [ id | name | category_id | warehouse_id ]
    warehouse [ id | name | city_id ]
    category  [ id | name ]
    city      [ id | name ]

The nested builder solution with combination of two design patterns abstract method and decorator enables you to make your fixtures simple and reusable. The fixtures can be reused in either tests or other fixtures easily. I highly recommend to combine this kit with a fixture factory to make your life easier.

In particular, we can create a two nested product fixtures as bellow:

    product.AdidasBallProductFixture struct
    ...........................................
    . Adidas ball product                     .
    . +++++++++++++++++++                     .
    . + Mitte warehouse +                     .
    . + *************** +                     .
    . + * Berlin City * +  """""""""""""""""" .
    . + *************** +  " Sport category " .
    . +++++++++++++++++++  """""""""""""""""" .
    ...........................................  

    product.MacbookProProductFixture struct
    ............................................
    . Macbook pro product                      .
    . +++++++++++++++++++                      .
    . + Mitte warehouse +                      .
    . + *************** +                      .
    . + * Berlin City * +  ~~~~~~~~~~~~~~~~~~~ .
    . + *************** +  ~ Hitech category ~ .
    . +++++++++++++++++++  ~~~~~~~~~~~~~~~~~~~ .
    ............................................ 

Definitely, we also have to create fixture structs for `Sport category`, `Hitech category`, `Berlin city` and `Mitte warehouse`.
But to keep it simple for now, we just assume that they are created already and we just simply reuse them as nested fixtures in `product.AdidasBallProductFixture` and `product.MacbookProProductFixture` structs.

Last step, let's see how easy to write tests the following example code.

    /* 
    Here you can extend your test from either file based SQLITE goit.ITsqlite or transaction based SQL goit.ITsql 
    */
    type productRepoTestSuite struct {
        goit.ITsqlite
        repo *productRepo
    }

    func (s *productRepoTestSuite) TesProductRepoGetByID_NoError {
        /*
        Setup Adidas ball product fixture eaisily with just one line of code
        You do not need to care about how to build fixture of category, city and warehouse
        Because product.AdidasBallProductFixture already takes care of them inside it.
        */
        product.NewAdidasBallProductFixture(s.DB()).Build()
        
        /* Define you expectation */
        expectedResult := *model.Product{ID: 1, Name: "Adidas ball", CategoryID: 1, WarehouseID: 1}
        
        /* Get the created fixture from a fixture store easily */
        adidasBall := s.GetFixture(product.AdidasBallFixtureReference)
        
        /* Call your function
        actualResult, err := s.productRepo.GetByID(adidasBall.(*model.Product).ID)
        
        /* Test your result */
        s.NoError(err)
        s.Equal(expectedResult, actualResult)
    }

## Test parallelism
Integration test is usually slow so we need to speed them up by taking advantages of Go test parallel feature. 
Basically we have two ways to run integration tests in parallel.
1. File based SQLite, each test case has it own separate database file.
2. Transaction based SQL, each test case starts a transaction and rollback the database when the test is done.

This tool already takes care all for you and you do not need to do anything else to implement parallel tests.

## Example code
More detail of how to use this tool kit is provided in example folder. You just simply clone this kit and run `go test ./example/repository` on a machine already have Go installed without database or any other requirements.

**__Ps: You are welcome to contribute your pull request for improvement.__**