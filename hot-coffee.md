### TODO list
    - manual testing after the completion of the each sub task or general task 
    - structured logging at each step
    - appropriate error handling
    - task must be done in the order they provided 

    - [x] Routing
    
    - [ ] Inventory
        - [x] Inventory storage
        - [x] Inventory service 
        - [x] Inventory Handler 
        - [ ] Manual testing

    - [ ] Menu 
        - [x] Menu storage
        - [x] Menu service
        - [x] Menu handler
        - [ ] Manual testing

    - [ ] Orders
        - [ ] Order service
        - [ ] Order storage 
        - [ ] Order handler
        - [ ] Manual testing

    - [ ] Aggregation functions
    - [ ] Flag

### Project structure
    - /internal/core contains domain related entities
    - /internal/repositories contains repository interfaces and implementations
    - /internal/services contains interfaces of services of the project and its implementations
    - /internal/infrastructure contains controllers and different handlers not related to service directly
    - /internal/utils contains utilities used by different packages
    - /internal/flag contains program argument parsing

### Backlog
    - [ ] Streaming reading of JSON files
    - [ ] Unit, integration tests
    - [ ] Dependency injection
    - [ ] Logger instance
    - [ ] Validator

