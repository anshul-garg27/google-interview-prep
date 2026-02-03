Core Concepts
=============

**Rest Layer**
----------
 - Provides a way to construct a service for a resource 
 - A resource is exposed as Entry object, and Persisted as Entity object
 - For example: CollectionEntry is exposed to the outside world while CollectionEntity is persisted to storage

 - **Layers**
    
    - **API**
        - Entry point of Request
        - Defines the routes for a given resource
        - Parses input and calls Service
        - Renders Response object from Service into JSON

 	- **Service**
        - Called by API, Talks to Manager
 		- Talks to Manager
 		- Maps Request Inputs to Entry objects and calls Manager
 		- Wraps Entries received from Manager into Response Objects and sends back to API

 	- **Manager**
        - Called by Service, Talks to DAO
 		- Home to business logic for a given resource
 		- Talks to DAO layer to fetch/persist entity objects
 		- Sends Entry object back to service

 	- **DAO**
 		- Called by Manager
 		- Delegates stuff to DB specific Dao Implementation
		- Available Delegates
		 	- PG Dao (for Postgres)
                - Uses Session inside App Context as the Entity Manager / Transaction context


- **Request Lifecycle Objects**
    - **Request Context**
        - Context Holder which is set by a middleware during request lifecycle
        - Contains 
            - Session (Transaction Context)
            - Request Headers
            - Other relevant request metadata (AccountID/PartnerID/IsLoggedId/etc)
            - Events pushed during request lifecycle   
                - These are later published _after successful commit_ of transaction
    
    - **Session**
        - Transaction Context set by a middleware during initial phase of request lifecycle
        - Session's responsibility is to
            - Manage the transaction object - Commit/Rollback transaction
                - Rollback is done in case of
                    - Panic
                    - Timeout
            - Perform after commit callback if-any
        - Currently Implemented: PG Session for Postgres
