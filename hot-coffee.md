### TODO list
    - manual testing after the completion of the each sub task or general task 
    - structured logging at each step
    - appropriate error handling
    - task must be done in the order they provided 

    - validation must include: 
    1. Empty value check
    2. ID collision check 
    3. Items dudplicate check 
    4. Negative quantity check 
    5. Empty items check

    - [x] Routing
    
    - [ ] Inventory
        - [x] Inventory storage
        - [x] Inventory service 
        - [x] Inventory Handler 
        - [ ] Strong Validation
        - [ ] Manual testing

    - [ ] Menu 
        - [x] Menu storage
        - [x] Menu service
        - [x] Menu handler
        - [ ] Strong Validation
        - [ ] Manual testing

    - [ ] Orders
        - [ ] Order service
        - [ ] Order storage 
        - [ ] Order handler
        - [ ] Strong Validation
        - [ ] Manual testing

    - [ ] Aggregation functions
    - [ ] Flag
    - [ ] Reread the project



### Project structure
    - /internal/core contains domain related entities
    - /internal/repositories contains repository interfaces and implementations
    - /internal/services contains interfaces of services of the project and its implementations
    - /internal/infrastructure contains controllers and different handlers not related to service directly
    - /internal/utils contains utilities used by different packages
    - /internal/flag contains program argument parsing

### Backlog
    - [ ] Refactor:
        - [ ] Generalize many methods
    - [ ] Streaming reading of JSON files
    - [ ] Unit, integration tests
    - [ ] Dependency injection
    - [ ] Logger instance
    - [ ] Validator

