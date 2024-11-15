### TODO list
    - manual testing after the completion of the each sub task or general task 
    - structured logging at each step
    - appropriate error handling
    - task must be done in the order they provided 

    - validation must include: 
    1. Empty value check
    2. ID divergence check
    3. ID collision check 
    4. Items dudplicate check 
    5. Negative quantity check 
    6. Empty items check

    - [x] Routing
    
    - [ ] Inventory
        - [x] Inventory storage
        - [x] Inventory service 
        - [x] Inventory Handler 
        - [x] Validation
        - [ ] Manual testing

    - [ ] Menu 
        - [x] Menu storage
        - [x] Menu service
        - [x] Menu handler
        - [x] Validation
        - [ ] Manual testing

    - [ ] Orders
        - [] Order service
        - [x] Ingridients sufficiency validation
        - [x] Ingridients deduction
        - [x] Order storage 
        - [x] Order handler
        - [x] Validation
        - [x] Manual testing

    - [ ] Aggregation functions
    - [ ] Flag
    - [ ] Manual testing



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

