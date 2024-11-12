### TODO list
    - manual testing after the completion of the each sub task or general task 
    - structured logging at each step
    - appropriate error handling
    - task must be done in the order they provided 

    - [x] Routing

    - [ ] Handlers
        - [ ] Order handler
        - [ ] Menu handler
        - [ ] Inventory Handler 
    
    - [ ] Services
        - [ ] Order service
        - [ ] Menu service
        - [ ] Inventory service 

    - [ ] DAL
        - [ ] Inventory storage
        - [ ] Menu storage
        - [ ] Order storage 
    

    - [ ] Aggregation functions
    - [ ] Flag

### Small documentation
/internal/core directory contains the essence of the project, its domain entities and interfaces of the services and repositories.
/internal/infrastructure directory contains the controllers(handlers), data access layer implementation of repositories and services, data models and mappers between model and entity

### Backlog
    - [ ] Streaming reading
